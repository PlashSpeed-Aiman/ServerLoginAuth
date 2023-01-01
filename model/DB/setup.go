package DB

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup() *pgxpool.Pool {

	dsn := "host=localhost user= password= dbname= port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, _ := pgxpool.New(context.Background(), dsn)

	err := db.Ping(context.Background())
	if err != nil {
		fmt.Println(err.Error())
	}

	return db
}
