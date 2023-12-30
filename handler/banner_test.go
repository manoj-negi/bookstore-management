package handler

import(
	"context"
	"testing"
	"time"
	"database/sql"
	db "github.com/vod/db/sqlc"
	"github.com/gorilla/mux"
	util "github.com/vod/utils"
	"github.com/stretchr/testify/require"

)

type Server struct {
	config util.Config
	store  db.Querier
	router *mux.Router
}

func(server *Server) CreateRandomBanner(t *testing.T) db.Banner{

	arg:= db.CreateBannerParams{
		Name: util.RandomName(),
		Image:  util.PgtypeText(),
		StartDate: util.PgtypeDate(),
		EndDate: util.PgtypeDate(),
		OfferID: util.RandomInt32(),
		IsDeleted: util.PgtypeBool(),
	}

	banner, err := server.store.CreateBanner(context.Background(), arg)
	require.NoError(t,err)
	require.NotEmpty(t,banner)
	require.Equal(t, arg.Name, banner.Name)
	require.Equal(t, arg.Image, banner.Image)
	require.Equal(t, arg.StartDate, banner.StartDate)
	require.Equal(t, arg.EndDate, banner.EndDate)
	require.Equal(t, arg.OfferID, banner.OfferID)
	require.Equal(t, arg.IsDeleted, banner.IsDeleted)

	require.NotZero(t,banner.ID)
	require.NotZero(t,banner.CreatedAt)
	require.NotZero(t,banner.UpdatedAt)

	return banner
}

func(server *Server) TestCreateBanner(t *testing.T) {
	server.CreateRandomBanner(t)
}

func(server * Server) TestGetBanner(t *testing.T){
	banner1 := server.CreateRandomBanner(t)

	banner2, err := server.store.GetBanner(context.Background(), banner1.ID)
	require.NoError(t,err)
	require.NotEmpty(t,banner2)
	require.Equal(t, banner1.ID, banner2.ID)
	require.Equal(t, banner1.Name, banner2.Name)
	require.Equal(t, banner1.Image, banner2.Image)
	require.Equal(t, banner1.StartDate, banner2.StartDate)
	require.Equal(t, banner1.EndDate, banner2.EndDate)
	require.Equal(t, banner1.OfferID, banner2.OfferID)
	require.Equal(t, banner1.IsDeleted, banner2.IsDeleted)

	createdTime1 := banner1.CreatedAt.Time
	createdTime2 := banner2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)


	updatedBanner1 := banner1.UpdatedAt.Time
	updatedBanner2 := banner2.UpdatedAt.Time
	require.WithinDuration(t, updatedBanner1, updatedBanner2, time.Second)
}

func(server *Server) TestUpdateBanner(t *testing.T){

	banner1 := server.CreateRandomBanner(t)

	arg := db.UpdateBannerParams{
		ID: banner1.ID,
		Name:  util.RandomName(),
		Image: util.PgtypeText(),
		StartDate: util.PgtypeDate(),
		EndDate: util.PgtypeDate(),
		OfferID: util.RandomInt32(),
		IsDeleted: util.PgtypeBool(),
	}

	banner2,err := server.store.UpdateBanner(context.Background(),arg)
	require.NoError(t,err)
	require.NotEmpty(t,banner2)
	require.Equal(t, arg.ID, banner2.ID)
	require.Equal(t, arg.Name, banner2.Name)
	require.Equal(t, arg.Image, banner2.Image)
	require.Equal(t, arg.StartDate, banner2.StartDate)
	require.Equal(t, arg.EndDate, banner2.EndDate)
	require.Equal(t, arg.OfferID, banner2.OfferID)
	require.Equal(t, arg.IsDeleted, banner2.IsDeleted)
	createdTime1 := banner1.CreatedAt.Time
	createdTime2 := banner2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)
	updatedBanner1 := banner1.UpdatedAt.Time
	updatedBanner2 := banner2.UpdatedAt.Time
	require.WithinDuration(t, updatedBanner1, updatedBanner2, time.Second)
}

func(server *Server) TestDeleteBanner(t *testing.T){
	banner1 := server.CreateRandomBanner(t)
	_,err := server.store.DeleteBanner(context.Background(),banner1.ID)
	require.NoError(t,err)

	banner2,err := server.store.DeleteBanner(context.Background(),banner1.ID)
	require.NoError(t,err)
	require.EqualError(t,err,sql.ErrNoRows.Error())
	require.Empty(t,banner2)

}

func(server *Server) TestGetAllBanner(t *testing.T){
	for i:=0; i<10; i++{
		server.CreateRandomBanner(t)
	}

	banner2,err := server.store.GetAllBanners(context.Background())
	require.NoError(t,err)
	require.Len(t,banner2,5)

	for _,banner := range banner2{
		require.NotEmpty(t,banner)
	}
}

