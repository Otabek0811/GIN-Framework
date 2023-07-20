package handler

import (
	"app/models"
	"app/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)
// Create Category godoc
//@ID create_category
//@Router /category [POST]
//@Summary Create Category
//@Description Create Category
//@Tags Category
//@Accept json
//@Procedure json
//@Param cateogry body models.CreateCategory true "CategoryRequest"
//@Success 200 {object} Response{data=string} "Success Request"
//@Response 400 {object} Response{data=string} "Bad Request"
//@Failure 500 {object} Response{data=string} "Server Error"
func (h *handler) CreateCategory(c *gin.Context) {

	var createCategory models.CreateCategory
	err := c.ShouldBindJSON(&createCategory)
	if err != nil {
		h.handlerResponse(c, "error category should bind json", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.strg.Category().CreateCategory(&createCategory)
	if err != nil {
		h.handlerResponse(c, "storage.category.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.strg.Category().GetCategoryByID(&models.CategoryPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.category.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create category response", http.StatusOK, resp)
}
// GetByID Category godoc
// @ID get_by_id_category
// @Router /category/{id} [GET]
// @Summary Get By ID Category
// @Description Get By ID Category
// @Tags Category
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdCategory(c *gin.Context) {

	var id string = c.Param("id")
	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	resp, err := h.strg.Category().GetCategoryByID(&models.CategoryPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.category.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get by id category response", http.StatusOK, resp)
}
// GetList Category godoc
// @ID get_list_category
// @Router /category [GET]
// @Summary Get List Category
// @Description Get List Category
// @Tags Category
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListCategory(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list category offset", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list category limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.Category().GetCategoryList(&models.CategoryGetListRequest{
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.category.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list category response", http.StatusOK, resp)
}
// Update Category godoc
//@ID update_category
//@Router /category [PUT]
//@Summary Update Category
//@Description Update Category
//@Tags Category
//@Accept json
//@Procedure json
//@Param update body models.UpdateCategory true "CategoryRequest"
//@Success 200 {object} Response{data=string} "Success Request"
//@Response 400 {object} Response{data=string} "Bad Request"
//@Failure 500 {object} Response{data=string} "Server Error"
func (h *handler) UpdateCategory(c *gin.Context){
	var updateCategory models.UpdateCategory
	err := c.ShouldBindJSON(&updateCategory)
	if err != nil {
		h.handlerResponse(c, "error updateCategory should bind json", http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.strg.Category().UpdateCategory(&updateCategory)
	if err != nil {
		h.handlerResponse(c, "storage.category.update", http.StatusInternalServerError, err.Error())
		return
	}
	
	resp, err := h.strg.Category().GetCategoryByID(&models.CategoryPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.category.response", http.StatusInternalServerError, err.Error())
		return
	}
	
	h.handlerResponse(c, "update category resposne", http.StatusOK, resp)
}

// Delete Category godoc
// @ID delete_category
// @Router /category/{id} [DELETE]
// @Summary DELETE Category
// @Description Delete Category
// @Tags Category
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteCategory(c *gin.Context){
	var id string = c.Param("id")
	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := h.strg.Category().DeleteCategory(&models.CategoryPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.category.delete", http.StatusInternalServerError, err.Error())
		return
	}
	
	h.handlerResponse(c, "delete category resposne", http.StatusOK,nil)

}