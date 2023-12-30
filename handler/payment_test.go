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

func(server *Server) CreateRandomPayment(t *testing.T) db.Payment{
	arg:= db.CreatePaymentParams{
		OrderID: util.RandomInt32(),
		Amount: util.RandomInt32(),
		PaymentStatus: util.GenerateRandomStatus1(),
		IsDeleted: util.PgtypeBool(),
	}

	payment, err := server.store.CreatePayment(context.Background(), arg)
	require.NoError(t,err)
	require.NotEmpty(t,payment)
	require.Equal(t, arg.OrderID, payment.OrderID)
	require.Equal(t, arg.Amount, payment.Amount)
	require.Equal(t, arg.PaymentStatus, payment.PaymentStatus)
	require.Equal(t, arg.IsDeleted, payment.IsDeleted)

	require.NotZero(t,payment.ID)
	require.NotZero(t,payment.CreatedAt)
	require.NotZero(t,payment.UpdatedAt)

	return payment
}

func(server *Server) TestCreatePayment(t *testing.T) {
	server.CreateRandomPayment(t)
}

func(server * Server) TestGetPayment(t *testing.T){
	payment1 := server.CreateRandomPayment(t)

	payment2, err := server.store.GetPayment(context.Background(), payment1.ID)
	require.NoError(t,err)
	require.NotEmpty(t,payment2)
	require.Equal(t, payment1.OrderID, payment2.OrderID)
	require.Equal(t, payment1.Amount, payment2.Amount)
	require.Equal(t, payment1.PaymentStatus, payment2.PaymentStatus)
	require.Equal(t, payment1.IsDeleted, payment2.IsDeleted)

	createdTime1 := payment1.CreatedAt.Time
	createdTime2 := payment2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)

	updatedTime1 := payment1.UpdatedAt.Time
	updatedTime2 := payment2.UpdatedAt.Time
	require.WithinDuration(t, updatedTime1, updatedTime2, time.Second)
}

func(server *Server) TestUpdatePayment(t *testing.T){
	payment1 := server.CreateRandomPayment(t)

	arg := db.UpdatePaymentParams{
		ID: payment1.ID,
		OrderID: util.RandomInt32(),
		Amount: util.RandomInt32(),
		PaymentStatus: util.GenerateRandomStatus1(),
		IsDeleted: util.PgtypeBool(),
	}

	payment2,err := server.store.UpdatePayment(context.Background(),arg)
	require.NoError(t,err)
	require.NotEmpty(t,payment2)
	require.Equal(t, arg.ID, payment2.ID)
	require.Equal(t, arg.OrderID, payment2.OrderID)
	require.Equal(t, arg.Amount, payment2.Amount)
	require.Equal(t, arg.PaymentStatus, payment2.PaymentStatus)
	require.Equal(t, arg.IsDeleted, payment2.IsDeleted)



	createdTime1 := payment1.CreatedAt.Time
	createdTime2 := payment2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)

	updatedTime1 := payment1.UpdatedAt.Time
	updatedTime2 := payment2.UpdatedAt.Time
	require.WithinDuration(t, updatedTime1, updatedTime2, time.Second)

}

func(server *Server) TestDeletePayment(t *testing.T){

	payment1 := server.CreateRandomPayment(t)

	_,err := server.store.DeletePayment(context.Background(),payment1.ID)
	require.NoError(t,err)

	payment2,err := server.store.GetPayment(context.Background(),payment1.ID)
	require.NoError(t,err)
	require.EqualError(t,err,sql.ErrNoRows.Error())
	require.Empty(t,payment2)

}

func(server *Server) TestGetAllPayment(t *testing.T){
	for i:=0; i<10; i++{
		server.CreateRandomPayment(t)
	}

	payment2,err := server.store.GetAllPayments(context.Background())
	require.NoError(t,err)
	require.Len(t,payment2,5)

	for _,payment := range payment2{
		require.NotEmpty(t,payment)
	}
}

