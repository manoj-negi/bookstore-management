package handler

import(
	"context"
	"testing"
	"time"
	db "github.com/vod/db/sqlc"
	"github.com/gorilla/mux"
	util "github.com/vod/utils"
	"database/sql"
	"github.com/stretchr/testify/require"

)

type Server struct {
	config util.Config
	store  db.Querier
	//tokenMaker token.Maker
	router *mux.Router
}

func(server *Server) CreateRandomOffer(t *testing.T) db.Offer{
	arg:= db.CreateOfferParams{
		BookID: util.RandomInt32(),
		DiscountPercentage: util.PgtypeText(),
		StartDate: util.PgtypeDate(),
		EndDate: util.PgtypeDate(),
		IsDeleted: util.PgtypeBool(),
	}

	offer, err := server.store.CreateOffer(context.Background(), arg)
	require.NoError(t,err)
	require.NotEmpty(t,offer)
	require.Equal(t, arg.BookID, offer.BookID)
	require.Equal(t, arg.DiscountPercentage, offer.DiscountPercentage)
	require.Equal(t, arg.StartDate, offer.StartDate)
	require.Equal(t, arg.EndDate, offer.EndDate)
	require.Equal(t, arg.IsDeleted, offer.IsDeleted)

	require.NotZero(t,offer.ID)
	require.NotZero(t,offer.CreatedAt)
	require.NotZero(t,offer.UpdatedAt)

	return offer
}

func(server *Server) TestCreateOffer(t *testing.T) {
	server.CreateRandomOffer(t)
}

func(server * Server) TestGetOffer(t *testing.T){
	offer1 := server.CreateRandomOffer(t)

	offer2, err := server.store.GetOffer(context.Background(), offer1.ID)
	require.NoError(t,err)
	require.NotEmpty(t,offer2)
	require.Equal(t, offer1.BookID, offer2.BookID)
	require.Equal(t, offer1.DiscountPercentage, offer2.DiscountPercentage)
	require.Equal(t, offer1.StartDate, offer2.StartDate)
	require.Equal(t, offer1.EndDate, offer2.EndDate)
	require.Equal(t, offer1.IsDeleted, offer2.IsDeleted)

	createdTime1 := offer1.CreatedAt.Time
	createdTime2 := offer2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)

	updatedTime1 := offer1.UpdatedAt.Time
	updatedTime2 := offer2.UpdatedAt.Time
	require.WithinDuration(t, updatedTime1, updatedTime2, time.Second)
}

func(server *Server) TestUpdateoffer(t *testing.T){
	offer1 := server.CreateRandomOffer(t)

	arg := db.UpdateOfferParams{
		ID: offer1.ID,
		BookID: util.RandomInt32(),
		DiscountPercentage: util.PgtypeText(),
		StartDate: util.PgtypeDate(),
		EndDate: util.PgtypeDate(),
		IsDeleted: util.PgtypeBool(),
	}

	offer2,err := server.store.UpdateOffer(context.Background(),arg)
	require.NoError(t,err)
	require.NotEmpty(t,offer2)
	require.Equal(t, arg.ID, offer2.ID)
	require.Equal(t, arg.BookID, offer2.BookID)
	require.Equal(t, arg.DiscountPercentage, offer2.DiscountPercentage)
	require.Equal(t, arg.StartDate, offer2.StartDate)
	require.Equal(t, arg.EndDate, offer2.EndDate)
	require.Equal(t, arg.IsDeleted, offer2.IsDeleted)

	createdTime1 := offer1.CreatedAt.Time
	createdTime2 := offer2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)

	updatedTime1 := offer1.UpdatedAt.Time
	updatedTime2 := offer2.UpdatedAt.Time
	require.WithinDuration(t, updatedTime1, updatedTime2, time.Second)

}

func(server *Server) TestDeleteOffer(t *testing.T){

	offer1 := server.CreateRandomOffer(t)

	_,err := server.store.DeleteOffer(context.Background(),offer1.ID)
	require.NoError(t,err)

	offer2,err := server.store.GetOffer(context.Background(),offer1.ID)
	require.NoError(t,err)
	require.EqualError(t,err,sql.ErrNoRows.Error())
	require.Empty(t,offer2)

}

func(server *Server) TestGetAllOffer(t *testing.T){
	for i:=0; i<10; i++{
		server.CreateRandomOffer(t)
	}


	offer2,err := server.store.GetAllOffers(context.Background())
	require.NoError(t,err)
	require.Len(t,offer2,5)

	for _,offer := range offer2{
		require.NotEmpty(t,offer)
	}
}

