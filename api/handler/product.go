package handler

import (
	"app/models"
	"app/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)
// Create Product godoc
//@ID create_product
//@Router /product [POST]
//@Summary Create Product
//@Description Create Product
//@Tags Product
//@Accept json
//@Procedure json
//@Param product body models.CreateProduct true "ProductRequest"
//@Success 200 {object} Response{data=string} "Success Request"
//@Response 400 {object} Response{data=string} "Bad Request"
//@Failure 500 {object} Response{data=string} "Server Error"
func (h *handler) CreateProduct(c *gin.Context) {
	var createproduct models.CreateProduct

	err := c.ShouldBindJSON(&createproduct)
	if err != nil {
		h.handlerResponse(c, "error product should bind json", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.strg.Product().CreateProduct(&createproduct)
	if err != nil {
		h.handlerResponse(c, "storage.product.create", http.StatusInternalServerError, err.Error())
		return
	}
	resp, err := h.strg.Product().GetProductByID(&models.ProductPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "error while Product response by id:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}
	h.handlerResponse(c, "create Product response", http.StatusOK, resp)

}

// GetByID Product godoc
// @ID get_by_id_product
// @Router /product/{id} [GET]
// @Summary Get By ID Product
// @Description Get By ID Product
// @Tags Product
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetProductByID(c *gin.Context) {
	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "error while give product id: invalid uuid ", http.StatusBadRequest, nil)
		return
	}

	resp, err := h.strg.Product().GetProductByID(&models.ProductPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.product.getById:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	h.handlerResponse(c, "get by id product resposne", http.StatusOK, resp)

}

// GetList Product godoc
// @ID get_list_product
// @Router /product [GET]
// @Summary Get List Product
// @Description Get List Product
// @Tags Product
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListProduct(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list product offset", http.StatusBadRequest, "invalid offset")
		return
	}
	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list product limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.Product().GetListProduct(&models.ProductGetListRequest{
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
	})

	if err != nil {
		h.handlerResponse(c, "error while storage product get list:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	h.handlerResponse(c, "Get List Product Response", http.StatusOK, resp)
}

// Update Product godoc
//@ID update_product
//@Router /product [PUT]
//@Summary Update Product
//@Description Update Product
//@Tags Product
//@Accept json
//@Procedure json
//@Param update body models.UpdateProduct true "ProductRequest"
//@Success 200 {object} Response{data=string} "Success Request"
//@Response 400 {object} Response{data=string} "Bad Request"
//@Failure 500 {object} Response{data=string} "Server Error"
func (h *handler) UpdateProduct(c *gin.Context) {

	var upProduct models.UpdateProduct

	err := c.ShouldBindJSON(&upProduct)

	id, err := h.strg.Product().UpdateProduct(&upProduct)
	if err != nil {
		h.handlerResponse(c, "error while  product upadate:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	resp, err := h.strg.Product().GetProductByID(&models.ProductPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "error while Update-Product response by id:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	h.handlerResponse(c, "Update Product Response", http.StatusOK, resp)

}

// Delete Product godoc
// @ID delete_product
// @Router /product/{id} [DELETE]
// @Summary DELETE Product
// @Description Delete Product
// @Tags Product
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteProduct(c *gin.Context) {
	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "error while delete product--> give product id: invalid uuid", http.StatusBadRequest, nil)
		return
	}

	err := h.strg.Product().DeleteProduct(&models.ProductPrimaryKey{Id: id})

	if err != nil {
		h.handlerResponse(c, "error while delete product:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}
	h.handlerResponse(c, "Delete Product Response", http.StatusOK, nil)

}


