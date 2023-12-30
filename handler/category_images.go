package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"path/filepath"
	"mime"
	"log/slog"
	"fmt"
	// "github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/gorilla/mux"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	db "github.com/vod/db/sqlc"
	util "github.com/vod/utils"
)

type CategoriesImage struct {
	ID         int32            `json:"id"`
	Image      string      		`json:"image"`
	CategoryID int32            `json:"category_id"`
	IsDeleted  pgtype.Bool      `json:"is_deleted"`
	CreatedAt  pgtype.Timestamp `json:"created_at"`
	UpdatedAt  pgtype.Timestamp `json:"updated_at"`
}


func (server *Server) handlerCreateCategoryImage(w http.ResponseWriter, r *http.Request) {
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


	categoryString := r.FormValue("category_id")

	categoryId, err := strconv.ParseInt(categoryString, 10, 32)
	if err != nil {
		http.Error(w, "Invalid categoryId", http.StatusBadRequest)
		return
	}

	categoryID32 := int32(categoryId)

	IsDeletedString := r.FormValue("is_deleted")

	var isDeleted pgtype.Bool
	if IsDeletedString == "true" {
		isDeleted.Bool = true
	} else {
		isDeleted.Bool = false
	}

	
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
	objectKey := "category/" + uniqueFilename

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
		// ACL:           aws.String("public-read"),             // Change this as needed, depending on your security requirements.
	}

	// Upload the image to S3.
	_, err = svc.PutObject(input)

	if err != nil {
		slog.Info("=====", err)
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Failed to upload the file")
		return
		
	}

	arg := db.CreateCategoryImageParams{
		CategoryID: categoryID32,
		Image:  objectKey,
		IsDeleted:isDeleted,
	}

	categoryInfo, err := server.store.CreateCategoryImage(ctx, arg)
	if err != nil {
		fmt.Println("err:",err)
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
		Data    []db.CategoriesImage `json:"data"`
	}{
		Status:  true,
		Message: "Category Image created successfully",
		Data:    []db.CategoriesImage{categoryInfo},
	}

	json.NewEncoder(w).Encode(response)
}

func (server *Server) handlerGetCategoryImageById(w http.ResponseWriter, r *http.Request) {
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
	categoryInfo, err:= server.store.GetCategoryImage(ctx, int32(id))
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch category image",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool      `json:"status"`
		Message string    `json:"message"`
		Data    []db.CategoriesImage `json:"data"`
	}{
		Status:  true,
		Message: "Category Image retrieved successfully",
		Data:    []db.CategoriesImage{categoryInfo},
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

func (server *Server) handlerGetAllCategoryImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only GET requests are allowed")
		return
	}
	ctx := r.Context()

	categoryInfo, err := server.store.GetAllCategoryImages(ctx)
	if err != nil {
		fmt.Println("errrrrrr",err)
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch category image",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	type ExtendedCategoriesImage struct {
		ID         int32            `json:"id"`
		Image      string      		`json:"image"`
		CategoryID int32            `json:"category_id"`
		IsDeleted  pgtype.Bool      `json:"is_deleted"`
		CreatedAt  pgtype.Timestamp `json:"created_at"`
		UpdatedAt  pgtype.Timestamp `json:"updated_at"`
	}

	extendedCategoryInfo := make([]ExtendedCategoriesImage, len(categoryInfo))

	for i, item := range categoryInfo {
		extendedItem := ExtendedCategoriesImage{
			ID:            item.ID,
			Image:         server.config.BUCKET_URL + item.Image,
			CategoryID:    item.CategoryID,
			IsDeleted:     item.IsDeleted,
			CreatedAt:     item.CreatedAt,
			UpdatedAt:     item.UpdatedAt,
		}
		extendedCategoryInfo[i] = extendedItem
	}


	response := struct {
		Status  bool      `json:"status"`
		Message string    `json:"message"`
		Data    []ExtendedCategoriesImage `json:"data"`
	}{
		Status:  true,
		Message: "Category Image retrieved successfully",
		Data:    extendedCategoryInfo,
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

func (server *Server) handlerUpdateCategoryImage(w http.ResponseWriter, r *http.Request){
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

	category := CategoriesImage{}
	err = json.NewDecoder(r.Body).Decode(&category)

	if err != nil {
		util.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON request")
		return
	}

	arg := db.UpdateCategoryImageParams{
		ID: int32(id),
	}

	if category.CategoryID > 0 {
		arg.SetCategoryID = true
		arg.CategoryID = category.CategoryID
	}

	if category.Image != "" {
		arg.SetImage = true
		arg.Image = category.Image
	}

    if category.IsDeleted.Valid && category.IsDeleted.Bool {
        arg.SetIsDeleted = true
        arg.IsDeleted = category.IsDeleted
    }


	categoryInfo, err := server.store.UpdateCategoryImage(ctx, arg)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch category image")
		return
	}

	response := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []db.CategoriesImage `json:"data"`
	}{
		Status:  true,
		Message: "Category Image updated successfully",
		Data:    []db.CategoriesImage{categoryInfo},
	}

	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (server *Server) handlerDeleteCategoryImage(w http.ResponseWriter, r *http.Request) {
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

	categoryInfo, err:= server.store.DeleteCategoryImage(ctx, int32(id))
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch category image",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	

	response := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    []db.CategoriesImage `json:"data"`
	}{
		Status:  true,
		Message: "category image deleted successfully",
		Data:     []db.CategoriesImage{categoryInfo},
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

