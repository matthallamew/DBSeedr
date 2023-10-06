package main

import (
	"DBSeedr/dataGenerator"
	"fmt"
)

func main() {
	generatedSentence, err := dataGenerator.GenerateRandomData("string", 100)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v is of type %[1]T \n", generatedSentence)

	generatedNum, err := dataGenerator.GenerateRandomData("int", 9998)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v is of type %[1]T \n", generatedNum)

	generatedBadType, err := dataGenerator.GenerateRandomData("badType", 10)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v is of type %[1]T \n", generatedBadType)
}
