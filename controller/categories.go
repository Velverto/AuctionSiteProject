package controller

import (
	"AuctionSiteProject/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//	@Tags Categories Controller
//  @Summary Tìm categories bằng id
//	@Description Tìm categories bằng id, nhận id dưới dạng query, trả về toàn bộ tên categories nếu để trống.
//  @Param id query string true "id of categories, if empty then return all"
//  @Success 200 {object} model.Categories
//	@Failure 500 {body} string "Error message"
//  @Router /categories [GET]
func SearchCategories(c *gin.Context) {
	db := GetDBInstance().Db
	var categories []model.Categories
	id := c.DefaultQuery("id", "all")

	errGetCategories := db.Table("categories").
		Where("categories_id = ? OR 'all' = ?", id, id).
		Select("*").
		Scan(&categories).
		Error
	if errGetCategories != nil {
		log.Println(errGetCategories)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errGetCategories,
			"message": "Error while fetching categories data",
		})
		return
	}
	c.JSON(200, categories)
	return
}

//	@Tags Categories Controller
//	@Summary Tạo Categories mới (Administrator only)
//  @Success 200 {body} string "Success message"
//	@Failure 500 {body} string "Error message"
//  @Router /categories [POST]
func NewCategories(c *gin.Context) {
	var headerInfo model.AuthorizationHeader
	if err := c.ShouldBindHeader(&headerInfo); err != nil {
		c.JSON(200, err)
	}
	var userID string
	var errtoken error
	if userID, errtoken = checkSessionToken(headerInfo.Token); errtoken != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   errtoken,
			"message": "Bad request",
		})
		return
	} else if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token không hợp lệ",
		})
		return
	}
	//check if user are adminitrator
	if check, err := checkAdministrator(userID); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "Error while fetching user data",
		})
		return
	} else if check == false {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Only Administrator can use this API!",
		})
		return
	}

	db := GetDBInstance().Db
	var categories model.Categories
	errJSON := c.BindJSON(&categories)
	if errJSON != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   errJSON,
			"message": "Not a valid JSON!",
		})
		return
	}

	errCreate := db.Table("categories").Create(categories).Error
	if errCreate != nil {
		log.Println(errCreate)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errCreate,
			"message": "Error while creating new Categories!",
		})
		return
	}
	c.JSON(200, "Successfully create new Categories!")
	return
}
