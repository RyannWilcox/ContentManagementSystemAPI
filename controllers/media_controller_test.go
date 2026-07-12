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
	//TODO: Add test for GetMediaByID
}

func TestCreateMedia(t *testing.T) {
	//TODO: Add test for CreateMedia
}

func TestDeleteMedia(t *testing.T) {
	//TODO: Add test for DeleteMedia
}
