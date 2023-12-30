package handler

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	db "github.com/vod/db/sqlc"
	util "github.com/vod/utils"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config util.Config
	store  db.Store
	//tokenMaker token.Maker
	router *mux.Router
}

func SetContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}


// NewServer creates a new HTTP server and set up routing.
func NewServer(store db.Store, config util.Config) (*Server, error) {

	server := &Server{
		store:  store,
		config: config,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) GetRouter() *mux.Router {
	return server.router
}

func (server *Server) setupRouter() {
	router := mux.NewRouter()
	apiV1Router := router.PathPrefix("/api/v1").Subrouter()
    apiV1Router.Use(SetContentTypeJSON)
	// router.Use(SetContentTypeJSON)
	// router.HandleFunc("/book/insert", server.handlerInsertBook)


	//router.HandleFunc("/upload/video", server.uploadVideoToS3).Methods("POST")
	//router.HandleFunc("/videos", server.listAllVideos).Methods("GET")
	
	apiV1Router.HandleFunc("/role", server.handlerCreateRole).Methods("POST")
	apiV1Router.HandleFunc("/role/{id}", server.handlerGetRoleById).Methods("GET")
	apiV1Router.HandleFunc("/role", server.handlerGetAllRole).Methods("GET")
	apiV1Router.HandleFunc("/role/{id}", server.handlerUpdateRole).Methods("PUT")
	apiV1Router.HandleFunc("/role/{id}", server.handlerDeleteRole).Methods("DELETE")


	apiV1Router.HandleFunc("/permission", server.handlerCreatePermission).Methods("POST")
	apiV1Router.HandleFunc("/permission/{id}", server.handlerGetPermissionById).Methods("GET")
	apiV1Router.HandleFunc("/permission", server.handlerGetAllPermission).Methods("GET")
	apiV1Router.HandleFunc("/permission/{id}", server.handlerUpdatePermission).Methods("PUT")
	apiV1Router.HandleFunc("/permission/{id}", server.handlerDeletePermission).Methods("DELETE")

	apiV1Router.HandleFunc("/rolepermission", server.handlerCreateRolePermission).Methods("POST")
	apiV1Router.HandleFunc("/rolepermission/{id}", server.handlerGetRolePermissionById).Methods("GET")
	apiV1Router.HandleFunc("/rolepermission", server.handlerGetAllRolePermission).Methods("GET")
	apiV1Router.HandleFunc("/rolepermission/{id}", server.handlerUpdateRolePermission).Methods("PUT")
	apiV1Router.HandleFunc("/rolepermission/{id}", server.handlerDeleteRolePermission).Methods("DELETE")

	apiV1Router.HandleFunc("/country", server.handlerCreateCountry).Methods("POST")
	apiV1Router.HandleFunc("/country/{id}", server.handlerGetCountryById).Methods("GET")
	apiV1Router.HandleFunc("/country", server.handlerGetAllCountry).Methods("GET")
	apiV1Router.HandleFunc("/country/{id}", server.handlerUpdateCountry).Methods("PUT")
	apiV1Router.HandleFunc("/country/{id}", server.handlerDeleteCountry).Methods("DELETE")

	apiV1Router.HandleFunc("/author", server.handlerCreateAuthor).Methods("POST")
	apiV1Router.HandleFunc("/author/{id}", server.handlerGetAuthorById).Methods("GET")
	router.HandleFunc("/author", server.handlerGetAllAuthor).Methods("GET")
	apiV1Router.HandleFunc("/author/{id}", server.handlerUpdateAuthor).Methods("PUT")
	apiV1Router.HandleFunc("/author/{id}", server.handlerDeleteAuthor).Methods("DELETE")

	apiV1Router.HandleFunc("/category", server.handlerCreateCategory).Methods("POST")
	apiV1Router.HandleFunc("/category/{id}", server.handlerGetCategoryById).Methods("GET")
	apiV1Router.HandleFunc("/category", server.handlerGetAllCategory).Methods("GET")
	apiV1Router.HandleFunc("/category/{id}", server.handlerUpdateCategory).Methods("PUT")
	apiV1Router.HandleFunc("/category/{id}", server.handlerDeleteCategory).Methods("DELETE")

	apiV1Router.HandleFunc("/categoryimage", server.handlerCreateCategoryImage).Methods("POST")
	apiV1Router.HandleFunc("/categoryimage/{id}", server.handlerGetCategoryImageById).Methods("GET")
	apiV1Router.HandleFunc("/categoryimage", server.handlerGetAllCategoryImage).Methods("GET")
	apiV1Router.HandleFunc("/categoryimage/{id}", server.handlerUpdateCategoryImage).Methods("PUT")
	apiV1Router.HandleFunc("/categoryimage/{id}", server.handlerDeleteCategoryImage).Methods("DELETE")

	apiV1Router.HandleFunc("/book", server.handlerCreateBook).Methods("POST")
	apiV1Router.HandleFunc("/book/{id}", server.handlerGetBookById).Methods("GET")
	apiV1Router.HandleFunc("/book", server.handlerGetAllBook).Methods("GET")
	apiV1Router.HandleFunc("/book/{id}", server.handlerUpdateBook).Methods("PUT")
	apiV1Router.HandleFunc("/book/{id}", server.handlerDeleteBook).Methods("DELETE")


	apiV1Router.HandleFunc("/bookcategory", server.handlerCreateBookCategory).Methods("POST")
	apiV1Router.HandleFunc("/bookcategory/{id}", server.handlerGetBookCategoryById).Methods("GET")
	apiV1Router.HandleFunc("/bookcategory", server.handlerGetAllBookCategory).Methods("GET")
	apiV1Router.HandleFunc("/bookcategory/{id}", server.handlerUpdateBookCategory).Methods("PUT")
	apiV1Router.HandleFunc("/bookcategory/{id}", server.handlerDeleteBookCategory).Methods("DELETE")

	apiV1Router.HandleFunc("/offer", server.handlerCreateOffer).Methods("POST")
	apiV1Router.HandleFunc("/offer/{id}", server.handlerGetOfferById).Methods("GET")
	apiV1Router.HandleFunc("/offer", server.handlerGetAllOffer).Methods("GET")
	apiV1Router.HandleFunc("/offer/{id}", server.handlerUpdateOffer).Methods("PUT")
	apiV1Router.HandleFunc("/offer/{id}", server.handlerDeleteOffer).Methods("DELETE")

	apiV1Router.HandleFunc("/banner", server.handlerCreateBanner).Methods("POST")
	apiV1Router.HandleFunc("/banner/{id}", server.handlerGetBannerById).Methods("GET")
	apiV1Router.HandleFunc("/banner", server.handlerGetAllBanner).Methods("GET")
	apiV1Router.HandleFunc("/banner/{id}", server.handlerUpdateBanner).Methods("PUT")
	apiV1Router.HandleFunc("/banner/{id}", server.handlerDeleteBanner).Methods("DELETE")

	apiV1Router.HandleFunc("/user", server.handlerCreateUser).Methods("POST")
	apiV1Router.HandleFunc("/user/{id}", server.handlerGetUserById).Methods("GET")
	apiV1Router.HandleFunc("/user", server.handlerGetAllUser).Methods("GET")
	apiV1Router.HandleFunc("/user/{id}", server.handlerUpdateUser).Methods("PUT")
	apiV1Router.HandleFunc("/user/{id}", server.handlerDeleteUser).Methods("DELETE")

	apiV1Router.HandleFunc("/order", server.handlerCreateOrder).Methods("POST")
	apiV1Router.HandleFunc("/order/{id}", server.handlerGetOrderById).Methods("GET")
	apiV1Router.HandleFunc("/order", server.handlerGetAllOrder).Methods("GET")
	apiV1Router.HandleFunc("/order/{id}", server.handlerUpdateOrder).Methods("PUT")
	apiV1Router.HandleFunc("/order/{id}", server.handlerDeleteOrder).Methods("DELETE")

	apiV1Router.HandleFunc("/payment", server.handlerCreatePayment).Methods("POST")
	apiV1Router.HandleFunc("/payment/{id}", server.handlerGetPaymentById).Methods("GET")
	apiV1Router.HandleFunc("/payment", server.handlerGetAllPayment).Methods("GET")
	apiV1Router.HandleFunc("/payment/{id}", server.handlerUpdatePayment).Methods("PUT")
	apiV1Router.HandleFunc("/payment/{id}", server.handlerDeletePayment).Methods("DELETE")
  
	apiV1Router.HandleFunc("/home", server.handlerGetHome).Methods("GET")
	// apiV1Router.HandleFunc("/image", server.uploadVideoToS3).Methods("POST")
	server.router = router 
}


// Start runs the HTTP server on a specific address.
func (server *Server) Start(port string) error {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	addr := "0.0.0.0"

	srv := &http.Server{
		Addr: addr + port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      server.router, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			slog.Info(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(ctx)

	slog.Info("shutting down")
	os.Exit(0)
	return nil

}

