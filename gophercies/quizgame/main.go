package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

type Problem struct {
	Question string
	Answer   string
}

func main() {
	// Open the CSV file
	file, err := os.Open("problems.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Read the CSV data
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	problems := []Problem{}
	// Parse the CSV data into Problem structs
	for _, row := range data {
		// Ensure that there are at least 2 columns (Question and Answer)
		if len(row) < 2 {
			fmt.Println("Skipping invalid row:", row)
			continue
		}
		problems = append(problems, Problem{
			Question: row[0],
			Answer:   row[1],
		})
	}

	t := flag.Int("time", 30, "Quiz Duration") // Default is 18
	flag.Parse()
	result := Quiz(problems, *t)
	total := len(data)
	var message string

	diff := int(math.Abs(float64(total) - float64(result)))

	if diff == 0 {
		message = "Hurray, Good Job!"
	} else if diff < (total / 2) {
		message = "Mhmmmm, Not bad"
	} else {
		message = "Really Bro you can do better"
	}

	fmt.Printf("%s: %d / %d ", message, result, total)
}

// Quiz function displays the problems
func Quiz(problems []Problem, t int) int {
	// Create a timer for t seconds
	timer := time.NewTimer(time.Duration(t) * time.Second)
	sum := 0

	fmt.Println("Starting the quiz... You have", t, "seconds!")

	for i, problem := range problems {
		fmt.Printf("Problem %d: %s = ", i+1, problem.Question)

		answerCh := make(chan int) // Channel to receive user's answer
		go func() {
			var ans int
			fmt.Scan(&ans)
			answerCh <- ans
		}()

		select {
		case <-timer.C:
			// Timer expired
			fmt.Println("\nTime's up!")
			return sum
		case ans := <-answerCh:
			// User provided an answer
			gave, err := strconv.Atoi(problem.Answer)
			if err != nil {
				fmt.Println("Invalid answer in problem list. Skipping question.")
				continue
			}
			if ans == gave {
				sum++
			}
		}
	}

	return sum
}
