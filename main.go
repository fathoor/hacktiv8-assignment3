package main

import (
	"encoding/json"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type Input struct {
	Status Status `json:"status"`
}

func reload() {
	const (
		max = 100
		min = 1
	)

	for {
		input := Input{
			Status: Status{
				Water: rand.Intn(max-min) + min,
				Wind:  rand.Intn(max-min) + min,
			},
		}

		inputJSON, err := json.MarshalIndent(&input, "", "  ")
		if err != nil {
			log.Fatal("Error marshalling input")
		}

		err = os.WriteFile("input.json", inputJSON, 0644)
		if err != nil {
			log.Fatal("Error writing input.json")
		}

		time.Sleep(15 * time.Second)
	}
}

func render(w http.ResponseWriter, r *http.Request) {
	inputJSON, err := os.ReadFile("input.json")
	if err != nil {
		log.Fatal("Error reading input.json")
	}

	var input Input
	err = json.Unmarshal(inputJSON, &input)
	if err != nil {
		log.Fatal("Error unmarshalling input.json")
	}

	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatal("Error parsing index.html")
	}

	var waterStatus, windStatus string
	water := input.Status.Water
	wind := input.Status.Wind

	if water < 6 {
		waterStatus = "Aman"
	} else if water >= 6 && water <= 8 {
		waterStatus = "Siaga"
	} else {
		waterStatus = "Bahaya"
	}

	if wind < 7 {
		windStatus = "Aman"
	} else if wind >= 7 && wind <= 15 {
		windStatus = "Siaga"
	} else {
		windStatus = "Bahaya"
	}

	data := map[string]interface{}{
		"water":       water,
		"wind":        wind,
		"waterStatus": waterStatus,
		"windStatus":  windStatus,
	}

	t.Execute(w, data)
}

func main() {
	go reload()

	http.HandleFunc("/", render)

	HOST := os.Getenv("APP_HOST")
	PORT := os.Getenv("APP_PORT")

	http.ListenAndServe(HOST+":"+PORT, nil)
}
