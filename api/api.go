package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var scores = make(map[string]int)

type Question struct {
	Text          string
	Options       []string
	CorrectAnswer int
}

var questions = []Question{
	{
		Text: "What is the capital of Portugal?",
		Options: []string{
			"Paris",
			"Lisbon",
			"Bucharest",
			"Madrid",
		},
		CorrectAnswer: 2,
	},
	{
		Text: "What is the second biggest city in Portugal?",
		Options: []string{
			"Porto",
			"Lisbon",
			"Braga",
			"Aveiro",
		},
		CorrectAnswer: 1,
	},
	{
		Text: "Which island is part of Portugal?",
		Options: []string{
			"Bahamas",
			"Islas Canarias",
			"Madeira",
			"Menorca",
		},
		CorrectAnswer: 3,
	},
	{
		Text: "In which region of Portugal is Portimao?",
		Options: []string{
			"Algarve",
			"North",
			"Alentejo",
			"Center",
		},
		CorrectAnswer: 1,
	},
}

func calculatePercentage(newScore int) float64 {
	if len(scores) == 0 {
		return 100
	}

	count := 0
	for _, score := range scores {
		if score < newScore {
			count++
		}
	}

	percentage := float64(count) / float64(len(scores)) * 100
	return percentage
}

func handleQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

func handleAnswers(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	fmt.Println("username:" + username)

	var answers []int
	json.NewDecoder(r.Body).Decode(&answers)

	correct := 0
	for i, answer := range answers {
		if answer == questions[i].CorrectAnswer {
			correct++
		}
	}

	var percentage float64 = calculatePercentage(correct)

	scores[username] = correct

	result := struct {
		Correct    int     `json:"correct"`
		Percentage float64 `json:"percentage"`
	}{
		Correct:    correct,
		Percentage: percentage,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/questions", handleQuestions)
	http.HandleFunc("/answers", handleAnswers)

	fmt.Println("Starting quiz API on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
