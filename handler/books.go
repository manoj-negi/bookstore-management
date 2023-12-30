package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/gorilla/mux"
	db "github.com/vod/db/sqlc"
	util "github.com/vod/utils"
)


type Book struct {
	ID              int32            `json:"id"`
	Title           string           `json:"title"`
	AuthorID        int32            `json:"author_id"`
	PublicationDate pgtype.Date      `json:"publication_date"`
	Price           int32            `json:"price"`
	StockQuantity   int32            `json:"stock_quantity"`
	IsDeleted       pgtype.Bool      `json:"is_deleted"`
	CreatedAt       pgtype.Timestamp `json:"created_at"`
	UpdatedAt       pgtype.Timestamp `json:"updated_at"`
	Bestseller      pgtype.Bool      `json:"bestseller"`
}
var emptyDate pgtype.Date

func (server *Server) handlerCreateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only POST requests are allowed")
		return
	}
	ctx := r.Context()

	book := Book{}
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

	arg := db.CreateBookParams{
		Title:  book.Title,
		AuthorID: book.AuthorID,
		PublicationDate: book.PublicationDate,
		Price:  book.Price,
		StockQuantity: book.StockQuantity,
		Bestseller: book.Bestseller,
		IsDeleted: book.IsDeleted,
	}

	bookInfo, err := server.store.CreateBook(ctx, arg)
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
		Data    []db.Book `json:"data"`
	}{
		Status:  true,
		Message: "Book created successfully",
		Data:    []db.Book{bookInfo},
	}

	json.NewEncoder(w).Encode(response)
}

func (server *Server) handlerGetBookById(w http.ResponseWriter, r *http.Request) {
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
	bookInfo, err:= server.store.GetBook(ctx, int32(id))
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch book",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool      `json:"status"`
		Message string    `json:"message"`
		Data    []db.Book `json:"data"`
	}{
		Status:  true,
		Message: "Book retrieved successfully",
		Data:    []db.Book{bookInfo},
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

func (server *Server) handlerGetAllBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only GET requests are allowed")
		return
	}
	ctx := r.Context()

	bookInfo, err := server.store.GetAllBooks(ctx)
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch book",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool      `json:"status"`
		Message string    `json:"message"`
		Data    []db.Book `json:"data"`
	}{
		Status:  true,
		Message: "Book retrieved successfully",
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

func (server *Server) handlerUpdateBook(w http.ResponseWriter, r *http.Request) {
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

    book := Book{}
    err = json.NewDecoder(r.Body).Decode(&book)

    if err != nil {
        util.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON request")
        return
    }

    arg := db.UpdateBookParams{
        ID: int32(id),
    }

    if book.Title != "" {
        arg.SetTitle = true
        arg.Title = book.Title
    }

    if book.AuthorID != 0 {
        arg.SetAuthorID = true
        arg.AuthorID = book.AuthorID
    }

   

	if book.PublicationDate != emptyDate{
		arg.SetPublicationDate = true
        arg.PublicationDate = book.PublicationDate
	} 

    if book.Price != 0 {
        arg.SetPrice = true
        arg.Price = book.Price
    }

    if book.StockQuantity != 0 {
        arg.SetStockQuantity = true
        arg.StockQuantity = book.StockQuantity
    }

    if book.IsDeleted.Valid && book.IsDeleted.Bool {
        arg.SetIsDeleted = true
        arg.IsDeleted = book.IsDeleted
    }

	if book.Bestseller.Valid && book.Bestseller.Bool {
        arg.SetBestseller = true
        arg.Bestseller = book.Bestseller
    }

    bookInfo, err := server.store.UpdateBook(ctx, arg)
    if err != nil {
        fmt.Println("error-------------", err)
        util.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch book")
        return
    }

    response := struct {
        Status  bool     `json:"status"`
        Message string   `json:"message"`
        Data    []db.Book `json:"data"`
    }{
        Status:  true,
        Message: "Book updated successfully",
        Data:    []db.Book{bookInfo},
    }

    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}

func (server *Server) handlerDeleteBook(w http.ResponseWriter, r *http.Request) {
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

	bookInfo, err:= server.store.DeleteBook(ctx, int32(id))
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch book",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []db.Book `json:"data"`
	}{
		Status:  true,
		Message: "book deleted successfully",
		Data:     []db.Book{bookInfo},
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

