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



func(server *Server) CreateRandomAuthor(t *testing.T) db.Author{
	arg:= db.CreateAuthorParams{
		Name: util.RandomName(),
		IsDeleted: util.PgtypeBool(),
	}

	author, err := server.store.CreateAuthor(context.Background(), arg)
	require.NoError(t,err)
	require.NotEmpty(t,author)
	require.Equal(t, arg.Name, author.Name)
	require.Equal(t, arg.IsDeleted, author.IsDeleted)

	require.NotZero(t,author.ID)
	require.NotZero(t,author.CreatedAt)
	require.NotZero(t,author.UpdatedAt)

	return author
}

func(server *Server) TestCreateAuthor(t *testing.T) {
	server.CreateRandomAuthor(t)
}

func(server * Server) TestGetAuthor(t *testing.T){
	author1 := server.CreateRandomAuthor(t)

	author2, err := server.store.GetAuthor(context.Background(), author1.ID)
	require.NoError(t,err)
	require.NotEmpty(t,author2)
	require.Equal(t, author1.ID, author2.ID)
	require.Equal(t, author1.Name, author2.Name)
	createdTime1 := author1.CreatedAt.Time
	createdTime2 := author2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)
}

func(server *Server) TestUpdateAuthor(t *testing.T){
	author1 := server.CreateRandomAuthor(t)

	arg := db.UpdateAuthorParams{
		ID: author1.ID,
		Name:  util.RandomName(),
	}

	author2,err := server.store.UpdateAuthor(context.Background(),arg)
	require.NoError(t,err)
	require.NotEmpty(t,author2)
	require.Equal(t, arg.ID, author2.ID)
	require.Equal(t, arg.Name, author2.Name)
	createdTime1 := author1.CreatedAt.Time
	createdTime2 := author2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)

}

func(server *Server) TestDeleteAuthor(t *testing.T){
	author1 := server.CreateRandomAuthor(t)
	_,err := server.store.DeleteAuthor(context.Background(),author1.ID)
	require.NoError(t,err)

	author2,err := server.store.GetAuthor(context.Background(),author1.ID)
	require.NoError(t,err)
	require.EqualError(t,err,sql.ErrNoRows.Error())
	require.Empty(t,author2)

}

func(server *Server) TestGetAllAuthor(t *testing.T){
	for i:=0; i<10; i++{
		server.CreateRandomAuthor(t)
	}

	author2,err := server.store.GetAllAuthors(context.Background())
	require.NoError(t,err)
	require.Len(t,author2,5)

	for _,author := range author2{
		require.NotEmpty(t,author)
	}
}

