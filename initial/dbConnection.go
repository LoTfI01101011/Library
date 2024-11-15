package initial

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	loadEnv()
}
func ConnectToDb() {

	var err error
	dsn := "host=localhost" + " user=" + os.Getenv("User") + " password=" + os.Getenv("PASSWORD") + " dbname=" + os.Getenv("DBNAME") + " port=5432 sslmode=disable "
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connecting the db")
	}

}
