package controllers

import (
	"bytes"
	"cms-backend/models"
	"cms-backend/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// TODO: Import required packages for:
// - HTTP testing (net/http, httptest)
// - JSON handling (encoding/json)
// - Database mocking (sqlmock)
// - Time handling
// - Your application packages (models, utils)

func TestGetPages(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)

	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "updated_at"}).
		AddRow(1, "First Page", "Content 1", time.Now(), time.Now()).
		AddRow(2, "Second Page", "Content 2", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "pages"`).WillReturnRows(rows)

	router.GET("/pages", GetPages)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/pages", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, but got %d", w.Code)
	}

	var response []models.Page
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}
	if len(response) != 2 {
		t.Fatalf("Expected 2 pages, but got %d", len(response))
	}
}

func TestGetPage(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)

	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "updated_at"}).
		AddRow(1, "First Page", "Content 1", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "pages" WHERE "pages"\."id" = \$1 ORDER BY "pages"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	router.GET("/pages/:id", GetPage)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/pages/1", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	var response models.Page
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Failed to unmarshal response")

	assert.Equal(t, uint(1), response.ID, "Expected page ID to be 1")
	assert.Equal(t, "First Page", response.Title, "Page title should match")
	assert.Equal(t, "Content 1", response.Content, "Page content should match")
	assert.NoError(t, mock.ExpectationsWereMet(), "All expectations were not met")
}

func TestCreatePage(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)

	router.POST("/pages", CreatePage)

	page := models.Page{
		Title:   "Test Page",
		Content: "Content 1",
	}

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "updated_at"}).
		AddRow(1, "Test Page", "Content 1", time.Now(), time.Now())

	mock.ExpectQuery(`INSERT INTO "pages"`).
		WithArgs(
			page.Title,
			page.Content,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).WillReturnRows(rows)

	mock.ExpectCommit()

	w := httptest.NewRecorder()

	body, _ := json.Marshal(page)

	req := httptest.NewRequest("POST", "/pages", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Expected status code 201")

	var returnedPage models.Page

	err := json.Unmarshal(w.Body.Bytes(), &returnedPage)

	assert.NoError(t, err, "Failed to unmarshal response")
	assert.Equal(t, uint(1), returnedPage.ID, "Expected page ID to be 1")
	assert.Equal(t, page.Title, returnedPage.Title, "Page title should match")
	assert.Equal(t, page.Content, returnedPage.Content, "Page content should match")
	assert.NoError(t, mock.ExpectationsWereMet(), "All expectations were not met")
}

func TestUpdatePage(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)

	router.PUT("/pages/:id", UpdatePage)

	selectRows := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "updated_at"}).
		AddRow(1, "Old Title", "Old Content", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "pages" WHERE "pages"\."id" = \$1 ORDER BY "pages"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(selectRows)

	mock.ExpectBegin()

	mock.ExpectExec(`UPDATE "pages" SET "title"=\$1,"content"=\$2,"updated_at"=\$3 WHERE "id" = \$4`).
		WithArgs("Updated Title", "Updated Content", sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	reloadRows := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "updated_at"}).
		AddRow(1, "Updated Title", "Updated Content", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "pages" WHERE "id" = \$1`).
		WithArgs(1).
		WillReturnRows(reloadRows)

	mock.ExpectCommit()

	updateData := models.Page{
		Title:   "Updated Title",
		Content: "Updated Content",
	}
	body, _ := json.Marshal(updateData)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/pages/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	var returnedPage models.Page
	err := json.Unmarshal(w.Body.Bytes(), &returnedPage)

	assert.NoError(t, err, "Failed to unmarshal response")
	assert.Equal(t, uint(1), returnedPage.ID, "Expected page ID to remain 1")
	assert.Equal(t, updateData.Title, returnedPage.Title, "Page title should be updated")
	assert.Equal(t, updateData.Content, returnedPage.Content, "Page content should be updated")
	assert.NoError(t, mock.ExpectationsWereMet(), "All expectations were not met")
}

func TestDeletePage(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	defer mock.ExpectClose()

	router.DELETE("/pages/:id", DeletePage)

	selectRows := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "updated_at"}).
		AddRow(1, "Test Title", "Content 1", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "pages" WHERE "pages"\."id" = \$1 ORDER BY "pages"\."id" LIMIT \$2`).
		WithArgs(1, 1).WillReturnRows(selectRows)

	mock.ExpectBegin()

	mock.ExpectExec(`DELETE FROM "pages" WHERE "pages"\."id" = \$1`).
		WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/pages/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200.")

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.NoError(t, err, "Failed to unmarshal response")
	assert.Equal(t, "Page successfully deleted", response["message"])
	assert.NoError(t, mock.ExpectationsWereMet(), "All expectations were not met")
}

/*
TESTING HINTS:
1. Use sqlmock.AnyArg() for timestamp fields
2. Remember to escape special characters in SQL patterns
3. Each database operation needs proper error handling
4. Content-Type header is required for POST/PUT requests
5. Transaction tests need Begin/Commit expectations
6. Use proper argument matching in mock expectations
7. Consider testing error cases:
   - Invalid IDs
   - Missing required fields
   - Database errors
*/
