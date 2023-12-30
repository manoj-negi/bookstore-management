package handler

import (
	"encoding/json"
	"net/http"
	db "github.com/vod/db/sqlc"
	util "github.com/vod/utils"
)

type HomeData struct {
    Banners    []db.Banner   `json:"banners"`
    Categories  []db.Category `json:"categories"`
	Offers      []db.Offer     `json:"offers"`
	BestSeller   []db.Book  `json:"bestseller"`
}

func (server *Server) handlerGetHome(w http.ResponseWriter, r *http.Request) {
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

	categoryInfo, err := server.store.GetAllCategories(ctx)
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch category",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

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

	bookInfo, err := server.store.GetBestseller(ctx)
	if err != nil {
		jsonResponse := JsonResponse{
			Status:     false,
			Message:    "Failed to fetch book",
			StatusCode: http.StatusInternalServerError,
		}
		util.WriteJSONResponse(w, http.StatusInternalServerError, jsonResponse)
		return
	}

	homeData := HomeData{
		Banners:    bannerInfo,
		Categories: categoryInfo,
		Offers:     offerInfo,
		BestSeller:      bookInfo,
	}
	

	response := struct {
		Status  bool     `json:"status"`
		Message string   `json:"message"`
		Data    HomeData `json:"data"`
	}{
		Status:  true,
		Message: "Data retrieved successfully",
		Data:    homeData,
	}
	
	

	// response := struct {
	// 	Status  bool      `json:"status"`
	// 	Message string    `json:"message"`
	// 	Data    []db.Banner `json:"data"`
	// }{
	// 	Status:  true,
	// 	Message: "banner retrieved successfully",
	// 	Data:    bannerInfo,
	// }

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



