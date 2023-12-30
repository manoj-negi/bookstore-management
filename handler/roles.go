package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/gorilla/mux"
	db "github.com/vod/db/sqlc"
	util "github.com/vod/utils"
)

type Role struct {
	ID          int32            `json:"id"`
	Name        string           `json:"name" validate:"required"`
	Description pgtype.Text      `json:"description"`
	IsDeleted   pgtype.Bool      `json:"is_deleted"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
}

var emptyText pgtype.Text

func (server *Server) handlerCreateRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only POST requests are allowed")
		return
	}
	ctx := r.Context()

	role := Role{}
	err := json.NewDecoder(r.Body).Decode(&role)

	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "invalid JSON request",
			StatusCode: http.StatusNotAcceptable,
		}
		
		util.WriteJSONResponse(w, http.StatusNotAcceptable, jsonResponse)
		return
	}

	validate := validator.New()
	err = validate.Struct(role)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err != nil {
				jsonResponse := JsonResponse{
					Status:     false,
					Message:    "Invalid value for " + err.Field(),
					StatusCode: http.StatusNotAcceptable,
				}
				
				json.NewEncoder(w).Encode(jsonResponse)
				return

			}
		}
	}

	arg := db.CreateRoleParams{
		Name:    role.Name,
		Description: role.Description,
		IsDeleted: role.IsDeleted,
	}

	roleInfo, err := server.store.CreateRole(ctx, arg)
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "invalid JSON request1",
			StatusCode: http.StatusNotAcceptable,
		}
		util.WriteJSONResponse(w, http.StatusNotAcceptable, jsonResponse)
		return
	}
	

	response := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []db.Role `json:"data"`
	}{
		Status:  true,
		Message: "Role created successfully",
		Data:    []db.Role{roleInfo},
	}

	json.NewEncoder(w).Encode(response)
}

func (server *Server) handlerGetRoleById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only GET requests are allowed")
		return
	}
	ctx := r.Context()
	vars := mux.Vars(r)
	idParam, ok := vars["id"]
	if !ok {
		util.ErrorResponse(w, http.StatusBadRequest, "Missing 'id' URL parameter")
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid 'id' URL parameter")
		return
	}
	roleInfo, err:= server.store.GetRole(ctx, int32(id))
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch role",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool      `json:"status"`
		Message string    `json:"message"`
		Data    []db.Role `json:"data"`
	}{
		Status:  true,
		Message: "Role retrieved successfully",
		Data:    []db.Role{roleInfo},
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to encode response",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}
}

func (server *Server) handlerGetAllRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only GET requests are allowed")
		return
	}
	ctx := r.Context()

	roleInfo, err := server.store.GetAllRoles(ctx)
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch role",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool      `json:"status"`
		Message string    `json:"message"`
		Data    []db.Role `json:"data"`
	}{
		Status:  true,
		Message: "Role retrieved successfully",
		Data:    roleInfo,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to encode response",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}
}

func (server *Server) handlerUpdateRole(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPut {
        util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only PUT requests are allowed")
        return
    }

    ctx := r.Context()

    vars := mux.Vars(r)
    idParam, ok := vars["id"]
    if !ok {
        util.ErrorResponse(w, http.StatusBadRequest, "Missing 'id' URL parameter")
        return
    }

    id, err := strconv.Atoi(idParam)
    if err != nil {
        util.ErrorResponse(w, http.StatusBadRequest, "Invalid 'id' URL parameter")
        return
    }

    role := Role{}
    err = json.NewDecoder(r.Body).Decode(&role)

    if err != nil {
        util.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON request")
        return
    }

    arg := db.UpdateRoleParams{
        ID: int32(id),
    }

    if role.Name != "" {
        arg.SetName = true
        arg.Name = role.Name
    }

    if role.Description != emptyText {
        arg.SetDescription = true
        arg.Description = role.Description
    }

	if role.IsDeleted.Valid && role.IsDeleted.Bool {
        arg.SetIsDeleted = true
        arg.IsDeleted = role.IsDeleted
    }

    roleInfo, err := server.store.UpdateRole(ctx, arg)
    if err != nil {
        util.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch role")
        return
    }

    response := struct {
        Status  bool   `json:"status"`
        Message string `json:"message"`
        Data    []db.Role `json:"data"`
    }{
        Status:  true,
        Message: "Role updated successfully",
        Data:     []db.Role{roleInfo},
    }

    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}

func (server *Server) handlerDeleteRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only DELETE requests are allowed")
		return
	}
	ctx := r.Context()

	vars := mux.Vars(r)
	idParam, ok := vars["id"]
	if !ok {
		util.ErrorResponse(w, http.StatusBadRequest, "Missing 'id' URL parameter")
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid 'id' URL parameter")
		return
	}

	roleInfo, err:= server.store.DeleteRole(ctx, int32(id))
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch role",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []db.Role `json:"data"`
	}{
		Status:  true,
		Message: "role deleted successfully",
		Data:     []db.Role{roleInfo},
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to encode response",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}
}

