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

func(server *Server) CreateRandomBookCategory(t *testing.T) db.BooksCategory{

	arg:= db.CreateBookCategoryParams{
		BookID: util.RandomInt32(),
		CategoryID: util.RandomInt32(),
		IsDeleted: util.PgtypeBool(),
	}

	bookcategory, err := server.store.CreateBookCategory(context.Background(), arg)
	require.NoError(t,err)
	require.NotEmpty(t,bookcategory)
	require.Equal(t, arg.BookID, bookcategory.BookID)
	require.Equal(t, arg.CategoryID, bookcategory.CategoryID)
	require.Equal(t, arg.IsDeleted, bookcategory.IsDeleted)

	require.NotZero(t,bookcategory.ID)
	require.NotZero(t,bookcategory.CreatedAt)
	require.NotZero(t,bookcategory.UpdatedAt)

	return bookcategory
}

func(server *Server) TestCreateBookCategory(t *testing.T) {
	server.CreateRandomBookCategory(t)
}

func(server * Server) TestGetBookCategory(t *testing.T){
	bookcategory1 := server.CreateRandomBookCategory(t)

	bookcategory2, err := server.store.GetBookCategory(context.Background(), bookcategory1.ID)
	require.NoError(t,err)
	require.NotEmpty(t,bookcategory2)
	require.Equal(t, bookcategory1.ID, bookcategory2.ID)
	require.Equal(t, bookcategory1.BookID, bookcategory2.BookID)
	require.Equal(t, bookcategory1.CategoryID, bookcategory2.CategoryID)
	require.Equal(t, bookcategory1.IsDeleted, bookcategory2.IsDeleted)

	createdTime1 := bookcategory1.CreatedAt.Time
	createdTime2 := bookcategory2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)


	updatedTime1 := bookcategory1.UpdatedAt.Time
	updatedTime2 := bookcategory2.UpdatedAt.Time
	require.WithinDuration(t, updatedTime1, updatedTime2, time.Second)
}

func(server *Server) TestUpdateBookCategory(t *testing.T){

	bookcategory1 := server.CreateRandomBookCategory(t)

	arg := db.UpdateBookCategoryParams{
		ID: bookcategory1.ID,
		BookID: util.RandomInt32(),
		CategoryID: util.RandomInt32(),
		IsDeleted: util.PgtypeBool(),
	}

	bookcategory2,err := server.store.UpdateBookCategory(context.Background(),arg)
	require.NoError(t,err)
	require.NotEmpty(t,bookcategory2)
	require.Equal(t, bookcategory1.ID, bookcategory2.ID)
	require.Equal(t, bookcategory1.BookID, bookcategory2.BookID)
	require.Equal(t, bookcategory1.CategoryID, bookcategory2.CategoryID)
	require.Equal(t, bookcategory1.IsDeleted, bookcategory2.IsDeleted)

	createdTime1 := bookcategory1.CreatedAt.Time
	createdTime2 := bookcategory2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)


	updatedBanner1 := bookcategory1.UpdatedAt.Time
	updatedBanner2 := bookcategory2.UpdatedAt.Time
	require.WithinDuration(t, updatedBanner1, updatedBanner2, time.Second)
}

func(server *Server) TestDeleteBookCategory(t *testing.T){
	bookcategory1 := server.CreateRandomBookCategory(t)
	_,err := server.store.DeleteBookCategory(context.Background(),bookcategory1.ID)
	require.NoError(t,err)

	bookcategory2,err := server.store.GetBookCategory(context.Background(),bookcategory1.ID)
	require.NoError(t,err)
	require.EqualError(t,err,sql.ErrNoRows.Error())
	require.Empty(t,bookcategory2)

}

func(server *Server) TestGetAllBookCategory(t *testing.T){
	for i:=0; i<10; i++{
		server.CreateRandomBookCategory(t)
	}

	bookcategory2,err := server.store.GetAllBookCategories(context.Background())
	require.NoError(t,err)
	require.Len(t,bookcategory2,5)

	for _,bookcategory := range bookcategory2{
		require.NotEmpty(t,bookcategory)
	}
}

