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

func(server *Server) CreateRandomUser(t *testing.T) db.User{
	arg:= db.CreateUserParams{
		FirstName: util.RandomName(),
		LastName: util.RandomName(),
		Gender: util.GenerateRandomGenderStatus1(),
		Dob: util.PgtypeDate(),
		Address: util.RandomName(),
		City: util.RandomName(),
		State: util.RandomName(),
		CountryID: util.RandomInt32(),
		MobileNo: util.RandomName(),
		Username: util.RandomName(),
		Email: util.RandomName(),
		Password: util.RandomName(),
		RoleID: util.RandomInt32(),
		Otp: util.RandomInt32(),
		IsDeleted: util.PgtypeBool(),
	}

	user, err := server.store.CreateUser(context.Background(), arg)
	require.NoError(t,err)
	require.NotEmpty(t,user)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)
	require.Equal(t, arg.Gender, user.Gender)
	require.Equal(t, arg.Dob, user.Dob)
	require.Equal(t, arg.Address, user.Address)
	require.Equal(t, arg.City, user.City)
	require.Equal(t, arg.State, user.State)
	require.Equal(t, arg.CountryID, user.CountryID)
	require.Equal(t, arg.MobileNo, user.MobileNo)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Password, user.Password)
	require.Equal(t, arg.RoleID, user.RoleID)
	require.Equal(t, arg.Otp, user.Otp)
	require.Equal(t, arg.IsDeleted, user.IsDeleted)

	require.NotZero(t,user.ID)
	require.NotZero(t,user.CreatedAt)
	require.NotZero(t,user.UpdatedAt)

	return user
}

func(server *Server) TestCreateUser(t *testing.T) {
	server.CreateRandomUser(t)
}

func(server * Server) TestGetUser(t *testing.T){
	user1 := server.CreateRandomUser(t)

	user2, err := server.store.GetUser(context.Background(), user1.ID)
	require.NoError(t,err)
	require.NotEmpty(t,user2)
	require.Equal(t, user1.FirstName, user2.FirstName)
	require.Equal(t, user1.LastName, user2.LastName)
	require.Equal(t, user1.Gender, user2.Gender)
	require.Equal(t, user1.Dob, user2.Dob)
	require.Equal(t, user1.Address, user2.Address)
	require.Equal(t, user1.City, user2.City)
	require.Equal(t, user1.State, user2.State)
	require.Equal(t, user1.CountryID, user2.CountryID)
	require.Equal(t, user1.MobileNo, user2.MobileNo)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.RoleID, user2.RoleID)
	require.Equal(t, user1.Otp, user2.Otp)
	require.Equal(t, user1.IsDeleted, user2.IsDeleted)

	createdTime1 := user1.CreatedAt.Time
	createdTime2 := user2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)

	updatedTime1 := user1.UpdatedAt.Time
	updatedTime2 := user2.UpdatedAt.Time
	require.WithinDuration(t, updatedTime1, updatedTime2, time.Second)
}

func(server *Server) TestUpdateUser(t *testing.T){
	user1 := server.CreateRandomUser(t)

	arg := db.UpdateUserParams{
		ID: user1.ID,
		FirstName: util.RandomName(),
		LastName: util.RandomName(),
		Gender: util.GenerateRandomGenderStatus1(),
		Dob: util.PgtypeDate(),
		Address: util.RandomName(),
		City: util.RandomName(),
		State: util.RandomName(),
		CountryID: util.RandomInt32(),
		MobileNo: util.RandomName(),
		Username: util.RandomName(),
		Email: util.RandomName(),
		Password: util.RandomName(),
		RoleID: util.RandomInt32(),
		Otp: util.RandomInt32(),
		IsDeleted: util.PgtypeBool(),
	}

	user2,err := server.store.UpdateUser(context.Background(),arg)
	require.NoError(t,err)
	require.NotEmpty(t,user2)
	require.Equal(t, arg.ID, user2.ID)
	require.Equal(t, arg.FirstName, user2.FirstName)
	require.Equal(t, arg.LastName, user2.LastName)
	require.Equal(t, arg.Gender, user2.Gender)
	require.Equal(t, arg.Dob, user2.Dob)
	require.Equal(t, arg.Address, user2.Address)
	require.Equal(t, arg.City, user2.City)
	require.Equal(t, arg.State, user2.State)
	require.Equal(t, arg.CountryID, user2.CountryID)
	require.Equal(t, arg.MobileNo, user2.MobileNo)
	require.Equal(t, arg.Username, user2.Username)
	require.Equal(t, arg.Email, user2.Email)
	require.Equal(t, arg.Password, user2.Password)
	require.Equal(t, arg.RoleID, user2.RoleID)
	require.Equal(t, arg.Otp, user2.Otp)
	require.Equal(t, arg.IsDeleted, user2.IsDeleted)

	createdTime1 := user1.CreatedAt.Time
	createdTime2 := user2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)

	updatedTime1 := user1.UpdatedAt.Time
	updatedTime2 := user2.UpdatedAt.Time
	require.WithinDuration(t, updatedTime1, updatedTime2, time.Second)

}

func(server *Server) TestDeleteUser(t *testing.T){

	user1 := server.CreateRandomUser(t)

	_,err := server.store.DeleteUser(context.Background(),user1.ID)
	require.NoError(t,err)

	user2,err := server.store.GetUser(context.Background(),user1.ID)
	require.NoError(t,err)
	require.EqualError(t,err,sql.ErrNoRows.Error())
	require.Empty(t,user2)

}

func(server *Server) TestGetAllUser(t *testing.T){
	for i:=0; i<10; i++{
		server.CreateRandomUser(t)
	}

	user2,err := server.store.GetAllUsers(context.Background())
	require.NoError(t,err)
	require.Len(t,user2,5)

	for _,user := range user2{
		require.NotEmpty(t,user)
	}
}

