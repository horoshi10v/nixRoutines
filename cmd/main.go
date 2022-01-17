package main

import (
	"NIXRutine/database"
	"NIXRutine/filebase"
	"NIXRutine/models"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
)

func GetRowId(db *sql.DB, selectQuery, insertQuery string, args ...interface{}) int64 {
	row := db.QueryRow(selectQuery, args...)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		if strings.HasPrefix(err.Error(), "no rows in") {
			log.Println(err)
		}
		result, err := db.Exec(insertQuery, args...)
		if err != nil {
			if strings.HasPrefix(err.Error(), "Error 1062") {
				return GetRowId(db, selectQuery, insertQuery, args...)
			}
			log.Fatalln(err)
		}
		id, err = result.LastInsertId()
		if err != nil {
			log.Fatalln(err)
		}
	}
	return id
}

func main() {
	conn, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	filePath := "./filebase/data/supplier_"
	pool := NewWorkerPool(4)
	wg := sync.WaitGroup{}
	wg.Add(pool.Count)
	for i := 0; i < pool.Count; i++ {
		fmt.Println("start")
		go pool.Run(&wg, func(rest models.Restaurant) {
			shopTypeId := GetRowId(conn, "SELECT id FROM shop_types WHERE name = ?",
				"INSERT INTO shop_types(name) VALUE (?)", rest.Type)
			_, err = conn.Exec("INSERT INTO restaurant VALUE (?, ?, ?, ?, ?, ?)",
				rest.Id, rest.Name, shopTypeId, rest.Image, rest.WorkingHours.Opening,
				rest.WorkingHours.Closing)

			if err != nil {
				log.Println(err)
			}

			for _, prod := range rest.Menu {
				prodTypeId := GetRowId(conn, "SELECT id FROM prod_types WHERE name = ?",
					"INSERT INTO prod_types(name) VALUE (?)", prod.Type)
				_, err = conn.Exec(
					"INSERT INTO products VALUE (?, ?, ?, ?, ?, ?)",
					prod.Id, prod.Name, prod.Price, prod.Image, prodTypeId, rest.Id)

				if err != nil {
					log.Println(err)
				}

				for _, ing := range prod.Ingredients {
					ingId := GetRowId(conn, "SELECT id FROM ingredients WHERE name = ?",
						"INSERT INTO ingredients(name) VALUE (?)", ing)
					_, err = conn.Exec("INSERT INTO product_ingredient VALUE (?, ?)", prod.Id, ingId)

					if err != nil {
						log.Println(err)
					}
				}
			}
		})
	}
	shops := [7]models.Restaurant{}
	var shop models.Restaurant
	for i := 1; i < 8; i++ {
		shop, err = filebase.ParseFiles(filePath+strconv.Itoa(i)+".json", i)
		if err != nil {
			log.Println(err)
		}
		shops[i-1] = shop
	}

	for _, shop = range shops {
		pool.Sender <- shop
	}

	pool.Stop()
	wg.Wait()
	/*
		var restaurants []models.Restaurant
		c1 := make(chan models.Restaurant, 7)
		for i := 0; i < 7; i++ {
			iterator := i + 1
			go func() {
				err := filebase.ParseFiles(filePath+strconv.Itoa(iterator)+".json", c1)
				if err != nil {
					log.Fatal(err)
				}
			}()
		}
		for i := 0; i < 7; i++ {
			restaurants = append(restaurants, <-c1)

		}

		database.DropTable1(conn)
		var wg sync.WaitGroup
		for i := 0; i < 7; i++ {
			iterator := i
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := restaurants[iterator].Insert(conn)
				if err != nil {
					log.Fatal(err)
				}
			}()
		}
		wg.Wait()*/
}
