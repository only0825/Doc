package helper

import (
	"github.com/ipipdotnet/ipdb-go"
	"log"
)

func IpdbInit(path string) *ipdb.City {

	db, err := ipdb.NewCity(path)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
