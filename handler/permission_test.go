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

func(server *Server) CreateRandomPermission(t *testing.T) db.Permission{
	arg:= db.CreatePermissionParams{
		Name: util.RandomName(),
		Permission: util.PgtypeText(),
		IsDeleted: util.PgtypeBool(),
	}

	permission, err := server.store.CreatePermission(context.Background(), arg)
	require.NoError(t,err)
	require.NotEmpty(t,permission)
	require.Equal(t, arg.Name, permission.Name)
	require.Equal(t, arg.Permission, permission.Permission)
	require.Equal(t, arg.IsDeleted, permission.IsDeleted)

	require.NotZero(t,permission.ID)
	require.NotZero(t,permission.CreatedAt)
	require.NotZero(t,permission.UpdatedAt)

	return permission
}

func(server *Server) TestCreatePermission(t *testing.T) {
	server.CreateRandomPermission(t)
}

func(server * Server) TestGetPermission(t *testing.T){
	permission1 := server.CreateRandomPermission(t)

	permission2, err := server.store.GetPermission(context.Background(), permission1.ID)
	require.NoError(t,err)
	require.NotEmpty(t,permission2)
	require.Equal(t, permission1.Name, permission2.Name)
	require.Equal(t, permission1.Permission, permission2.Permission)
	require.Equal(t, permission1.IsDeleted, permission2.IsDeleted)

	createdTime1 := permission1.CreatedAt.Time
	createdTime2 := permission2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)

	updatedTime1 := permission1.UpdatedAt.Time
	updatedTime2 := permission2.UpdatedAt.Time
	require.WithinDuration(t, updatedTime1, updatedTime2, time.Second)
}

func(server *Server) TestUpdatePermission(t *testing.T){
	permission1 := server.CreateRandomPermission(t)

	arg := db.UpdatePermissionParams{
		ID: permission1.ID,
		Name: util.RandomName(),
		Permission: util.PgtypeText(),
		IsDeleted: util.PgtypeBool(),
	}

	permission2,err := server.store.UpdatePermission(context.Background(),arg)
	require.NoError(t,err)
	require.NotEmpty(t,permission2)
	require.Equal(t, arg.ID, permission2.ID)
	require.Equal(t, arg.Name, permission2.Name)
	require.Equal(t, arg.Permission, permission2.Permission)
	require.Equal(t, arg.IsDeleted, permission2.IsDeleted)

	createdTime1 := permission1.CreatedAt.Time
	createdTime2 := permission2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)

	updatedTime1 := permission1.UpdatedAt.Time
	updatedTime2 := permission2.UpdatedAt.Time
	require.WithinDuration(t, updatedTime1, updatedTime2, time.Second)

}

func(server *Server) TestDeletePermission(t *testing.T){

	permission1 := server.CreateRandomPermission(t)

	_,err := server.store.DeletePermission(context.Background(),permission1.ID)
	require.NoError(t,err)

	permission2,err := server.store.GetPermission(context.Background(),permission1.ID)
	require.NoError(t,err)
	require.EqualError(t,err,sql.ErrNoRows.Error())
	require.Empty(t,permission2)

}

func(server *Server) TestGetAllPermission(t *testing.T){
	for i:=0; i<10; i++{
		server.CreateRandomPermission(t)
	}

	permission2,err := server.store.GetAllPermissions(context.Background())
	require.NoError(t,err)
	require.Len(t,permission2,5)

	for _,permission := range permission2{
		require.NotEmpty(t,permission)
	}
}

