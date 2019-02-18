package psql

import (
	"fmt"
	"github.com/OIT-ads-web/widgets_import"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

var Database *sqlx.DB

func GetConnection() *sqlx.DB {
	return Database
}

func MakeConnection(conf widgets_import.Config) error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Database.Server, conf.Database.Port,
		conf.Database.User, conf.Database.Password,
		conf.Database.Database)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Println("m=GetPool,msg=connection has failed", err)
	}

    Database = db
    return err
}
