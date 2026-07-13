package integration

import (
	"cms-backend/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: Import required packages for:
// - JSON handling
// - HTTP testing
// - Your application models
// - Testing package

/*
MEDIA INTEGRATION TESTS

These tests verify the complete flow of media operations through the API.
Each test should:
1. Start with a clean database state
2. Perform API operations
3. Verify the responses
4. Check database state if needed
*/

func TestMediaIntegration(t *testing.T) {
	// STEP 1: Clear Database
	// - Clear all tables before starting tests
	clearTables()

	t.Run("Create Media", func(t *testing.T) {
		//TODO: Implement test logic
		// STEP 1: Prepare Test Data
		// - Create JSON body with:
		//   * URL (e.g., "http://example.com/test.jpg")
		//   * Type (e.g., "image")
		body := `{
			"url": "http://example.com/test.jpg",
			"type": "image"
		}`

		// STEP 2: Create HTTP Request
		// - Create POST request to /api/v1/media
		// - Set Content-Type header
		// - Add request body
		req := httptest.NewRequest("POST", "/api/v1/media", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		// STEP 3: Execute Request
		// - Create response recorder
		// - Send request through router
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Expected status 201, got %d: %s", w.Code, w.Body.String())
		}

		// STEP 4: Verify Response
		// - Check status code (should be 201 Created)
		// - Parse response JSON
		// - Verify media properties (URL, type)
		// - Handle any parsing errors
		var response models.Media
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.URL != "http://example.com/test.jpg" {
			t.Errorf("Expected URL 'http://example.com/test.jpg', got %s", response.URL)
		}
	})

	t.Run("Get All Media", func(t *testing.T) {
		// STEP 1: Setup Test Data
		// - Create test media entries if needed
		clearTables()
		createTestMedia(t)
		createTestMedia(t)

		// STEP 2: Create HTTP Request
		req := httptest.NewRequest("GET", "/api/v1/media", nil)

		// STEP 3: Execute Request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// STEP 4: Verify Response
		// - Check status code (should be 200 OK)
		// - Parse response JSON array
		// - Verify media list properties
		// - Check number of items
		assert.Equal(t, http.StatusOK, w.Code, "Expected status 200 OK")

		var response []models.Media
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Failed to unmarshal response")
		assert.Equal(t, 2, len(response), "Expected 2 media items")
	})

	t.Run("Get Media by ID", func(t *testing.T) {
		clearTables()
		id := createTestMedia(t)

		// STEP 1: Create HTTP Request
		req := httptest.NewRequest("GET", "/api/v1/media/"+strconv.FormatUint(uint64(id), 10), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Expected status 200 OK")

		//	STEP 2: Verify Response
		var response models.Media
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Failed to unmarshal response")

		assert.Equal(t, "http://example.com/test.jpg", response.URL, "Expected URL 'http://example.com/test.jpg'")
		assert.Equal(t, "image", response.Type, "Expected type 'image'")
		assert.Equal(t, id, response.ID, "Expected ID to match")
	})

	t.Run("Create Media with Invalid Data", func(t *testing.T) {
		clearTables()

		// Invalid url
		invalidBody := `{
			"url": "",
			"type": "a type"
		}`
		// STEP 1: Create HTTP Request with Invalid Data
		req := httptest.NewRequest("POST", "/api/v1/media", strings.NewReader(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status 400 Bad Request")
	})

	t.Run("Delete Media", func(t *testing.T) {
		clearTables()
		id := createTestMedia(t)

		// STEP 1: Create HTTP Request for Deletion
		req := httptest.NewRequest("DELETE", "/api/v1/media/"+strconv.FormatUint(uint64(id), 10), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Expected status 200 OK")

		// STEP 2: Verify Response
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Failed to unmarshal response")
		assert.Equal(t, "Media successfully deleted", response["message"], "Expected deletion message")

		// Verify that the media is deleted
		req = httptest.NewRequest("GET", "/api/v1/media/"+strconv.FormatUint(uint64(id), 10), nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code, "Expected status 404 Not Found after deletion")
	})

	t.Run("Create Media with Duplicate URL", func(t *testing.T) {
		clearTables()
		createTestMedia(t)

		duplicateBody := `{
			"url": "http://example.com/test.jpg",
			"type": "image"
		}`

		req := httptest.NewRequest("POST", "/api/v1/media", strings.NewReader(duplicateBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Expected status 201")
	})

}

// TODO: Additional test cases to consider:
// - Get single media by ID
// - Create media with invalid data
// - Delete media
// - Create media with duplicate URL

/*
TESTING HINTS:
1. Request Creation:
   - Use httptest.NewRequest for creating requests
   - Remember to set Content-Type for POST requests
   - Use strings.NewReader for request bodies

2. Response Handling:
   - Use httptest.NewRecorder for capturing responses
   - Parse JSON responses carefully
   - Check both status codes and response bodies

3. Test Data:
   - Use meaningful test data
   - Clean up between tests
   - Consider edge cases

4. Error Cases:
   - Test invalid inputs
   - Test missing required fields
   - Test invalid content types
*/
