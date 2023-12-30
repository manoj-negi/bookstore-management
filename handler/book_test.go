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

func(server *Server) CreateRandomBook(t *testing.T) db.Book{

	arg:= db.CreateBookParams{
		Title: util.RandomName(),
		AuthorID: util.RandomInt32(),
		PublicationDate: util.PgtypeDate(),
		Price: util.RandomInt32(),
		StockQuantity: util.RandomInt32(),
		Bestseller: util.PgtypeBool(),
		IsDeleted: util.PgtypeBool(),
	}

	book, err := server.store.CreateBook(context.Background(), arg)
	require.NoError(t,err)
	require.NotEmpty(t,book)
	require.Equal(t, arg.Title, book.Title)
	require.Equal(t, arg.AuthorID, book.AuthorID)
	require.Equal(t, arg.PublicationDate, book.PublicationDate)
	require.Equal(t, arg.Price, book.Price)
	require.Equal(t, arg.StockQuantity, book.StockQuantity)
	require.Equal(t, arg.Bestseller, book.Bestseller)
	require.Equal(t, arg.IsDeleted, book.IsDeleted)

	require.NotZero(t,book.ID)
	require.NotZero(t,book.CreatedAt)
	require.NotZero(t,book.UpdatedAt)

	return book
}

func(server *Server) TestCreateBook(t *testing.T) {
	server.CreateRandomBook(t)
}

func(server * Server) TestGetBook(t *testing.T){
	book1 := server.CreateRandomBook(t)

	book2, err := server.store.GetBook(context.Background(), book1.ID)
	require.NoError(t,err)
	require.NotEmpty(t,book2)
	require.Equal(t, book1.ID, book2.ID)
	require.Equal(t, book1.Title, book2.Title)
	require.Equal(t, book1.AuthorID, book2.AuthorID)
	require.Equal(t, book1.PublicationDate, book2.PublicationDate)
	require.Equal(t, book1.Price, book2.Price)
	require.Equal(t, book1.StockQuantity, book2.StockQuantity)
	require.Equal(t, book1.Bestseller, book2.Bestseller)
	require.Equal(t, book1.IsDeleted, book2.IsDeleted)

	createdTime1 := book1.CreatedAt.Time
	createdTime2 := book2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)


	updatedTime1 := book1.UpdatedAt.Time
	updatedTime2 := book2.UpdatedAt.Time
	require.WithinDuration(t, updatedTime1, updatedTime2, time.Second)
}

func(server *Server) TestUpdateBook(t *testing.T){

	book1 := server.CreateRandomBook(t)

	arg := db.UpdateBookParams{
		ID: book1.ID,
		Title: util.RandomName(),
		AuthorID: util.RandomInt32(),
		PublicationDate: util.PgtypeDate(),
		Price: util.RandomInt32(),
		StockQuantity: util.RandomInt32(),
		Bestseller: util.PgtypeBool(),
		IsDeleted: util.PgtypeBool(),
	}

	book2,err := server.store.UpdateBook(context.Background(),arg)
	require.NoError(t,err)
	require.NotEmpty(t,book2)
	require.Equal(t, book1.ID, book2.ID)
	require.Equal(t, arg.Title, book2.Title)
	require.Equal(t, arg.AuthorID, book2.AuthorID)
	require.Equal(t, arg.PublicationDate, book2.PublicationDate)
	require.Equal(t, arg.Price, book2.Price)
	require.Equal(t, arg.StockQuantity, book2.StockQuantity)
	require.Equal(t, arg.Bestseller, book2.Bestseller)
	require.Equal(t, arg.IsDeleted, book2.IsDeleted)

	createdTime1 := book1.CreatedAt.Time
	createdTime2 := book2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)


	updatedBanner1 := book1.UpdatedAt.Time
	updatedBanner2 := book2.UpdatedAt.Time
	require.WithinDuration(t, updatedBanner1, updatedBanner2, time.Second)
}

func(server *Server) TestDeleteBook(t *testing.T){
	book1 := server.CreateRandomBook(t)

	_,err := server.store.DeleteBook(context.Background(),book1.ID)
	require.NoError(t,err)

	book2,err := server.store.GetBook(context.Background(),book1.ID)
	require.NoError(t,err)
	require.EqualError(t,err,sql.ErrNoRows.Error())
	require.Empty(t,book2)

}

func(server *Server) TestGetAllBook(t *testing.T){
	for i:=0; i<10; i++{
		server.CreateRandomBook(t)
	}

	book2,err := server.store.GetAllBooks(context.Background())
	require.NoError(t,err)
	require.Len(t,book2,5)

	for _,book := range book2{
		require.NotEmpty(t,book)
	}
}

