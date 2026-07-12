package controllers

import (
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
	//TODO: Add test for GetPost
	// STEP 1: Test Setup

	// STEP 2: Mock Data Creation

	// STEP 3: Database Expectations

	// STEP 4: HTTP Test Setup
}

func TestCreatePost(t *testing.T) {
	//TODO: Add test for CreatePost
	// STEP 1: Test Setup

	// STEP 2: Mock Data Creation

	// STEP 3: Database Expectations

	// STEP 4: HTTP Test Setup
}

func TestUpdatePost(t *testing.T) {
	//TODO: Add test for UpdatePost
	// STEP 1: Test Setup

	// STEP 2: Mock Data Creation

	// STEP 3: Database Expectations

	// STEP 4: HTTP Test Setup
}

func TestDeletePost(t *testing.T) {
	//TODO: Add test for DeletePost
	// STEP 1: Test Setup

	// STEP 2: Mock Data Creation

	// STEP 3: Database Expectations

	// STEP 4: HTTP Test Setup
}
