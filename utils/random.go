package util

import(
	"math/rand"
	"strings"
	"time"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/vod/db/sqlc"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init(){
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min,max int64) int64{
	return min + rand.Int63n(max-min+1)
}

func RandomInt32() int32{
	var min int32
	var max int32
	return min + int32(rand.Intn(int(max - min + 1)))
}


func GenerateRandomStatus() db.StatusEnum {
    statuses := []db.StatusEnum{
        db.StatusEnumPending,
        db.StatusEnumInProcess,
        db.StatusEnumCompleted,
    }
    return statuses[rand.Intn(len(statuses))]
}

func GenerateRandomStatus1() db.PaymentStatusEnum {
    statuses := []db.PaymentStatusEnum{
        db.PaymentStatusEnumPending,
        db.PaymentStatusEnumInProcess,
        db.PaymentStatusEnumCompleted,
    }
    return statuses[rand.Intn(len(statuses))]
}

func GenerateRandomGenderStatus1() db.GenderEnum {
    statuses := []db.GenderEnum{
        db.GenderEnumMale,
        db.GenderEnumFemale,
    }
    return statuses[rand.Intn(len(statuses))]
}


func RandomString(n int) string{
	var sb strings.Builder
	k := len(alphabet)

	for i:=0; i<n; i++{
		c:= alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomName() string{
	return RandomString(6)
}

func PgtypeDate() pgtype.Date {
	var date time.Time
    return pgtype.Date{Time: date}
}

func PgtypeText() pgtype.Text {
	return pgtype.Text{String: "abca.jpg"}
}

func PgtypeBool() pgtype.Bool{
	return pgtype.Bool{Bool: false}
}
