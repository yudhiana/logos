package main

import (
	"context"
	"fmt"
)

func main() {
	// Defining keys
	IrisContextKey := "iris_key"
	PanicContextKey := "panic_key"

	// Creating contexts with values
	firstCtx := context.WithValue(context.TODO(), IrisContextKey, "first data")
	newCtx := context.WithValue(firstCtx, PanicContextKey, "new ctx")

	// Retrieve the value associated with PanicContextKey
	value := newCtx.Value(IrisContextKey)

	// Print the value
	fmt.Println("Value for PanicContextKey:", value)

}
