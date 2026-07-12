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

func TestGetPosts(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)

	router.GET("/posts", GetPosts)

	postsRows := sqlmock.NewRows([]string{"id", "title", "content", "author", "created_at", "updated_at"}).
		AddRow(uint(1), "Test Title 1", "Content 1", "The Author", time.Now(), time.Now()).
		AddRow(uint(2), "Test Title 2", "Content 2", "An Author", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT .* FROM "posts"`).
		WillReturnRows(postsRows)

	joinRows := sqlmock.NewRows([]string{"post_id", "media_id"}).
		AddRow(uint(1), uint(10)).
		AddRow(uint(2), uint(20))

	mock.ExpectQuery(`SELECT .* FROM "post_media" WHERE "post_media"\."post_id" IN \(\$1,\$2\)`).
		WithArgs(1, 2).
		WillReturnRows(joinRows)

	mediaRows := sqlmock.NewRows([]string{"id", "url"}).
		AddRow(uint(10), "https://thing1.com").
		AddRow(uint(20), "https://thing2.com")

	mock.ExpectQuery(`SELECT .* FROM "media" WHERE "media"\."id" IN \(\$1,\$2\)`).
		WithArgs(10, 20).
		WillReturnRows(mediaRows)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/posts", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	var response []models.Post
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Failed to unmarshal response")
	assert.Equal(t, 2, len(response), "Expected 2 posts")
	assert.NotEmpty(t, response[0].Media, "Expected media relationship to be preloaded for post 1")
	assert.NoError(t, mock.ExpectationsWereMet(), "All expectations were not met")
}

func TestGetPostsWithFilters(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	router.GET("/posts", GetPosts)

	postsRows := sqlmock.NewRows([]string{"id", "title", "content", "author", "created_at", "updated_at"}).
		AddRow(uint(1), "Test Title 1", "Content 1", "The Author", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT .* FROM "posts" WHERE title ILIKE \$1 AND author = \$2`).
		WithArgs("%Test%", "The Author").
		WillReturnRows(postsRows)

	joinRows := sqlmock.NewRows([]string{"post_id", "media_id"}).AddRow(uint(1), uint(10))
	mock.ExpectQuery(`SELECT .* FROM "post_media"`).WithArgs(1).WillReturnRows(joinRows)

	mediaRows := sqlmock.NewRows([]string{"id", "url"}).AddRow(uint(10), "https://thing1.com")
	mock.ExpectQuery(`SELECT .* FROM "media"`).WithArgs(10).WillReturnRows(mediaRows)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/posts?title=Test&author=The+Author", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")
	assert.NoError(t, mock.ExpectationsWereMet(), "All expectations were not met")
}

func TestGetPost(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)

	// STEP 2: Mock Data Creation
	rows := sqlmock.NewRows([]string{"id", "title", "content", "author", "created_at", "updated_at"}).
		AddRow(uint(1), "Test Title 1", "Content 1", "The Author", time.Now(), time.Now())

		// STEP 3: Database Expectations
	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE "posts"\."id" = \$1 ORDER BY "posts"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	router.GET("/posts/:id", GetPost)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/posts/1", nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	var response models.Post
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Failed to unmarshal response")

	assert.Equal(t, uint(1), response.ID, "Expected post ID to be 1")
	assert.Equal(t, "Test Title 1", response.Title, "post title should match")
	assert.Equal(t, "Content 1", response.Content, "post content should match")
	assert.Equal(t, "The Author", response.Author, "post author should match")
	assert.NoError(t, mock.ExpectationsWereMet(), "All expectations were not met")
}

func TestCreatePost(t *testing.T) {
	//TODO: Add test for CreatePost
	// STEP 1: Test Setup
	router, _, mock := utils.SetupRouterAndMockDB(t)
	router.POST("/posts", CreatePost)

	post := models.Post{
		Title:   "Test Title",
		Content: "Content 1",
		Author:  "Test Author",
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author", "created_at", "updated_at"}).
		AddRow(uint(1), "Test Title", "Content 1", "Test Author", time.Now(), time.Now())

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "posts"`).
		WithArgs(post.Title, post.Content, post.Author, sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(rows)

	mock.ExpectCommit()

	w := httptest.NewRecorder()
	body, _ := json.Marshal(post)

	req := httptest.NewRequest("POST", "/posts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Expected status code 201")

	var returnedPost models.Post

	err := json.Unmarshal(w.Body.Bytes(), &returnedPost)
	assert.NoError(t, err, "Failed to unmarshal response")
	assert.Equal(t, uint(1), returnedPost.ID, "Expected post ID to be 1")
	assert.Equal(t, post.Title, returnedPost.Title, "post title should match")
	assert.Equal(t, post.Content, returnedPost.Content, "post content should match")
	assert.Equal(t, post.Author, returnedPost.Author, "post author should match")
	assert.NoError(t, mock.ExpectationsWereMet(), "All expectations were not met")
}

func TestUpdatePost(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	router.PUT("/posts/:id", UpdatePost)

	updateData := models.Post{
		Title:   "Updated Title",
		Content: "Updated Content",
		Author:  "Updated Author",
	}

	selectRows := sqlmock.NewRows([]string{"id", "title", "content", "author", "created_at", "updated_at"}).
		AddRow(uint(1), "Old Title", "Old Content", "Old Author", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT .* FROM "posts" WHERE "posts"\."id" = \$1 ORDER BY "posts"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(selectRows)

	mock.ExpectBegin()

	mock.ExpectExec(`UPDATE "posts" SET "title"=\$1,"content"=\$2,"author"=\$3,"updated_at"=\$4 WHERE "id" = \$5`).
		WithArgs(updateData.Title, updateData.Content, updateData.Author, sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	reloadRows := sqlmock.NewRows([]string{"id", "title", "content", "author", "created_at", "updated_at"}).
		AddRow(uint(1), "Updated Title", "Updated Content", "Updated Author", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT .* FROM "posts" WHERE "id" = \$1`).
		WithArgs(1).
		WillReturnRows(reloadRows)

	mock.ExpectCommit()

	body, _ := json.Marshal(updateData)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/posts/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	var returnedPost models.Post
	err := json.Unmarshal(w.Body.Bytes(), &returnedPost)
	assert.NoError(t, err, "Failed to unmarshal response")

	assert.Equal(t, uint(1), returnedPost.ID, "Expected post ID to remain 1")
	assert.Equal(t, updateData.Title, returnedPost.Title, "Post title should be updated")
	assert.Equal(t, updateData.Content, returnedPost.Content, "post content should be updated")
	assert.Equal(t, updateData.Author, returnedPost.Author, "post author should be updated")
	assert.NoError(t, mock.ExpectationsWereMet(), "All expectations were not met")
}

func TestDeletePost(t *testing.T) {
	router, _, mock := utils.SetupRouterAndMockDB(t)
	router.DELETE("/posts/:id", DeletePost)

	selectRows := sqlmock.NewRows([]string{"id", "title", "content", "author", "created_at", "updated_at"}).
		AddRow(1, "Test Title", "Content 1", "Author 1", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE "posts"\."id" = \$1 ORDER BY "posts"\."id" LIMIT \$2`).
		WithArgs(1, 1).WillReturnRows(selectRows)

	mock.ExpectBegin()

	mock.ExpectExec(`DELETE FROM "posts" WHERE "posts"\."id" = \$1`).
		WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/posts/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.NoError(t, err, "Failed to unmarshal response")
	assert.Equal(t, "Post successfully deleted", response["message"])
	assert.NoError(t, mock.ExpectationsWereMet(), "All expectations were not met")
}
