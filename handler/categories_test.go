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

func(server *Server) CreateRandomCategory(t *testing.T) db.Category{
	arg:= db.CreateCategoryParams{
		Name: util.RandomName(),
		IsSpecial: util.PgtypeText(),
		IsDeleted: util.PgtypeBool(),
	}

	category, err := server.store.CreateCategory(context.Background(), arg)
	require.NoError(t,err)
	require.NotEmpty(t,category)
	require.Equal(t, arg.Name, category.Name)
	require.Equal(t, arg.IsSpecial, category.IsSpecial)
	require.Equal(t, arg.IsDeleted, category.IsDeleted)

	require.NotZero(t,category.ID)
	require.NotZero(t,category.CreatedAt)
	require.NotZero(t,category.UpdatedAt)

	return category
}

func(server *Server) TestCreateCategory(t *testing.T) {
	server.CreateRandomCategory(t)
}

func(server * Server) TestGetCategory(t *testing.T){
	category1 := server.CreateRandomCategory(t)

	category2, err := server.store.GetCategory(context.Background(), category1.ID)
	require.NoError(t,err)
	require.NotEmpty(t,category2)
	require.Equal(t, category1.ID, category2.ID)
	require.Equal(t, category1.Name, category2.Name)
	require.Equal(t, category1.IsSpecial, category2.IsSpecial)
	require.Equal(t, category1.IsDeleted, category2.IsDeleted)
	createdTime1 := category1.CreatedAt.Time
	createdTime2 := category2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)

	updatedTime1 := category1.UpdatedAt.Time
	updatedTime2 := category2.UpdatedAt.Time
	require.WithinDuration(t, updatedTime1, updatedTime2, time.Second)
}

func(server *Server) TestUpdateCategory(t *testing.T){
	category1 := server.CreateRandomCategory(t)

	arg := db.UpdateCategoryParams{
		ID: category1.ID,
		Name:  util.RandomName(),
		IsSpecial: util.PgtypeText(),
		IsDeleted: util.PgtypeBool(),
	}

	category2,err := server.store.UpdateCategory(context.Background(),arg)
	require.NoError(t,err)
	require.NotEmpty(t,category2)
	require.Equal(t, arg.ID, category2.ID)
	require.Equal(t, arg.Name, category2.Name)
	require.Equal(t, arg.IsSpecial, category2.IsSpecial)
	require.Equal(t, arg.IsDeleted, category2.IsDeleted)

	createdTime1 := category1.CreatedAt.Time
	createdTime2 := category2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)

	updatedTime1 := category1.UpdatedAt.Time
	updatedTime2 := category2.UpdatedAt.Time
	require.WithinDuration(t, updatedTime1, updatedTime2, time.Second)
}


func(server *Server) TestDeleteCategory(t *testing.T){

	category1 := server.CreateRandomCategory(t)

	_,err := server.store.DeleteCategory(context.Background(),category1.ID)
	require.NoError(t,err)

	category2,err := server.store.GetCategory(context.Background(),category1.ID)
	require.NoError(t,err)
	require.EqualError(t,err,sql.ErrNoRows.Error())
	require.Empty(t,category2)

}

func(server *Server) TestGetAllCategory(t *testing.T){
	for i:=0; i<10; i++{
		server.CreateRandomCategory(t)
	}

	category2,err := server.store.GetAllCategories(context.Background())
	require.NoError(t,err)
	require.Len(t,category2,5)

	for _,category := range category2{
		require.NotEmpty(t,category)
	}
}

