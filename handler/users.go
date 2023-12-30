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
	"golang.org/x/crypto/bcrypt"
	db "github.com/vod/db/sqlc"
	util "github.com/vod/utils"
)


type GenderEnum string

const (
	GenderEnumMale   GenderEnum = "Male"
	GenderEnumFemale GenderEnum = "Female"
)

func (e *GenderEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = GenderEnum(s)
	case string:
		*e = GenderEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for GenderEnum: %T", src)
	}
	return nil
}

type NullGenderEnum struct {
	GenderEnum GenderEnum `json:"gender_enum"`
	Valid      bool       `json:"valid"` // Valid is true if GenderEnum is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullGenderEnum) Scan(value interface{}) error {
	if value == nil {
		ns.GenderEnum, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.GenderEnum.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullGenderEnum) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.GenderEnum), nil
}

type User struct {
	ID        int32            `json:"id"`
	FirstName string           `json:"first_name" validate:"required"`
	LastName  string           `json:"last_name" validate:"required"`
	Gender    GenderEnum       `json:"gender" validate:"required"`
	Dob       pgtype.Date      `json:"dob" validate:"required"`
	Address   string           `json:"address" validate:"required"`
	City      string           `json:"city"`
	State     string           `json:"state" validate:"required"`
	CountryID int32            `json:"country_id" validate:"required"`
	MobileNo  string           `json:"mobile_no" validate:"required"`
	Username  string           `json:"username" validate:"required"`
	Email     string           `json:"email" validate:"required"`
	Password  string           `json:"password" validate:"required"`
	RoleID    int32            `json:"role_id" validate:"required"`
	Otp       int32            `json:"otp"`
	IsDeleted pgtype.Bool      `json:"is_deleted"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}


func (server *Server) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only POST requests are allowed")
		return
	}
	ctx := r.Context()

	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)

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

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	otp := util.RandomGenerateOtp()

	validate := validator.New()
	err = validate.Struct(user)
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


	arg := db.CreateUserParams{
		FirstName: user.FirstName,
		LastName: user.LastName,
		Gender: db.GenderEnum(user.Gender),
		Dob: user.Dob,
		Address: user.Address,
		City: user.City,
		State: user.State,
		CountryID: user.CountryID,
		MobileNo: user.MobileNo,
		Username: user.MobileNo,
		Email: user.Email,
		Password: string(hash),
		RoleID: user.RoleID,
		Otp: otp,
		IsDeleted: user.IsDeleted,
	}

	userInfo, err := server.store.CreateUser(ctx, arg)
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
		Data    []db.User `json:"data"`
	}{
		Status:  true,
		Message: "user created successfully",
		Data:    []db.User{userInfo},
	}

	json.NewEncoder(w).Encode(response)
}

func (server *Server) handlerGetUserById(w http.ResponseWriter, r *http.Request) {
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
	userInfo, err:= server.store.GetUser(ctx, int32(id))
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch user",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool      `json:"status"`
		Message string    `json:"message"`
		Data    []db.User `json:"data"`
	}{
		Status:  true,
		Message: "user retrieved successfully",
		Data:    []db.User{userInfo},
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

func (server *Server) handlerGetAllUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only GET requests are allowed")
		return
	}
	ctx := r.Context()

	userInfo, err := server.store.GetAllUsers(ctx)
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch user",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool      `json:"status"`
		Message string    `json:"message"`
		Data    []db.User `json:"data"`
	}{
		Status:  true,
		Message: "user retrieved successfully",
		Data:    userInfo,
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


func (server *Server) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
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

    user := User{}
    err = json.NewDecoder(r.Body).Decode(&user)

    if err != nil {
        util.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON request")
        return
    }

    arg := db.UpdateUserParams{
        ID: int32(id),
    }

    if user.FirstName != "" {
        arg.SetFirstName = true
        arg.FirstName = user.FirstName
    }

    if user.LastName != "" {
        arg.SetLastName = true
        arg.LastName = user.LastName
    }

    if user.Gender != "" {
        arg.SetGender = true
        arg.Gender = db.GenderEnum(user.Gender)
    }

    if user.Dob != emptyDate {
        arg.SetDob = true
        arg.Dob = user.Dob
    }

    if user.Address != "" {
        arg.SetAddress = true
        arg.Address = user.Address
    }

    if user.City != "" {
        arg.SetCity = true
        arg.City = user.City
    }

    if user.State != "" {
        arg.SetState = true
        arg.State = user.State
    }

    if user.CountryID != 0 {
        arg.SetCountryID = true
        arg.CountryID = user.CountryID
    }

    if user.MobileNo != "" {
        arg.SetMobileNo = true
        arg.MobileNo = user.MobileNo
    }

    if user.Username != "" {
        arg.SetUsername = true
        arg.Username = user.Username
    }

    if user.Email != "" {
        arg.SetEmail = true
        arg.Email = user.Email
    }

    if user.Password != "" {
        arg.SetPassword = true
        arg.Password = user.Password
    }

    if user.RoleID != 0 {
        arg.SetRoleID = true
        arg.RoleID = user.RoleID
    }

    if user.IsDeleted.Valid && user.IsDeleted.Bool {
        arg.SetIsDeleted = true
        arg.IsDeleted = user.IsDeleted
    }

    userInfo, err := server.store.UpdateUser(ctx, arg)
    if err != nil {
        fmt.Println("------err1------", err)
        util.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch user")
        return
    }

    response := struct {
        Status  bool   `json:"status"`
        Message string `json:"message"`
        Data    []db.User `json:"data"`
    }{
        Status:  true,
        Message: "User updated successfully",
        Data:     []db.User{userInfo},
    }

    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}

func (server *Server) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
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

	userInfo, err:= server.store.DeleteUser(ctx, int32(id))
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch user",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []db.User `json:"data"`
	}{
		Status:  true,
		Message: "user deleted successfully",
		Data:     []db.User{userInfo},
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

