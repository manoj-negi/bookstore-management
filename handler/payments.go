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


type PaymentStatusEnum string

const (
	PaymentStatusEnumPending   PaymentStatusEnum = "Pending"
	PaymentStatusEnumInProcess PaymentStatusEnum = "In-Process"
	PaymentStatusEnumCompleted PaymentStatusEnum = "Completed"
)

func (e *PaymentStatusEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PaymentStatusEnum(s)
	case string:
		*e = PaymentStatusEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for PaymentStatusEnum: %T", src)
	}
	return nil
}

type NullPaymentStatusEnum struct {
	PaymentStatusEnum PaymentStatusEnum `json:"payment_status_enum"`
	Valid             bool              `json:"valid"` 
}


func (ns *NullPaymentStatusEnum) Scan(value interface{}) error {
	if value == nil {
		ns.PaymentStatusEnum, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PaymentStatusEnum.Scan(value)
}

func (ns NullPaymentStatusEnum) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PaymentStatusEnum), nil
}


type Payment struct {
	ID            int32             `json:"id"`
	OrderID       int32             `json:"order_id" validate:"required"`
	Amount        int32             `json:"amount" validate:"required"`
	PaymentStatus PaymentStatusEnum `json:"payment_status" validate:"required"`
	IsDeleted     pgtype.Bool       `json:"is_deleted"`
	CreatedAt     pgtype.Timestamp  `json:"created_at"`
	UpdatedAt     pgtype.Timestamp  `json:"updated_at"`
}

func (server *Server) handlerCreatePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only POST requests are allowed")
		return
	}
	ctx := r.Context()

	payment := Payment{}
	err := json.NewDecoder(r.Body).Decode(&payment)

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
	err = validate.Struct(payment)
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

	arg := db.CreatePaymentParams{
		OrderID:  payment.OrderID,
		Amount:   payment.Amount,
		PaymentStatus: db.PaymentStatusEnum(payment.PaymentStatus),
		IsDeleted: payment.IsDeleted,
	}

	paymentInfo, err := server.store.CreatePayment(ctx, arg)
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
		Data    []db.Payment `json:"data"`
	}{
		Status:  true,
		Message: "payment created successfully",
		Data:    []db.Payment{paymentInfo},
	}

	json.NewEncoder(w).Encode(response)
}

func (server *Server) handlerGetPaymentById(w http.ResponseWriter, r *http.Request) {
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
	paymentInfo, err:= server.store.GetPayment(ctx, int32(id))
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch payment",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool      `json:"status"`
		Message string    `json:"message"`
		Data    []db.Payment `json:"data"`
	}{
		Status:  true,
		Message: "payment retrieved successfully",
		Data:    []db.Payment{paymentInfo},
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

func (server *Server) handlerGetAllPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only GET requests are allowed")
		return
	}
	ctx := r.Context()

	paymentInfo, err := server.store.GetAllPayments(ctx)
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch payment",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool      `json:"status"`
		Message string    `json:"message"`
		Data    []db.Payment `json:"data"`
	}{
		Status:  true,
		Message: "payment retrieved successfully",
		Data:    paymentInfo,
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

func (server *Server) handlerUpdatePayment(w http.ResponseWriter, r *http.Request) {
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

    payment := Payment{}
    err = json.NewDecoder(r.Body).Decode(&payment)

    if err != nil {
        util.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON request")
        return
    }

    arg := db.UpdatePaymentParams{
        ID: int32(id),
    }

    if payment.OrderID != 0 {
        arg.SetOrderID = true
        arg.OrderID = payment.OrderID
    }

    if payment.Amount != 0 {
        arg.SetAmount = true
        arg.Amount = payment.Amount
    }

    if payment.PaymentStatus != "" {
        arg.SetPaymentStatus = true
        arg.PaymentStatus = db.PaymentStatusEnum(payment.PaymentStatus)
    }

    if payment.IsDeleted.Valid && payment.IsDeleted.Bool {
        arg.SetIsDeleted = true
        arg.IsDeleted = payment.IsDeleted
    }

    paymentInfo, err := server.store.UpdatePayment(ctx, arg)
    if err != nil {
        fmt.Println("------err1------", err)
        util.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch payment")
        return
    }

    response := struct {
        Status  bool   `json:"status"`
        Message string `json:"message"`
        Data    []db.Payment `json:"data"`
    }{
        Status:  true,
        Message: "Payment updated successfully",
        Data:     []db.Payment{paymentInfo},
    }

    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}

func (server *Server) handlerDeletePayment(w http.ResponseWriter, r *http.Request) {
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

	paymentInfo, err:= server.store.DeletePayment(ctx, int32(id))
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch payment",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []db.Payment `json:"data"`
	}{
		Status:  true,
		Message: "payment deleted successfully",
		Data:     []db.Payment{paymentInfo},
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

