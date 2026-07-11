package interfaces

import (
	"fmt"
	"net/http"
	"server/internal/application/category"
	"time"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService category.CategoryService
}

func NewCategoryHandler(categoryService *category.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: *categoryService,
	}
}

func (ch *CategoryHandler) CreateCategory(c *gin.Context) {
	var req CreateCateReq

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		errorResponse(c, 400, "Can not bind request")
		return
	}

	fmt.Println(req)
	username := c.GetHeader("KIT-DEV-USERNAME")

	res, err := ch.categoryService.CreateCategory(c.Request.Context(), username, req.Name, req.Des)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	successResponse(c, http.StatusCreated, "Create Category success", GetCateRes{
		ID:        res.ID,
		Name:      res.Name,
		Desc:      res.Desc,
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
	})
}

func (ch *CategoryHandler) GetCategoryList(c *gin.Context) {
	username := c.GetHeader("KIT-DEV-USERNAME")

	res, err := ch.categoryService.GetCategoryListByID(c.Request.Context(), username)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var data []GetCateRes
	for _, element := range res {
		data = append(data, GetCateRes{
			ID:        element.ID,
			Name:      element.Name,
			Desc:      element.Desc,
			CreatedAt: element.CreatedAt,
			UpdateAt:  element.UpdateAt,
		})
	}
	successResponse(c, http.StatusCreated, "Get list Category success", data)
}
