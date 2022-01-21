package main

import (
	"NIXRutine/database"
	"NIXRutine/filebase"
	"NIXRutine/models"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
)

func main() {
	conn, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	filePath := "./filebase/data/supplier_"
	database.DeleteTables(conn)
	pool := NewWorkerPool(4)
	wg := sync.WaitGroup{}
	wg.Add(pool.Count)
	for i := 0; i < pool.Count; i++ {
		fmt.Println("start")
		go pool.Run(&wg, func(rest models.Restaurant) {
			restID, err := rest.Insert(conn)
			if err != nil {
				log.Println(err)
			}
			for _, prod := range rest.Menu {
				prodID, err := prod.Insert(conn)
				if err != nil {
					if strings.HasPrefix(err.Error(), "Error 1062") {
						continue
					} else {
						log.Fatal(err)
					}
				}
				restProductQuery := `INSERT INTO menu_products(rest_id, product_id, price) 
					  VALUES (?,?,?)`
				_, err = conn.Exec(restProductQuery, restID, prodID, prod.Price)

				for _, ing := range prod.Ingredients {
					var ingredientsQuery = `INSERT INTO ingredient (name) VALUE(?)`
					//conn.Exec("INSERT INTO ingredient (name) VALUE(?)", ing)
					ingredientRes, err := conn.Exec(ingredientsQuery, ing)
					if err != nil {
						if strings.HasPrefix(err.Error(), "Error 1062") {
							continue
						} else {
							log.Fatal(err)
						}
					}
					ingID, err := ingredientRes.LastInsertId()
					if err != nil {
						log.Fatal(err)
					}
					productIngredientsQuery := `INSERT INTO product_ingredient (product_id, ingredient_id)
											VALUES (?,?)`
					_, err = conn.Exec(productIngredientsQuery, prodID, ingID)
					if err != nil {
						if strings.HasPrefix(err.Error(), "Error 1062") {
							continue
						} else {
							log.Fatal(err)
						}
					}
				}
			}
		})
	}
	restaurants := [7]models.Restaurant{}
	var rest models.Restaurant
	for i := 1; i < 8; i++ {
		rest, err = filebase.ParseFiles(filePath + strconv.Itoa(i) + ".json")
		if err != nil {
			log.Println(err)
		}
		restaurants[i-1] = rest
	}

	for _, rest = range restaurants {
		pool.Sender <- rest
	}

	pool.Stop()
	wg.Wait()
}

/*var restaurants []models.Restaurant
c1 := make(chan models.Restaurant, 7)
for i := 0; i < 7; i++ {
	iterator := i + 1
	go func() {
		err := filebase.ParseFiles(filePath+strconv.Itoa(iterator)+".json", c1)
		if err != nil {
			log.Fatal(err)
		}
	}()
}*/
/*for i := 0; i < 7; i++ {
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
