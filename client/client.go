package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var answers []int

const API_URL = "http://localhost:8080"

type Question struct {
	Text          string
	Options       []string
	CorrectAnswer int
}

var quizCmd = &cobra.Command{
	Use:   "quiz",
	Short: "Take the quiz",
	Run: func(cmd *cobra.Command, args []string) {
		var username string
		fmt.Println("Enter your username: ")
		fmt.Scanf("%s \n", &username)

		// Get the questions from the API
		resp, err := http.Get(API_URL + "/questions")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		var questions []Question
		err = json.NewDecoder(resp.Body).Decode(&questions)
		if err != nil {
			panic(err)
		}

		// Present the questions to the user
		for i, question := range questions {
			fmt.Printf("\n%d. %s\n", i+1, question.Text)
			for j, option := range question.Options {
				fmt.Printf("%d) %s\n", j+1, option)
			}

			var answer int
			fmt.Println("Enter your answer: ")
			fmt.Scanf("%d \n", &answer) // TODO validate answer 1,2,3,4

			answers = append(answers, answer)
		}

		b := new(bytes.Buffer)
		err = json.NewEncoder(b).Encode(answers)
		if err != nil {
			panic(err)
		}

		// Submit the answers to the API
		resp, err = http.Post(API_URL+"/answers?username="+username, "application/json", b)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		var result struct {
			Correct    int     `json:"correct"`
			Percentage float64 `json:"percentage"`
		}

		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			panic(err)
		}

		// Display the result to the user
		fmt.Printf("You got %d out of %d answers correct\n", result.Correct, len(questions))
		fmt.Printf("You scored higher than %.2f%% of all quizzers\n", result.Percentage)
	},
}

func main() {
	quizCmd.Execute()
}
