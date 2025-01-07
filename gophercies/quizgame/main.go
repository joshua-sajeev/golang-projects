package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
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

	result := Quiz(problems)

	total := len(data)
	fmt.Printf("Your final score is: %d / %d ", result, total)
}

// Quiz function displays the problems
func Quiz(problems []Problem) int {
	fmt.Println("Starting the quiz...")
	sum := 0
	for i, problem := range problems {
		var ans int
		fmt.Printf("Problem %d: %s\n", i+1, problem.Question)
		// For simplicity, let's print the answer too, but in a real quiz scenario, you can prompt the user for answers.
		fmt.Scan(&ans)
		gave, err := strconv.Atoi(problem.Answer)
		if err != nil {
			continue
		}
		if ans == gave {
			sum += 1
		}
	}
	return sum
}
