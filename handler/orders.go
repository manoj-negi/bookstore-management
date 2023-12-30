package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"fmt"
	"database/sql/driver"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/gorilla/mux"
	db "github.com/vod/db/sqlc"
	util "github.com/vod/utils"
)


type StatusEnum string

const (
	StatusEnumPending   StatusEnum = "Pending"
	StatusEnumInProcess StatusEnum = "In-Process"
	StatusEnumCompleted StatusEnum = "Completed"
)

func (e *StatusEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = StatusEnum(s)
	case string:
		*e = StatusEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for StatusEnum: %T", src)
	}
	return nil
}

type NullStatusEnum struct {
	StatusEnum StatusEnum `json:"status_enum"`
	Valid      bool       `json:"valid"` // Valid is true if StatusEnum is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullStatusEnum) Scan(value interface{}) error {
	if value == nil {
		ns.StatusEnum, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.StatusEnum.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullStatusEnum) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.StatusEnum), nil
}

type Order struct {
	ID         int32            `json:"id"`
	BookID     int32            `json:"book_id" validate:"required"`
	UserID     int32            `json:"user_id" validate:"required"`
	OrderNo    pgtype.Text      `json:"order_no" validate:"required"`
	Quantity   int32            `json:"quantity" validate:"required"`
	TotalPrice int32            `json:"total_price" validate:"required"`
	Status     StatusEnum       `json:"status" validate:"required"`
	IsDeleted  pgtype.Bool      `json:"is_deleted"`
	CreatedAt  pgtype.Timestamp `json:"created_at"`
	UpdatedAt  pgtype.Timestamp `json:"updated_at"`
}

func (server *Server) handlerCreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only POST requests are allowed")
		return
	}
	ctx := r.Context()

	order := Order{}
	err := json.NewDecoder(r.Body).Decode(&order)

	if err != nil {
		fmt.Println("------err1------",err)
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "invalid JSON request",
			StatusCode: http.StatusNotAcceptable,
		}
		
		util.WriteJSONResponse(w, http.StatusNotAcceptable, jsonResponse)
		return
	}

	validate := validator.New()
	err = validate.Struct(order)
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

	arg := db.CreateOrderParams{
		BookID:  order.BookID,
		UserID:  order.UserID,
		OrderNo: order.OrderNo,
		Quantity: order.Quantity,
		TotalPrice: order.TotalPrice,
		Status:   db.StatusEnum(order.Status),
		IsDeleted: order.IsDeleted,
	}

	orderInfo, err := server.store.CreateOrder(ctx, arg)
	if err != nil {
		fmt.Println("------err1------",err)
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
		Data    []db.Order `json:"data"`
	}{
		Status:  true,
		Message: "order created successfully",
		Data:    []db.Order{orderInfo},
	}

	json.NewEncoder(w).Encode(response)
}

func (server *Server) handlerGetOrderById(w http.ResponseWriter, r *http.Request) {
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
	orderInfo, err:= server.store.GetOrder(ctx, int32(id))
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch order",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool      `json:"status"`
		Message string    `json:"message"`
		Data    []db.Order `json:"data"`
	}{
		Status:  true,
		Message: "order retrieved successfully",
		Data:    []db.Order{orderInfo},
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

func (server *Server) handlerGetAllOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only GET requests are allowed")
		return
	}
	ctx := r.Context()

	orderInfo, err := server.store.GetAllOrders(ctx)
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch order",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool      `json:"status"`
		Message string    `json:"message"`
		Data    []db.Order `json:"data"`
	}{
		Status:  true,
		Message: "order retrieved successfully",
		Data:    orderInfo,
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

func (server *Server) handlerUpdateOrder(w http.ResponseWriter, r *http.Request) {
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

    order := Order{}
    err = json.NewDecoder(r.Body).Decode(&order)

    if err != nil {
        util.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON request")
        return
    }

    arg := db.UpdateOrderParams{
        ID: int32(id),
    }

    if order.BookID != 0 {
        arg.SetBookID = true
        arg.BookID = order.BookID
    }

    if order.UserID != 0 {
        arg.SetUserID = true
        arg.UserID = order.UserID
    }

    if order.OrderNo != emptyText {
        arg.SetOrderNo = true
        arg.OrderNo = order.OrderNo
    }

    if order.Quantity != 0 {
        arg.SetQuantity = true
        arg.Quantity = order.Quantity
    }

    if order.TotalPrice != 0 {
        arg.SetTotalPrice = true
        arg.TotalPrice = order.TotalPrice
    }

    if order.Status != "" {
        arg.SetStatus = true
        arg.Status = db.StatusEnum(order.Status)
    }

    if order.IsDeleted.Valid && order.IsDeleted.Bool {
        arg.SetIsDeleted = true
        arg.IsDeleted = order.IsDeleted
    }

    orderInfo, err := server.store.UpdateOrder(ctx, arg)
    if err != nil {
        fmt.Println("------err1------", err)
        util.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch order")
        return
    }

    response := struct {
        Status  bool   `json:"status"`
        Message string `json:"message"`
        Data    []db.Order `json:"data"`
    }{
        Status:  true,
        Message: "Order updated successfully",
        Data:    []db.Order{orderInfo},
    }

    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}

func (server *Server) handlerDeleteOrder(w http.ResponseWriter, r *http.Request) {
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

	orderInfo, err:= server.store.DeleteOrder(ctx, int32(id))
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch order",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []db.Order `json:"data"`
	}{
		Status:  true,
		Message: "order deleted successfully",
		Data:     []db.Order{orderInfo},
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

