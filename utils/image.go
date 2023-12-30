package util

import (
	"fmt"
	"time"
	"strings"
	"math/rand"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	// "github.com/jackc/pgx/v5/pgtype"
	"github.com/google/uuid"
)

func GenerateUniqueFilename(originalFilename string) string {
	timestamp := time.Now().Format("20060102150405")
	randomUUID := uuid.New().String()
	randomComponent := rand.Intn(1000000)

	// return fmt.Sprintf("%s_%s_%d_%s", timestamp, randomUUID, randomComponent, originalFilename)

	filename := fmt.Sprintf("%s_%s_%d_%s", timestamp, randomUUID, randomComponent, originalFilename)
	filename = strings.Replace(filename, " ", "-", -1)

	return filename
}


func CreateS3Client(key, secret string) (*session.Session, *s3.S3, error) {
	awsConfig := aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
	}

	// Create a new AWS session.
	sess, err := session.NewSession(&awsConfig)
	if err != nil {
		return nil, nil, err
	}

	svc := s3.New(sess)
	return sess, svc, nil
}


// func (server *Server) uploadVideoToS3(w http.ResponseWriter, r *http.Request) {

// 	if r.Method != http.MethodPost {
// 		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only POST requests are allowed")
// 		return
// 	}
// 	ctx := r.Context()

// 	file, header, err := r.FormFile("image_url")
// 	if err != nil {
// 		fmt.Println(err)
// 		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Failed to get file from request")
// 		return
// 	}

// 	defer file.Close()

// 	filePath := header.Filename             // Replace this with the actual path to your video file.
// 	bucketName := server.config.BUCKET_NAME // Replace this with the name of your S3 bucket.
// 	uniqueFilename := GenerateUniqueFilename(filePath)
// 	objectKey := uniqueFilename

// 	ext := filepath.Ext(filePath)

// 	sess, svc, err := createS3Client(server.config.AWS_KEY, server.config.AWS_SECRET)
// 	if err != nil {
// 		ErrorResponse(w, http.StatusMethodNotAllowed, "Failed to stream video content")
// 		return
// 	}
// 	defer sess.Config.Credentials.Expire()

// 	// Prepare the input parameters for the S3 uploa
// 	input := &s3.PutObjectInput{
// 		Bucket:        aws.String(bucketName),
// 		Key:           aws.String(objectKey),
// 		Body:          file,
// 		ContentLength: aws.Int64(header.Size),
// 		ContentType:   aws.String(mime.TypeByExtension(ext)), // Replace with the appropriate content type for your video file.
// 		// ACL:           aws.String("public-read"),             // Change this as needed, depending on your security requirements.
// 	}

// 	// Upload the video to S3.
// 	_, err = svc.PutObject(input)

// 	if err != nil {
// 		slog.Info("=====", err)
// 		util.ErrorResponse(w, http.StatusMethodNotAllowed, "Failed to upload video content")
// 		return
		
// 	}

// 	arg := db.CreateImageParams{
// 		ImageUrl:       objectKey,
// 	}

// 	_, err = server.store.CreateImage(ctx, arg)
// 	if err != nil {
// 		fmt.Println(err)
// 		ErrorResponse(w, http.StatusMethodNotAllowed, "Failed to insert video name into PostgreSQL")
// 		return
		
// 	}

// 	response := struct {
// 		Status  bool      `json:"status"`
// 		Message string    `json:"message"`
// 	}{
// 		Status:  true,
// 		Message: "Video uploaded successfully to S3",
// 	}

// 	// Respond with a JSON response indicating success.
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)

// }