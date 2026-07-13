package integration

import (
	"cms-backend/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
POST INTEGRATION TESTS

These tests verify the complete flow of post operations through the API,
including relationships with media items.
Each test should:
1. Start with a clean database
2. Set up required relationships (media)
3. Perform post operations
4. Verify responses and relationships
*/

func TestPostIntegration(t *testing.T) {
	// TODO: Clear Database
	// - Clear all tables before starting tests
	clearTables()

	// TODO: Create Test Media
	// - Create media item to use in posts
	// - Store media ID for later use
	mediaId := createTestMedia(t)

	var postId uint

	t.Run("Create Post with Media", func(t *testing.T) {
		body := fmt.Sprintf(`{
			"title": "Test Post",
			"content": "This is a test post",
			"author": "Test Author",
			"media": [{"id": %d, "url": "http://example.com/test.jpg", "type": "image"}]
		}`, mediaId)

		req := httptest.NewRequest("POST", "/api/v1/posts", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Expected status 201 Created")

		var response models.Post
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Failed to unmarshal response")

		assert.Equal(t, "Test Post", response.Title, "Expected title should match")
		assert.Equal(t, "This is a test post", response.Content, "Expected content should match")
		assert.Equal(t, "Test Author", response.Author, "Expected author should match")
		assert.Len(t, response.Media, 1, "Expected one media item associated with the post")
		assert.Equal(t, mediaId, response.Media[0].ID, "Expected media ID should match")
		assert.Equal(t, "http://example.com/test.jpg", response.Media[0].URL, "Expected media URL should match")
		assert.Equal(t, "image", response.Media[0].Type, "Expected media type should match")

		postId = response.ID
	})

	t.Run("Get Posts with Filter", func(t *testing.T) {
		// Create a second post with a different author, so we can prove the
		// filter actually excludes non-matching posts rather than just
		// happening to include the one we're looking for.

		body := `{
			"title": "Post",
			"content": "This post should appear in the filtered results",
			"author": "Jane Doe"
		}`
		otherBody := `{
			"title": "Unrelated Post",
			"content": "This post should not appear in the filtered results",
			"author": "Someone Else"
		}`

		otherReq := httptest.NewRequest("POST", "/api/v1/posts", strings.NewReader(otherBody))
		otherReq.Header.Set("Content-Type", "application/json")
		otherW := httptest.NewRecorder()
		router.ServeHTTP(otherW, otherReq)

		firstReq := httptest.NewRequest("POST", "/api/v1/posts", strings.NewReader(body))
		firstReq.Header.Set("Content-Type", "application/json")
		firstW := httptest.NewRecorder()
		router.ServeHTTP(firstW, firstReq)

		if otherW.Code != http.StatusCreated {
			t.Fatalf("Failed to create unrelated post, status: %d, body: %s", otherW.Code, otherW.Body.String())
		}

		if firstW.Code != http.StatusCreated {
			t.Fatalf("Failed to create unrelated post, status: %d, body: %s", firstW.Code, firstW.Body.String())
		}

		var otherPost models.Post
		if err := json.Unmarshal(otherW.Body.Bytes(), &otherPost); err != nil {
			t.Fatalf("Failed to unmarshal unrelated post: %v", err)
		}

		// Filtering posts by author
		url := fmt.Sprintf("/api/v1/posts?author=%s", "Jane+Doe")
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		var response []models.Post
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		found := false
		unrelatedFound := false
		for _, p := range response {
			if p.ID == 3 {
				found = true
			}
			if p.ID == otherPost.ID {
				unrelatedFound = true
			}
		}

		assert.Equal(t, 1, len(response), "Expected exactly one post in filtered results")
		assert.True(t, found, "Expected to find post %d (author 'Jane Doe') in filtered results", postId)
		assert.False(t, unrelatedFound, "Expected unrelated post %d (author 'Someone Else') to be excluded from filtered results", otherPost.ID)
	})
}

// Helper function to create test media
func createTestMedia(t *testing.T) uint {
	body := `{
		"url": "http://example.com/test.jpg",
		"type": "image"
	}`

	req := httptest.NewRequest("POST", "/api/v1/media", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create test media, status: %d, body: %s", w.Code, w.Body.String())
	}

	var response models.Media
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to create test media: %v", err)
	}

	return response.ID
}

/*
TESTING HINTS:
1. Request Creation:
   - Use proper JSON formatting for relationships
   - Handle URL encoding for query parameters
   - Set appropriate headers

2. Response Validation:
   - Check both status codes and response content
   - Verify relationship data is correct
   - Validate filtered results carefully

3. Test Data:
   - Create meaningful test data
   - Handle relationships properly
   - Clean up between tests

4. Error Cases to Consider:
   - Invalid media IDs
   - Missing required fields
   - Invalid filter parameters
   - Non-existent relationships
*/
