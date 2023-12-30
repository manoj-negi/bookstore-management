package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/vod/db/sqlc"
	util "github.com/vod/utils"
)

type Country struct {
	ID          int32            `json:"id"`
	Iso2        string           `json:"iso2" validate:"required"`
	ShortName   string           `json:"short_name" validate:"required"`
	LongName    string           `json:"long_name" validate:"required"`
	Numcode     pgtype.Text      `json:"numcode" validate:"required"`
	CallingCode string           `json:"calling_code" validate:"required"`
	Cctld       string           `json:"cctld" validate:"required"`
	IsDeleted   pgtype.Bool      `json:"is_deleted"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
}

func (server *Server) handlerCreateCountry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only POST requests are allowed")
		return
	}
	ctx := r.Context()

	country := Country{}
	err := json.NewDecoder(r.Body).Decode(&country)

	if err != nil {
		fmt.Println("------error1------", err)

		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "invalid JSON request",
			StatusCode: http.StatusTeapot,
		}
		util.WriteJSONResponse(w, http.StatusTeapot, jsonResponse)
		return
	}


	validate := validator.New()
	err = validate.Struct(country)
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

	arg := db.CreateCountryParams{
		Iso2:        country.Iso2,
		ShortName:   country.ShortName,
		LongName:    country.LongName,
		Numcode:     country.Numcode,
		CallingCode: country.CallingCode,
		Cctld:       country.Cctld,
		IsDeleted:   country.IsDeleted,
	}

	countryInfo, err := server.store.CreateCountry(ctx, arg)
	if err != nil {
		fmt.Println("------error2------", err)

		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "invalid JSON request",
			StatusCode: http.StatusTeapot,
		}
		util.WriteJSONResponse(w, http.StatusTeapot, jsonResponse)
		return
	}
	

	response := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []db.Country `json:"data"`
	}{
		Status:  true,
		Message: "Country created successfully",
		Data:    []db.Country{countryInfo},
	}

	json.NewEncoder(w).Encode(response)

}

func (server *Server) handlerGetCountryById(w http.ResponseWriter, r *http.Request) {
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

	country, err := server.store.GetCountry(ctx, int32(id))
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch countries",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool         `json:"status"`
		Message string       `json:"message"`
		Data    []db.Country `json:"data"`
	}{
		Status:  true,
		Message: "Country retrieved successfully",
		Data:    []db.Country{country},
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

func (server *Server) handlerGetAllCountry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only GET requests are allowed")
		return
	}
	ctx := r.Context()

	countries, err := server.store.GetAllCountries(ctx)
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch countries",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool         `json:"status"`
		Message string       `json:"message"`
		Data    []db.Country `json:"data"`
	}{
		Status:  true,
		Message: "Countries retrieved successfully",
		Data:    countries,
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

func (server *Server) handlerUpdateCountry(w http.ResponseWriter, r *http.Request) {
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

	countries := db.Country{}
	err = json.NewDecoder(r.Body).Decode(&countries)

	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON request")
		return
	}

	arg := db.UpdateCountryParams{
		ID: int32(id),
	}

	if countries.Iso2 != "" {
		arg.SetIso2 = true
		arg.Iso2 = countries.Iso2
	}

	if countries.ShortName != "" {
		arg.SetShortName = true
		arg.ShortName = countries.ShortName
	}

	if countries.LongName != "" {
		arg.SetLongName = true
		arg.LongName = countries.LongName
	}

	if countries.Numcode != emptyText {
		arg.SetNumcode = true
		arg.Numcode = countries.Numcode
	}

	if countries.CallingCode != "" {
		arg.SetCallingCode = true
		arg.CallingCode = countries.CallingCode
	}

	if countries.Cctld != "" {
		arg.SetCctld = true
		arg.Cctld = countries.Cctld
	}

    if countries.IsDeleted.Valid && countries.IsDeleted.Bool {
        arg.SetIsDeleted = true
        arg.IsDeleted = countries.IsDeleted
    }


	countriesInfo, err := server.store.UpdateCountry(ctx, arg)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch Country")
		return
	}

	response := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []db.Country `json:"data"`
	}{
		Status:  true,
		Message: "Country updated successfully",
		Data:    []db.Country{countriesInfo},
	}

	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (server *Server) handlerDeleteCountry(w http.ResponseWriter, r *http.Request) {
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

	countriesInfo, err:= server.store.DeleteCountry(ctx, int32(id))
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch Country",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []db.Country `json:"data"`
	}{
		Status:  true,
		Message: "Country deleted successfully",
		Data:    []db.Country{countriesInfo},
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
