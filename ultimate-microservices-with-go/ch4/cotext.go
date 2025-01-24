package main

import (
	"context"
	"fmt"
	"time"
)

func ContextExample() {
	ctx := context.Background()
	childContext := context.WithValue(ctx, "apiKey", 123456)
	printAPIKey(childContext)

	contextWithCancel, cancel := context.WithCancel(ctx)
	cancel()
	fmt.Println("What Happened ?", contextWithCancel.Err())

	contextWithTimeout, _ := context.WithTimeout(ctx, 30*time.Second)
	fmt.Println("What Happened ?", contextWithTimeout.Err())
	time.Sleep(35 * time.Second)
	fmt.Println("What Happened ?", contextWithTimeout.Err())

}

func printAPIKey(ctx context.Context) {
	apiKey := ctx.Value("apiKey")
	fmt.Println("API Key: ", apiKey)
}
