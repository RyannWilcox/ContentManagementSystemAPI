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

func TestGetMedia(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	router.GET("/media", GetMedia)

	rows := sqlmock.NewRows([]string{"id", "url", "type", "created_at", "updated_at"}).
		AddRow(1, "http://thing1.com", "Test Site 1", time.Now(), time.Now()).
		AddRow(2, "http://thing2.com", "Test Site 2", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "media"`).WillReturnRows(rows)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/media", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Test failed with status %d. Server Error: %s", w.Code, w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	var response []models.Media
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Failed to unmarshal response")
	assert.Equal(t, 2, len(response), "Expected 2 media")
	assert.NoError(t, mock.ExpectationsWereMet(), "All expectations were not met")
}

func TestGetMediaByID(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	router.GET("/media/:id", GetMediaByID)

	rows := sqlmock.NewRows([]string{"id", "url", "type", "created_at", "updated_at"}).
		AddRow(1, "www.testsite.com", "website", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "media" WHERE "media"\."id" = \$1 ORDER BY "media"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/media/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	var response models.Media
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Failed to unmarshal response")

	assert.Equal(t, uint(1), response.ID, "Expected media ID to be 1")
	assert.Equal(t, "www.testsite.com", response.URL, "media url expected to match")
	assert.Equal(t, "website", response.Type, "Media type expected to match")
	assert.NoError(t, mock.ExpectationsWereMet(), "All expectations were not met")
}

func TestCreateMedia(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	router.POST("/media", CreateMedia)

	media := models.Media{
		URL:  "www.newmedia.com",
		Type: "website",
	}

	mock.ExpectBegin()

	rows := sqlmock.NewRows([]string{"id", "url", "type", "created_at", "updated_at"}).
		AddRow(1, "www.newmedia.com", "website", time.Now(), time.Now())

	mock.ExpectQuery(`INSERT INTO "media"`).
		WithArgs(media.URL, media.Type, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	mock.ExpectCommit()

	w := httptest.NewRecorder()
	body, _ := json.Marshal(media)

	req := httptest.NewRequest("POST", "/media", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Expected status code 201")

	var returnedMedia models.Media
	err := json.Unmarshal(w.Body.Bytes(), &returnedMedia)
	assert.NoError(t, err, "Failed to unmarshal response")
	assert.Equal(t, uint(1), returnedMedia.ID, "Expected media ID to be 1")
	assert.Equal(t, media.URL, returnedMedia.URL, "media url should match")
	assert.Equal(t, media.Type, returnedMedia.Type, "media type should match")
	assert.NoError(t, mock.ExpectationsWereMet(), "All expectations were not met")
}

func TestDeleteMedia(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	router.DELETE("/media/:id", DeleteMedia)

	rows := sqlmock.NewRows([]string{"id", "url", "type", "created_at", "updated_at"}).
		AddRow(1, "www.deleteme.com", "website", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "media" WHERE "media"\."id" = \$1 ORDER BY "media"\."id" LIMIT \$2`).
		WithArgs(1, 1).WillReturnRows(rows)

	mock.ExpectBegin()

	mock.ExpectExec(`DELETE FROM "media" WHERE "media"\."id" = \$1`).
		WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/media/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200.")

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.NoError(t, err, "Failed to unmarshal response")
	assert.Equal(t, "Media successfully deleted", response["message"])
	assert.NoError(t, mock.ExpectationsWereMet(), "All expectations were not met")
}
