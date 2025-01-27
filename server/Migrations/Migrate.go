package Migrations

import (
	"log"
	"os"

	"forum/GlobVar"
)

func Migrate() {
	query, err := os.ReadFile("../Database/modules.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = GlobVar.DB.Exec(string(query))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database migrated successfully!")
}
