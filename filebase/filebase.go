package filebase

import (
	"NIXRutine/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func ParseFiles(filePath string) (models.Restaurant, error) {
	var restaurant models.Restaurant
	jsonFile, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
	}
	byteValue, err := ioutil.ReadAll(jsonFile)

	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.Println(err)
		}
	}(jsonFile)

	err = json.Unmarshal(byteValue, &restaurant)
	if err != nil {
		log.Println(err)
	}
	//c1 <- restaurant
	return restaurant, nil
}
