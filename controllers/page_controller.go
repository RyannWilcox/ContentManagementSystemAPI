package controllers

import (
	"cms-backend/models"
	"cms-backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetPages retrieves all pages
func GetPages(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var pages []models.Page

	title := c.Query("title")
	author := c.Query("author")

	query := db
	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}
	if author != "" {
		query = query.Where("author = ?", author)
	}

	if err := query.Find(&pages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch pages",
		})
		return
	}

	c.JSON(http.StatusOK, pages)
}

// GetPage retrieves a specific page by ID
func GetPage(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	strId := c.Param("id")

	numId, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID provided : " + err.Error(),
		})
		return
	}

	var page models.Page
	if result := db.First(&page, numId).Error; result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Page could not be found.",
		})
		return
	}

	c.JSON(http.StatusOK, page)
}

// CreatePage creates a new page
func CreatePage(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var page models.Page
	if err := c.ShouldBindJSON(&page); err != nil {
		c.JSON(http.StatusBadRequest, utils.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	if err := page.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, utils.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	tx := db.Begin()

	if err := tx.Create(&page).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Unable to add a page to the database : " + err.Error(),
		})
		return
	}
	tx.Commit()

	c.JSON(http.StatusCreated, page)
}

// UpdatePage updates an existing page by ID
func UpdatePage(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	strId := c.Param("id")

	numId, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID provided : " + err.Error(),
		})
		return
	}

	var page models.Page
	if result := db.First(&page, numId); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.HTTPError{
			Code:    http.StatusNotFound,
			Message: "Page could not be found",
		})
		return
	}

	var inputData models.Page
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, utils.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid input data provided :" + err.Error(),
		})
		return
	}

	if result := inputData.Validate(); result.Error != nil {
		c.JSON(http.StatusBadRequest, utils.HTTPError{
			Code:    http.StatusBadRequest,
			Message: result.Error(),
		})
		return
	}

	tx := db.Begin()

	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to start transaction",
		})
		return
	}

	if err := db.Model(&page).Updates(inputData).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, utils.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Database operation failed : " + err.Error(),
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, utils.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to commit transaction: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, page)
}

// DeletePage deletes a page by ID
func DeletePage(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	strId := c.Param("id")

	numId, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID provided : " + err.Error(),
		})
		return
	}

	var page models.Page
	if result := db.First(&page, numId); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.HTTPError{
			Code:    http.StatusNotFound,
			Message: "Page could not be found",
		})
		return
	}

	tx := db.Begin()

	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to start transaction",
		})
		return
	}

	if err := db.Delete(&page).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, utils.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete page: " + err.Error(),
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, utils.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to commit transaction: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.MessageResponse{
		Message: "Page successfully deleted",
	})
}
