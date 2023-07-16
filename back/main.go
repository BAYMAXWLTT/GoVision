package main

import (
	"log"
)

func main() {
	frontAddress := `localhost:8000`
	classificationAddress := `http://localhost:4000`
	styleAddress := `http://localhost:3000`
	app := NewApplication(frontAddress, classificationAddress, styleAddress)
	log.Fatal(app.Run())
}
