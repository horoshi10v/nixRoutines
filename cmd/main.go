package main

import (
	"NIXRutine/database"
	"NIXRutine/filebase"
)

func main() {
	database.Connect()
	filebase.OpenFiles()
}
