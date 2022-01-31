package kom_arangodb_paginator

import (
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"log"
	"testing"
)

func TestPaginator_Asd(t *testing.T) {
	var dbConn driver.Database

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{fmt.Sprintf("http://localhost:8529")},
	})
	if err != nil {
		log.Fatal("something went wrong while opening connection...")
	}

	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication("root", ""),
	})
	if err != nil {
		log.Fatal("something went wrong while creating client...")
	}

	if ok, _ := client.DatabaseExists(nil, "kom_product"); !ok {
		// Open {database} database
		dbConn, err = client.CreateDatabase(nil, "kom_product", &driver.CreateDatabaseOptions{})
	} else {
		dbConn, err = client.Database(nil, "kom_product")
	}

	collection, _ := dbConn.Collection(nil, "product")

	var d []*Product
	New(collection, PaginatorOptions{}).TypeSafe().PaginateSafely(&Product{}, func(doc interface{}) {
		d = append(d, doc.(*Product))
	})

	log.Println(d)
}
