package controllers

import (
	"cms-backend/models"
	"cms-backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMedia(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var media []models.Media

	url := c.Query("url")
	mediaType := c.Query("type")

	query := db
	if url != "" {
		query = query.Where("url ILIKE ?", "%"+url+"%")
	}
	if mediaType != "" {
		query = query.Where("type = ? ", mediaType)
	}

	if err := query.Find(&media).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch media",
		})
		return
	}

	c.JSON(http.StatusOK, media)

}

func GetMediaByID(c *gin.Context) {
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

	var media models.Media
	if err := db.First(&media, numId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve media: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, media)
}

func CreateMedia(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var media models.Media

	if err := c.ShouldBindJSON(&media); err != nil {
		c.JSON(http.StatusBadRequest, utils.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid media data provided : " + err.Error(),
		})
		return
	}

	if err := media.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, utils.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
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

	if err := tx.Create(&media).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, utils.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Unable to add media to the database : " + err.Error(),
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, utils.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to commit transaction : " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, media)

}

func DeleteMedia(c *gin.Context) {
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

	var media models.Media

	if result := db.First(&media, numId); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.HTTPError{
			Code:    http.StatusNotFound,
			Message: "Media could not be found",
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

	if err := tx.Delete(&media).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, utils.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete media: " + err.Error(),
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

	c.JSON(http.StatusAccepted, utils.MessageResponse{
		Message: "Media successfully deleted",
	})

}
