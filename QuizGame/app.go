package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type question struct {
	question string
	answer   string
}

func main() {
	// Initialization (set how many seconds here)
	rand.Seed(time.Now().UnixNano())
	seconds := 5
	correct, wrong := 0, 0
	reader := bufio.NewReader(os.Stdin)
	file, err := os.Open("problems.csv")
	if err != nil {
		fmt.Println("Error trying to open problems.csv.")
		return
	}

	// Read questions and save to slice
	contents, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatal(err)
		return
	}
	var questions []question
	questions = make([]question, len(contents))
	for i, value := range contents {
		questions[i] = question{value[0], value[1]}
	}

	// Shuffle questions in slice
	rand.Shuffle(len(questions), func(i, j int) { questions[i], questions[j] = questions[j], questions[i] })

	fmt.Println("Hi! Welcome to go-quiz, a quiz game made in go.")
	fmt.Printf("You are given %d seconds to solve %d questions, good luck!\n", seconds, len(questions))
	fmt.Println("Press enter to begin the quiz.")
	reader.ReadLine()

	timer1 := time.NewTimer(time.Second * time.Duration(seconds))
	go func() {
		<-timer1.C
		fmt.Println("You ran out of time!")
		fmt.Printf("Your final score is %d out of %d.\n", correct, len(questions))
		os.Exit(3)
	}()

	for _, question := range questions {
		fmt.Printf("What is the answer to %s?\n", question.question)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(response)
		if response == question.answer {
			correct += 1
			fmt.Println("Correct!")
		} else {
			fmt.Printf("Incorrect. You answered %s, the correct answer was %s.\n", response, question.answer)
			wrong += 1
		}
	}
	timer1.Stop()
	fmt.Println("You managed to solve all questions!")
	fmt.Printf("Your final score is %d out of %d.", correct, len(questions))
}
