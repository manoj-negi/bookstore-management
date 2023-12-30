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

func(server *Server) CreateRandomCountry(t *testing.T) db.Country{
	arg:= db.CreateCountryParams{
		Iso2: util.RandomName(),
		ShortName: util.RandomName(),
		LongName:  util.RandomName(),
		Numcode:   util.PgtypeText(),
		CallingCode: util.RandomName(),
		Cctld:  util.RandomName(),
		IsDeleted: util.PgtypeBool(),
	}

	country, err := server.store.CreateCountry(context.Background(), arg)
	require.NoError(t,err)
	require.NotEmpty(t,country)
	require.Equal(t, arg.Iso2, country.Iso2)
	require.Equal(t, arg.ShortName, country.ShortName)
	require.Equal(t, arg.LongName, country.LongName)
	require.Equal(t, arg.Numcode, country.Numcode)
	require.Equal(t, arg.CallingCode, country.CallingCode)
	require.Equal(t, arg.Cctld, country.Cctld)
	require.Equal(t, arg.IsDeleted, country.IsDeleted)

	require.NotZero(t,country.ID)
	require.NotZero(t,country.CreatedAt)
	require.NotZero(t,country.UpdatedAt)

	return country
}

func(server *Server) TestCreateCountry(t *testing.T) {
	server.CreateRandomCountry(t)
}

func(server * Server) TestGetCountry(t *testing.T){
	country1 := server.CreateRandomCountry(t)

	country2, err := server.store.GetCountry(context.Background(), country1.ID)
	require.NoError(t,err)
	require.NotEmpty(t,country2)
	require.Equal(t, country1.ID, country2.ID)
	require.Equal(t, country1.Iso2, country2.Iso2)
	require.Equal(t, country1.ShortName, country2.ShortName)
	require.Equal(t, country1.LongName, country2.LongName)
	require.Equal(t, country1.Numcode, country2.Numcode)
	require.Equal(t, country1.CallingCode, country2.CallingCode)
	require.Equal(t, country1.Cctld, country2.Cctld)
	require.Equal(t, country1.IsDeleted, country2.IsDeleted)

	createdTime1 := country1.CreatedAt.Time
	createdTime2 := country2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)

	updatedTime1 := country1.UpdatedAt.Time
	updatedTime2 := country2.UpdatedAt.Time
	require.WithinDuration(t, updatedTime1, updatedTime2, time.Second)
}

func(server *Server) TestUpdateCountry(t *testing.T){
	country1 := server.CreateRandomCountry(t)

	arg := db.UpdateCountryParams{
		ID: country1.ID,
		Iso2: util.RandomName(),
		ShortName: util.RandomName(),
		LongName:  util.RandomName(),
		Numcode:   util.PgtypeText(),
		CallingCode: util.RandomName(),
		Cctld:  util.RandomName(),
		IsDeleted: util.PgtypeBool(),
	}

	country2,err := server.store.UpdateCountry(context.Background(),arg)
	require.NoError(t,err)
	require.NotEmpty(t,country2)
	require.Equal(t, arg.ID, country2.ID)
	require.Equal(t, arg.Iso2, country2.Iso2)
	require.Equal(t, arg.ShortName, country2.ShortName)
	require.Equal(t, arg.LongName, country2.LongName)
	require.Equal(t, arg.Numcode, country2.Numcode)
	require.Equal(t, arg.CallingCode, country2.CallingCode)
	require.Equal(t, arg.Cctld, country2.Cctld)
	require.Equal(t, arg.IsDeleted, country2.IsDeleted)

	createdTime1 := country1.CreatedAt.Time
	createdTime2 := country2.CreatedAt.Time
	require.WithinDuration(t, createdTime1, createdTime2, time.Second)

	updatedTime1 := country1.UpdatedAt.Time
	updatedTime2 := country2.UpdatedAt.Time
	require.WithinDuration(t, updatedTime1, updatedTime2, time.Second)
}

func(server *Server) TestDeleteCountry(t *testing.T){

	country1 := server.CreateRandomCountry(t)

	_,err := server.store.DeleteCountry(context.Background(),country1.ID)
	require.NoError(t,err)

	country2,err := server.store.GetCountry(context.Background(),country1.ID)
	require.NoError(t,err)
	require.EqualError(t,err,sql.ErrNoRows.Error())
	require.Empty(t,country2)

}

func(server *Server) TestGetAllCountry(t *testing.T){
	for i:=0; i<10; i++{
		server.CreateRandomCountry(t)
	}


	country2,err := server.store.GetAllCountries(context.Background())
	require.NoError(t,err)
	require.Len(t,country2,5)

	for _,country := range country2{
		require.NotEmpty(t,country)
	}
}

