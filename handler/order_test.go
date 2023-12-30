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

func(server *Server) CreateRandomOrder(t *testing.T) db.Order{
	arg:= db.CreateOrderParams{
		BookID: util.RandomInt32(),
		UserID: util.RandomInt32(),
		OrderNo: util.PgtypeText(),
		Quantity: util.RandomInt32(),
		TotalPrice: util.RandomInt32(),
		Status: util.GenerateRandomStatus(),
		IsDeleted: util.PgtypeBool(),
	}

	order, err := server.store.CreateOrder(context.Background(), arg)
	require.NoError(t,err)
	require.NotEmpty(t,order)
	require.Equal(t, arg.BookID, order.BookID)
	require.Equal(t, arg.UserID, order.UserID)
	require.Equal(t, arg.OrderNo, order.OrderNo)
	require.Equal(t, arg.Quantity, order.Quantity)
	require.Equal(t, arg.TotalPrice, order.TotalPrice)
	require.Equal(t, arg.Status, order.Status)
	require.Equal(t, arg.IsDeleted, order.IsDeleted)

	require.NotZero(t,order.ID)
	require.NotZero(t,order.CreatedAt)
	require.NotZero(t,order.UpdatedAt)

	return order
}

func(server *Server) TestCreateOrder(t *testing.T) {
	server.CreateRandomOrder(t)
}

func(server * Server) TestGetOrder(t *testing.T){
	order1 := server.CreateRandomOrder(t)

	order2, err := server.store.GetOrder(context.Background(), order1.ID)
	require.NoError(t,err)
	require.NotEmpty(t,order2)
	require.Equal(t, order1.BookID, order2.BookID)
	require.Equal(t, order1.UserID, order2.UserID)
	require.Equal(t, order1.OrderNo, order2.OrderNo)
	require.Equal(t, order1.Quantity, order2.Quantity)
	require.Equal(t, order1.TotalPrice, order2.TotalPrice)
	require.Equal(t, order1.Status, order2.Status)
	require.Equal(t, order1.IsDeleted, order2.IsDeleted)

	createdTime1 := order1.CreatedAt.Time
	createdTime2 := order2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)

	updatedTime1 := order1.UpdatedAt.Time
	updatedTime2 := order2.UpdatedAt.Time
	require.WithinDuration(t, updatedTime1, updatedTime2, time.Second)
}

func(server *Server) TestUpdateOrder(t *testing.T){
	order1 := server.CreateRandomOrder(t)

	arg := db.UpdateOrderParams{
		ID: order1.ID,
		BookID: util.RandomInt32(),
		UserID: util.RandomInt32(),
		OrderNo: util.PgtypeText(),
		Quantity: util.RandomInt32(),
		TotalPrice: util.RandomInt32(),
		Status: util.GenerateRandomStatus(),
		IsDeleted: util.PgtypeBool(),
	}

	order2,err := server.store.UpdateOrder(context.Background(),arg)
	require.NoError(t,err)
	require.NotEmpty(t,order2)
	require.Equal(t, arg.ID, order2.ID)
	require.Equal(t, arg.BookID, order2.BookID)
	require.Equal(t, arg.UserID, order2.UserID)
	require.Equal(t, arg.OrderNo, order2.OrderNo)
	require.Equal(t, arg.Quantity, order2.Quantity)
	require.Equal(t, arg.TotalPrice, order2.TotalPrice)
	require.Equal(t, arg.Status, order2.Status)
	require.Equal(t, arg.IsDeleted, order2.IsDeleted)


	createdTime1 := order1.CreatedAt.Time
	createdTime2 := order2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)

	updatedTime1 := order1.UpdatedAt.Time
	updatedTime2 := order2.UpdatedAt.Time
	require.WithinDuration(t, updatedTime1, updatedTime2, time.Second)

}

func(server *Server) TestDeleteOrder(t *testing.T){

	order1 := server.CreateRandomOrder(t)

	_,err := server.store.DeleteOrder(context.Background(),order1.ID)
	require.NoError(t,err)

	order2,err := server.store.GetOrder(context.Background(),order1.ID)
	require.NoError(t,err)
	require.EqualError(t,err,sql.ErrNoRows.Error())
	require.Empty(t,order2)

}

func(server *Server) TestGetAllOrder(t *testing.T){
	for i:=0; i<10; i++{
		server.CreateRandomOrder(t)
	}

	order2,err := server.store.GetAllOrders(context.Background())
	require.NoError(t,err)
	require.Len(t,order2,5)

	for _,order := range order2{
		require.NotEmpty(t,order)
	}
}

