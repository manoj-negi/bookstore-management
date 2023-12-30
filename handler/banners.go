package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"fmt"
	"path/filepath"
	"time"
	"log/slog"
	"mime"
	// "github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	db "github.com/vod/db/sqlc"
	util "github.com/vod/utils"
)

type Banner struct {
	ID        int32            `json:"id"`
	Name      string           `json:"name" validate:"required"`
	Image     string           `json:"image" validate:"required"`
	StartDate time.Time        `json:"start_date" validate:"required"`
	EndDate   time.Time        `json:"end_date" validate:"required"`
	OfferID   int32            `json:"offer_id" validate:"required"`
	IsDeleted pgtype.Bool      `json:"is_deleted"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

func (server *Server) handlerCreateBanner(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only POST requests are allowed")
        return
    }
    ctx := r.Context()

	err := r.ParseMultipartForm(10 << 20)
            if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }

	nameString := r.FormValue("name")	

	file, header, err := r.FormFile("image")
	if err != nil {
		fmt.Println(err)
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Failed to get file from request")
		return
	}

	defer file.Close()

	filePath := header.Filename             
	bucketName := server.config.BUCKET_NAME 
	uniqueFilename := util.GenerateUniqueFilename(filePath)
	objectKey := "banner/" + uniqueFilename

	ext := filepath.Ext(filePath)

	sess, svc, err := util.CreateS3Client(server.config.AWS_KEY, server.config.AWS_SECRET)
	if err != nil {
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Failed to create the s3 client")
		return
	}
	defer sess.Config.Credentials.Expire()

	input := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(objectKey),
		Body:          file,
		ContentLength: aws.Int64(header.Size),
		ContentType:   aws.String(mime.TypeByExtension(ext)),
		// ACL:           aws.String("public-read"), 
	}

	_, err = svc.PutObject(input)

	if err != nil {
		slog.Info("=====", err)
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Failed to upload the file")
		return
		
	}

	startDateString := r.FormValue("start_date")

	startDate, err := time.Parse("2006-01-02", startDateString)
		if err != nil {
		http.Error(w, "Invalid end Date", http.StatusBadRequest)
		return
		}

	endDateString := r.FormValue("end_date")

	endDate, err := time.Parse("2006-01-02", endDateString)

	if err != nil {
		http.Error(w, "Invalid end Date", http.StatusBadRequest)
		return
	}

	offerIdString := r.FormValue("offer_id")

	offerId, err := strconv.ParseInt(offerIdString, 10, 32)
	if err != nil {
		http.Error(w, "Invalid categoryId", http.StatusBadRequest)
		return
	}

	offerID32 := int32(offerId)

	isDeletedString := r.FormValue("is_deleted")	
	
	var isDeleted pgtype.Bool

	if isDeletedString == "true" {
		isDeleted.Bool = true
	} else {
		isDeleted.Bool = false
	}

    arg := db.CreateBannerParams{
        Name:   nameString,
		Image:  objectKey,
        OfferID: offerID32,
        StartDate: startDate,
        EndDate:  endDate,
        IsDeleted: isDeleted,
    }

    bannerInfo, err := server.store.CreateBanner(ctx, arg)
    if err != nil {
		fmt.Println("------log---1",err)
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
        Data    []db.Banner `json:"data"`
    }{
        Status:  true,
        Message: "Banner created successfully",
        Data:    []db.Banner{bannerInfo},
    }

    json.NewEncoder(w).Encode(response)
}


func (server *Server) handlerGetBannerById(w http.ResponseWriter, r *http.Request) {
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
	bannerInfo, err:= server.store.GetBanner(ctx, int32(id))
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch banner",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool      `json:"status"`
		Message string    `json:"message"`
		Data    []db.Banner `json:"data"`
	}{
		Status:  true,
		Message: "banner retrieved successfully",
		Data:    []db.Banner{bannerInfo},
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

func (server *Server) handlerGetAllBanner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only GET requests are allowed")
		return
	}
	ctx := r.Context()

	bannerInfo, err := server.store.GetAllBanners(ctx)
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch banner",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	type ExtendedBanner struct {
		ID        int32            `json:"id"`
		Name      string           `json:"name"`
		Image     string           `json:"image"`
		StartDate time.Time        `json:"start_date"`
		EndDate   time.Time        `json:"end_date"`
		OfferID   int32            `json:"offer_id"`
		IsDeleted pgtype.Bool      `json:"is_deleted"`
		CreatedAt pgtype.Timestamp `json:"created_at"`
		UpdatedAt pgtype.Timestamp `json:"updated_at"`
	}

	extendedBannerInfo := make([]ExtendedBanner, len(bannerInfo))

	for i, item := range bannerInfo {
		extendedItem := ExtendedBanner{
			ID:            item.ID,
			Name:          item.Name,
			Image:         server.config.BUCKET_URL + item.Image,
			StartDate:     item.StartDate,
			EndDate:       item.EndDate,
			OfferID:       item.OfferID,
			IsDeleted:     item.IsDeleted,
			CreatedAt:     item.CreatedAt,
			UpdatedAt:     item.UpdatedAt,
		}
		extendedBannerInfo[i] = extendedItem
	}

	

	response := struct {
		Status  bool      `json:"status"`
		Message string    `json:"message"`
		Data    []ExtendedBanner `json:"data"`
	}{
		Status:  true,
		Message: "banner retrieved successfully",
		Data:    extendedBannerInfo,
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

func (server *Server) handlerUpdateBanner(w http.ResponseWriter, r *http.Request) {
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


	banner := Banner{}
	err = json.NewDecoder(r.Body).Decode(&banner)

	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON request")
		return
	}

	arg := db.UpdateBannerParams{
		ID: int32(id),
	}

	if banner.Name  != "" {
		arg.SetName = true
		arg.Name = banner.Name
	}

	if banner.Image != "" {
		arg.SetImage = true
		arg.Image = banner.Image
	}

	// if banner.StartDate != emptyDate {
	// 	arg.SetStartDate = true
	// 	arg.StartDate = banner.StartDate
	// }

	if !banner.StartDate.IsZero() {
		arg.SetStartDate = true
		arg.StartDate = banner.StartDate
	}

	// if banner.EndDate != emptyDate {
	// 	arg.SetEndDate = true
	// 	arg.EndDate = banner.EndDate
	// }

	if !banner.EndDate.IsZero() {
		arg.SetStartDate = true
		arg.StartDate = banner.EndDate
	}


	if banner.OfferID != 0 {
		arg.SetOfferID = true
		arg.OfferID = banner.OfferID
	}

	if banner.IsDeleted.Valid && banner.IsDeleted.Bool {
        arg.SetIsDeleted = true
        arg.IsDeleted = banner.IsDeleted
    }

	bannerInfo, err := server.store.UpdateBanner(ctx, arg)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch banner")
		return
	}

	response := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []db.Banner `json:"data"`
	}{
		Status:  true,
		Message: "banner updated successfully",
		Data:    []db.Banner{bannerInfo},
	}

	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (server *Server) handlerDeleteBanner(w http.ResponseWriter, r *http.Request) {
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

	bannerInfo, err:= server.store.DeleteBanner(ctx, int32(id))
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch banner",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []db.Banner `json:"data"`
	}{
		Status:  true,
		Message: "banner deleted successfully",
		Data:     []db.Banner{bannerInfo},
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

