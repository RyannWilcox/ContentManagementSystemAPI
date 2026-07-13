package integration

import (
	"cms-backend/models"
	"cms-backend/routes"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/*
INTEGRATION TEST SETUP GUIDE

This file sets up the integration test environment for your CMS backend.
It handles database connections, schema migrations, and cleanup.

Key Components:
1. Test database connection
2. Router setup
3. Schema migrations
4. Test cleanup
*/

var testDB *gorm.DB
var router *gin.Engine

func TestMain(m *testing.M) {
	// STEP 1: Environment Setup
	setup()

	// STEP 2: Run Tests
	result := m.Run()

	// STEP 3: Cleanup
	cleanup()

	os.Exit(result)
}

func setup() {
	gin.SetMode(gin.TestMode)

	godotenv.Load(".env.test") // only way i was able to get the test env vars to load correctly

	// Get the test env
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"),
		os.Getenv("TEST_DB_PORT"),
	)
	var err error
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := testDB.AutoMigrate(&models.Media{}, &models.Page{}, &models.Post{}); err != nil {
		log.Fatal("Failed to automigrate test database:", err)
	}

	router = gin.Default()
	routes.InitializeRoutes(router, testDB)
}

func cleanup() {
	// make sure post_media is dropped first..
	tables := []string{"post_media", "posts", "media", "pages"}
	for _, table := range tables {
		if err := testDB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
			log.Fatal(fmt.Errorf("failed to drop table %s: %w", table, err))
		}
	}
	sqlDB, err := testDB.DB()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err))
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatal(fmt.Errorf("failed to close test database connection: %w", err))
	}
}

func clearTables() {
	// make sure post_media is dropped first..
	tables := []string{"post_media", "posts", "media", "pages"}
	for _, table := range tables {
		if err := testDB.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
			log.Fatal(fmt.Errorf("failed to clear table %s: %w", table, err))
		}
	}
}

/*
TESTING HINTS:
1. Database Connection:
   - Use a separate test database
   - Consider environment variables for credentials
   - Handle connection errors properly

2. Table Management:
   - Drop tables in correct order (foreign key constraints)
   - Clear data between tests
   - Consider using transactions for tests

3. Error Handling:
   - Log setup/cleanup errors
   - Ensure proper resource cleanup
   - Handle database operation errors

4. Best Practices:
   - Use constants for connection strings
   - Consider test helper functions
   - Add proper logging for debugging
   - Document any required setup steps
*/
