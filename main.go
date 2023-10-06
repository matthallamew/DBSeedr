package main

import (
	"DBSeedr/dataGenerator"
	"fmt"
)

func main() {
	generatedSentence := dataGenerator.GenerateRandomData(100)
	fmt.Println(generatedSentence)
}
