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

type BooksCategory struct {
	ID         int32            `json:"id"`
	BookID     int32            `json:"book_id" validate:"required"`
	CategoryID int32            `json:"category_id" validate:"required"`
	IsDeleted  pgtype.Bool      `json:"is_deleted"`
	CreatedAt  pgtype.Timestamp `json:"created_at"`
	UpdatedAt  pgtype.Timestamp `json:"updated_at"`
}
func (server *Server) handlerCreateBookCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only POST requests are allowed")
		return
	}
	ctx := r.Context()

	book := BooksCategory{}
	err := json.NewDecoder(r.Body).Decode(&book)

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
	err = validate.Struct(book)
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

	arg := db.CreateBookCategoryParams{
		BookID: book.BookID,
		CategoryID: book.CategoryID,
		IsDeleted: book.IsDeleted,
	}

	bookInfo, err := server.store.CreateBookCategory(ctx, arg)
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
		Data    []db.BooksCategory `json:"data"`
	}{
		Status:  true,
		Message: "Book Category created successfully",
		Data:    []db.BooksCategory{bookInfo},
	}

	json.NewEncoder(w).Encode(response)
}

func (server *Server) handlerGetBookCategoryById(w http.ResponseWriter, r *http.Request) {
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
	bookInfo, err:= server.store.GetBookCategory(ctx, int32(id))
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch book category",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool      `json:"status"`
		Message string    `json:"message"`
		Data    []db.BooksCategory `json:"data"`
	}{
		Status:  true,
		Message: "Book Category retrieved successfully",
		Data:    []db.BooksCategory{bookInfo},
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

func (server *Server) handlerGetAllBookCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only GET requests are allowed")
		return
	}
	ctx := r.Context()

	bookInfo, err := server.store.GetAllBookCategories(ctx)
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch book category",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool      `json:"status"`
		Message string    `json:"message"`
		Data    []db.BooksCategory `json:"data"`
	}{
		Status:  true,
		Message: "Book Category retrieved successfully",
		Data:    bookInfo,
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

func (server *Server) handlerUpdateBookCategory(w http.ResponseWriter, r *http.Request){
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

	book := BooksCategory{}
	err = json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON request")
		return
	}

	arg := db.UpdateBookCategoryParams{
		ID: int32(id),
	}

	if book.BookID != 0 {
		arg.SetBookID = true
		arg.BookID = book.BookID
	}

	if book.CategoryID != 0 {
		arg.SetCategoryID = true
		arg.CategoryID = book.CategoryID
	}

	if book.IsDeleted.Valid && book.IsDeleted.Bool {
        arg.SetIsDeleted = true
        arg.IsDeleted = book.IsDeleted
    }

	bookInfo, err := server.store.UpdateBookCategory(ctx, arg)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch book category")
		return
	}

	response := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []db.BooksCategory `json:"data"`
	}{
		Status:  true,
		Message: "Book Category updated successfully",
		Data:    []db.BooksCategory{bookInfo},
	}

	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (server *Server) handlerDeleteBookCategory(w http.ResponseWriter, r *http.Request) {
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

	bookInfo, err:= server.store.DeleteBookCategory(ctx, int32(id))
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch book category",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []db.BooksCategory `json:"data"`
	}{
		Status:  true,
		Message: "book category deleted successfully",
		Data:     []db.BooksCategory{bookInfo},
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

