package main

import (
	"calculator/mathop"
	"fmt"
)

func main() {
	var a, b int
	var op string
	operations := map[string]func(int, int) (float32, error){
		"+": func(a, b int) (float32, error) { return mathops.Add(a, b), nil },
		"-": func(a, b int) (float32, error) { return mathops.Subtract(a, b), nil },
		"*": func(a, b int) (float32, error) { return mathops.Multiply(a, b), nil },
		"/": mathops.Divide,
	}

	for {

		fmt.Print("Enter the first number: ")
		fmt.Scan(&a)
		fmt.Print("Enter the second number: ")
		fmt.Scan(&b)
		fmt.Println("Available operations: + - * /\nDefault operation is +")
		fmt.Print("Select your operation: ")
		fmt.Scan(&op)
		// Use default operation if invalid
		if _, exists := operations[op]; !exists {
			fmt.Println("Invalid operation, defaulting to +")
			op = "+"
		}

		result, err := operations[op](a, b)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("Result of %d %s %d is: %.2f\n", a, op, b, result)
		}

		var cont int
		fmt.Println("Do you want to continue\n1: Yes\n0: No")
		fmt.Scan(&cont)
		if cont == 1 {
			continue
		} else {
			fmt.Println("Exiting...")
			break
		}
	}
}
