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

type Offer struct {
	ID                 int32            `json:"id"`
	BookID             int32            `json:"book_id" validate:"required"`
	DiscountPercentage pgtype.Text      `json:"discount_percentage" validate:"required"`
	StartDate          pgtype.Date      `json:"start_date" validate:"required"`
	EndDate            pgtype.Date      `json:"end_date" validate:"required"`
	IsDeleted          pgtype.Bool      `json:"is_deleted"`
	CreatedAt          pgtype.Timestamp `json:"created_at"`
	UpdatedAt          pgtype.Timestamp `json:"updated_at"`
}

func (server *Server) handlerCreateOffer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only POST requests are allowed")
		return
	}
	ctx := r.Context()

	offer := Offer{}
	err := json.NewDecoder(r.Body).Decode(&offer)

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
	err = validate.Struct(offer)
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

	arg := db.CreateOfferParams{
		BookID:   offer.BookID,
		DiscountPercentage: offer.DiscountPercentage,
		StartDate: offer.StartDate,
		EndDate:  offer.EndDate,
		IsDeleted: offer.IsDeleted,
	}

	offerInfo, err := server.store.CreateOffer(ctx, arg)
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
		Data    []db.Offer `json:"data"`
	}{
		Status:  true,
		Message: "Offer created successfully",
		Data:    []db.Offer{offerInfo},
	}

	json.NewEncoder(w).Encode(response)
}

func (server *Server) handlerGetOfferById(w http.ResponseWriter, r *http.Request) {
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
	offerInfo, err:= server.store.GetOffer(ctx, int32(id))
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch offer",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool      `json:"status"`
		Message string    `json:"message"`
		Data    []db.Offer `json:"data"`
	}{
		Status:  true,
		Message: "Offer retrieved successfully",
		Data:    []db.Offer{offerInfo},
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

func (server *Server) handlerGetAllOffer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only GET requests are allowed")
		return
	}
	ctx := r.Context()

	offerInfo, err := server.store.GetAllOffers(ctx)
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch offer",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool      `json:"status"`
		Message string    `json:"message"`
		Data    []db.Offer `json:"data"`
	}{
		Status:  true,
		Message: "Offer retrieved successfully",
		Data:    offerInfo,
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

func (server *Server) handlerUpdateOffer(w http.ResponseWriter, r *http.Request) {
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

    offer := Offer{}
    err = json.NewDecoder(r.Body).Decode(&offer)

    if err != nil {
        util.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON request")
        return
    }

    arg := db.UpdateOfferParams{
        ID: int32(id),
    }

    if offer.BookID != 0 {
        arg.SetBookID = true
        arg.BookID = offer.BookID
    }

    if offer.DiscountPercentage != emptyText {
        arg.SetDiscountPercentage = true
        arg.DiscountPercentage = offer.DiscountPercentage
    }

    if offer.StartDate != emptyDate{
        arg.SetStartDate = true
        arg.StartDate = offer.StartDate
    }

    if offer.EndDate != emptyDate {
        arg.SetEndDate = true
        arg.EndDate = offer.EndDate
    }

    if offer.IsDeleted.Valid && offer.IsDeleted.Bool {
        arg.SetIsDeleted = true
        arg.IsDeleted = offer.IsDeleted
    }

    offerInfo, err := server.store.UpdateOffer(ctx, arg)
    if err != nil {
        util.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch offer")
        return
    }

    response := struct {
        Status  bool   `json:"status"`
        Message string `json:"message"`
        Data    []db.Offer `json:"data"`
    }{
        Status:  true,
        Message: "Offer updated successfully",
        Data:    []db.Offer{offerInfo},
    }

    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}

func (server *Server) handlerDeleteOffer(w http.ResponseWriter, r *http.Request) {
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

	offerInfo, err:= server.store.DeleteOffer(ctx, int32(id))
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch offer",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []db.Offer `json:"data"`
	}{
		Status:  true,
		Message: "offer deleted successfully",
		Data:     []db.Offer{offerInfo},
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

