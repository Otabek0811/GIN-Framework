package handler

import (
	"app/models"
	"app/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create User godoc
// @ID create_user
// @Router /user [POST]
// @Summary Create User
// @Description Create User
// @Tags User
// @Accept json
// @Procedure json
// @Param user body models.CreateUser true "UserRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *handler) CreateUser(c *gin.Context) {
	var createuser models.CreateUser

	err := c.ShouldBindJSON(&createuser)
	if err != nil {
		h.handlerResponse(c, "error user should bind json", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.strg.User().CreateUser(&createuser)
	if err != nil {
		h.handlerResponse(c, "storage.user.create", http.StatusInternalServerError, err.Error())
		return
	}
	resp, err := h.strg.User().GetUserByID(&models.UserPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "error while User response by id:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}
	h.handlerResponse(c, "create User response", http.StatusOK, resp)

}

// GetByID User godoc
// @ID get_by_id_user
// @Router /user/{id} [GET]
// @Summary Get By ID User
// @Description Get By ID User
// @Tags User
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetUserByID(c *gin.Context) {
	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "error while give user id: invalid uuid ", http.StatusBadRequest, nil)
		return
	}

	resp, err := h.strg.User().GetUserByID(&models.UserPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.user.getById:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	h.handlerResponse(c, "get by id user resposne", http.StatusOK, resp)

}

// GetList User godoc
// @ID get_list_user
// @Router /user [GET]
// @Summary Get List User
// @Description Get List User
// @Tags User
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListUser(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list user offset", http.StatusBadRequest, "invalid offset")
		return
	}
	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list user limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.User().GetListUser(&models.UserGetListRequest{
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
	})

	if err != nil {
		h.handlerResponse(c, "error while storage user get list:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	h.handlerResponse(c, "Get List User Response", http.StatusOK, resp)
}

// Update User godoc
// @ID update_user
// @Router /user [PUT]
// @Summary Update User
// @Description Update User
// @Tags User
// @Accept json
// @Procedure json
// @Param update body models.UpdateUser true "UserRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *handler) UpdateUser(c *gin.Context) {

	var upUser models.UpdateUser

	err := c.ShouldBindJSON(&upUser)

	id, err := h.strg.User().UpdateUser(&upUser)
	if err != nil {
		h.handlerResponse(c, "error while  user upadate:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	resp, err := h.strg.User().GetUserByID(&models.UserPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "error while Update-Userresponse by id:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	h.handlerResponse(c, "Update User Response", http.StatusOK, resp)

}

// Delete User godoc
// @ID delete_user
// @Router /user/{id} [DELETE]
// @Summary DELETE User
// @Description Delete User
// @Tags User
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteUser(c *gin.Context) {
	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "error while delete user--> give user id: invalid uuid", http.StatusBadRequest, nil)
		return
	}

	err := h.strg.User().DeleteUser(&models.UserPrimaryKey{Id: id})

	if err != nil {
		h.handlerResponse(c, "error while delete user:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}
	h.handlerResponse(c, "Delete User Response", http.StatusOK, nil)

}
