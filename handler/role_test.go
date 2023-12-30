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

func(server *Server) CreateRandomRole(t *testing.T) db.Role{
	arg:= db.CreateRoleParams{
		Name: util.RandomName(),
		Description: util.PgtypeText(),
		IsDeleted: util.PgtypeBool(),
	}

	role, err := server.store.CreateRole(context.Background(), arg)
	require.NoError(t,err)
	require.NotEmpty(t,role)
	require.Equal(t, arg.Name, role.Name)
	require.Equal(t, arg.Description, role.Description)
	require.Equal(t, arg.IsDeleted, role.IsDeleted)

	require.NotZero(t,role.ID)
	require.NotZero(t,role.CreatedAt)
	require.NotZero(t,role.UpdatedAt)

	return role
}

func(server *Server) TestCreateRole(t *testing.T) {
	server.CreateRandomRole(t)
}

func(server * Server) TestGetRole(t *testing.T){
	role1 := server.CreateRandomRole(t)

	role2, err := server.store.GetRole(context.Background(), role1.ID)
	require.NoError(t,err)
	require.NotEmpty(t,role2)
	require.Equal(t, role1.Name, role2.Name)
	require.Equal(t, role1.Description, role2.Description)
	require.Equal(t, role1.IsDeleted, role2.IsDeleted)

	createdTime1 := role1.CreatedAt.Time
	createdTime2 := role2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)

	updatedTime1 := role1.UpdatedAt.Time
	updatedTime2 := role2.UpdatedAt.Time
	require.WithinDuration(t, updatedTime1, updatedTime2, time.Second)
}

func(server *Server) TestUpdateRole(t *testing.T){
	role1 := server.CreateRandomRole(t)

	arg := db.UpdateRoleParams{
		ID: role1.ID,
		Name: util.RandomName(),
		Description: util.PgtypeText(),
		IsDeleted: util.PgtypeBool(),
	}

	role2,err := server.store.UpdateRole(context.Background(),arg)
	require.NoError(t,err)
	require.NotEmpty(t,role2)
	require.Equal(t, arg.ID, role2.ID)
	require.Equal(t, arg.Name, role2.Name)
	require.Equal(t, arg.Description, role2.Description)
	require.Equal(t, arg.IsDeleted, role2.IsDeleted)

	createdTime1 := role1.CreatedAt.Time
	createdTime2 := role2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)

	updatedTime1 := role1.UpdatedAt.Time
	updatedTime2 := role2.UpdatedAt.Time
	require.WithinDuration(t, updatedTime1, updatedTime2, time.Second)

}

func(server *Server) TestDeleteRole(t *testing.T){

	role1 := server.CreateRandomRole(t)

	_,err := server.store.DeleteRole(context.Background(),role1.ID)
	require.NoError(t,err)

	role2,err := server.store.GetRole(context.Background(),role1.ID)
	require.NoError(t,err)
	require.EqualError(t,err,sql.ErrNoRows.Error())
	require.Empty(t,role2)

}

func(server *Server) TestGetAllRole(t *testing.T){
	for i:=0; i<10; i++{
		server.CreateRandomRole(t)
	}

	role2,err := server.store.GetAllRoles(context.Background())
	require.NoError(t,err)
	require.Len(t,role2,5)

	for _,role := range role2{
		require.NotEmpty(t,role)
	}
}

