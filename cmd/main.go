package main

import (
	"NIXRutine/database"
	"NIXRutine/filebase"
	"NIXRutine/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func GetProducts(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(resp, "not alowed", http.StatusMethodNotAllowed)
	}
	conn, _ := database.Connect()
	prod := make([]*models.Menu, 0)
	rows, err := conn.Query("SELECT * FROM product")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		pr := new(models.Menu)
		err := rows.Scan(&pr.Id, &pr.Name, &pr.Price, &pr.Image, &pr.Type)
		if err != nil {
			log.Fatal(err)
		}
		prod = append(prod, pr)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	data, _ := json.Marshal(prod)
	resp.Header().Add("Access-Control-Allow-Origin", "*")
	fmt.Fprintln(resp, string(data))
}
func main() {
	server := http.NewServeMux()
	server.HandleFunc("/get_products", GetProducts)
	http.ListenAndServe(":8080", server)
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
				if err != nil {
					fmt.Println(err)
				}
				for _, ing := range prod.Ingredients {
					ingredientsQuery := `INSERT INTO ingredient(name) VALUE(?)`
					//conn.Exec("INSERT INTO ingredient(name) VALUE(?)", ing)
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
					productIngredientsQuery := `INSERT INTO product_ingredient(product_id, ingredient_id)
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
