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

func(server *Server) CreateRandomRolePermission(t *testing.T) db.RolesPermission{
	arg:= db.CreateRolePermissionParams{
		RoleID: util.RandomInt32(),
		PermissionID: util.RandomInt32(),
		IsDeleted: util.PgtypeBool(),
	}

	rolepermission, err := server.store.CreateRolePermission(context.Background(), arg)
	require.NoError(t,err)
	require.NotEmpty(t,rolepermission)
	require.Equal(t, arg.RoleID, rolepermission.RoleID)
	require.Equal(t, arg.PermissionID, rolepermission.PermissionID)
	require.Equal(t, arg.IsDeleted, rolepermission.IsDeleted)

	require.NotZero(t,rolepermission.ID)
	require.NotZero(t,rolepermission.CreatedAt)
	require.NotZero(t,rolepermission.UpdatedAt)

	return rolepermission
}

func(server *Server) TestCreateRolePermission(t *testing.T) {
	server.CreateRandomRolePermission(t)
}

func(server * Server) TestGetRolePermission(t *testing.T){
	rolepermission1 := server.CreateRandomRolePermission(t)

	rolepermission2, err := server.store.GetRolePermission(context.Background(), rolepermission1.ID)
	require.NoError(t,err)
	require.NotEmpty(t,rolepermission2)
	require.Equal(t, rolepermission1.RoleID, rolepermission2.RoleID)
	require.Equal(t, rolepermission1.PermissionID, rolepermission2.PermissionID)
	require.Equal(t, rolepermission1.IsDeleted, rolepermission2.IsDeleted)

	createdTime1 := rolepermission1.CreatedAt.Time
	createdTime2 := rolepermission2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)

	updatedTime1 := rolepermission1.UpdatedAt.Time
	updatedTime2 := rolepermission2.UpdatedAt.Time
	require.WithinDuration(t, updatedTime1, updatedTime2, time.Second)
}

func(server *Server) TestUpdateRolePermission(t *testing.T){
	rolepermission1 := server.CreateRandomRolePermission(t)

	arg := db.UpdateRolePermissionParams{
		ID: rolepermission1.ID,
		RoleID: util.RandomInt32(),
		PermissionID: util.RandomInt32(),
		IsDeleted: util.PgtypeBool(),
	}

	rolepermission2,err := server.store.UpdateRolePermission(context.Background(),arg)
	require.NoError(t,err)
	require.NotEmpty(t,rolepermission2)
	require.Equal(t, arg.ID, rolepermission2.ID)
	require.Equal(t, arg.RoleID, rolepermission2.RoleID)
	require.Equal(t, arg.PermissionID, rolepermission2.PermissionID)
	require.Equal(t, arg.IsDeleted, rolepermission2.IsDeleted)

	createdTime1 := rolepermission1.CreatedAt.Time
	createdTime2 := rolepermission2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)

	updatedTime1 := rolepermission1.UpdatedAt.Time
	updatedTime2 := rolepermission2.UpdatedAt.Time
	require.WithinDuration(t, updatedTime1, updatedTime2, time.Second)

}

func(server *Server) TestDeleteRolePermission(t *testing.T){

	rolepermission1 := server.CreateRandomRolePermission(t)

	_,err := server.store.DeleteRolePermission(context.Background(),rolepermission1.ID)
	require.NoError(t,err)

	rolepermission2,err := server.store.GetRolePermission(context.Background(),rolepermission1.ID)
	require.NoError(t,err)
	require.EqualError(t,err,sql.ErrNoRows.Error())
	require.Empty(t,rolepermission2)

}

func(server *Server) TestGetAllRolePermission(t *testing.T){
	for i:=0; i<10; i++{
		server.CreateRandomRolePermission(t)
	}


	rolepermission2,err := server.store.GetAllRolePermissions(context.Background())
	require.NoError(t,err)
	require.Len(t,rolepermission2,5)

	for _,rolepermission := range rolepermission2{
		require.NotEmpty(t,rolepermission)
	}
}

