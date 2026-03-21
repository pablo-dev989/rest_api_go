package sqlconnect

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDb() (*sql.DB, error) {
	// err := godotenv.Load()
	// if err != nil {
	// 	return nil, err
	// }

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("HOST")
	dbport := os.Getenv("DB_PORT")

	// connectionString := "root:sasa123@tcp(host.docker.internal:3306)/" + dbname
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, dbport, dbname)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		// panic(err)
		return nil, err
	}

	return db, nil
}
