package filebase

import (
	"NIXRutine/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func OpenFiles() {

	for i := 0; i < 7; i++ {
		filename := "./filebase/data/supplier_"
		//var full = filename + i + ".json"
		jsonFile, err := os.Open(filename + i + ".json")
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println("Successfully Opened %#v.json", i)
		// read our opened xmlFile as a byte array.
		byteValue, _ := ioutil.ReadAll(jsonFile)
		defer jsonFile.Close()
		// we initialize our Restaurant array
		var menu models.Menu

		// we unmarshal our byteArray which contains our
		err = json.Unmarshal(byteValue, &menu)
		if err != nil {
			log.Println(err)
		}

		for i := 0; i < len(menu.Menu); i++ {
			fmt.Println("ID: ", menu.Menu[i].Id)
			fmt.Println("NAME", menu.Menu[i].Name)
			fmt.Println("PRICE", menu.Menu[i].Price)
			fmt.Println("IMAGE", menu.Menu[i].Image)
		}
	}

}
