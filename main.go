package main

import "fmt"

func main() {
	variableName := "value"
	var pointerToVariable *string = &variableName

	fmt.Printf("value '%v' of type %T at address %v\n",
		*pointerToVariable, *pointerToVariable, &variableName)
}
