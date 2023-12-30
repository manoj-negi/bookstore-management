package main

import (
	"context"

	"log/slog"
	"net/http"
	conn "github.com/vod/config/database"
	db "github.com/vod/db/sqlc"
	"github.com/vod/handler"
	util "github.com/vod/utils"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}



func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		slog.Info("cannot load config", err)
	}

	dbConn, err := conn.NewPostgres(context.Background(), config.DB_URI)
	if err != nil {
		slog.Info("cannot connect to db", err)
	}
	store := db.NewStore(dbConn.DB) //db.New(dbConn.DB)

	server, err := handler.NewServer(store, config)
	if err != nil {
		slog.Info("cannot create server")
	}

	handler := corsMiddleware(server.GetRouter())
	
	// err = server.Start(":8080",handler)
	// if err != nil {
	// 	slog.Info("=======", err)
	// }

	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		slog.Info("Error starting server:", err)
	} else {
		slog.Info("Server is running on :8080")
	}

}
