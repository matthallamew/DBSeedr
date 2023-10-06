// Package dataGenerator implements utility functions to generate random data that can be utilized to
// seed a database table(s) with random data for UI testing
package dataGenerator

import (
	"math/rand"
	"time"
)

// GenerateRandomData takes in an int that indicates how many total characters should be generated.
// Spaces are added in between each group of random characters to form a mock sentence.
// The generatedLettersAndSpacesSlice is returned as a string.
func GenerateRandomData(amountToGenerate int) string {
	const lowerAndUpperCaseLettersEng = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	generatedLettersAndSpacesSlice := make([]byte, amountToGenerate)
	stopLen := len(generatedLettersAndSpacesSlice)
	savedLetLen, remainderLen := 0, 0
	for savedLetLen < stopLen {
		randBetween1And9 := getRandomNumBetween(1, 9)
		if (savedLetLen + randBetween1And9) < stopLen {
			// Add as many random letters as is equal to randBetween1And9
			for idx := 0; idx < randBetween1And9; idx++ {
				rando := getRandomNumBetween(0, 51)
				generatedLettersAndSpacesSlice[savedLetLen] = lowerAndUpperCaseLettersEng[rando]
				savedLetLen++
			}
			// Add a space after we add random letters to separate out the random letter groups
			if savedLetLen > 0 {
				generatedLettersAndSpacesSlice[savedLetLen] = ' '
			}
		}
		savedLetLen++

		// Bail out of the for loop if savedLetLen + randBetween1And9
		// is greater than the total length of the slice
		if (savedLetLen + randBetween1And9) >= stopLen {
			remainderLen = stopLen - savedLetLen
			savedLetLen = stopLen
		}
	}

	// Fill up the remaining space in the array with characters
	lastToAddLen := stopLen - remainderLen
	if remainderLen > 0 {
		for idx := 0; idx < remainderLen; idx++ {
			rando := getRandomNumBetween(0, 51)
			generatedLettersAndSpacesSlice[lastToAddLen] = lowerAndUpperCaseLettersEng[rando]
			lastToAddLen++
		}
	}
	return string(generatedLettersAndSpacesSlice)
}

// getRandomNumBetween takes in two integers, one to indicate the minimum and one to indicate the maximum.
// A randomly generated integer between the minimum the maximum (inclusive) is returned.
func getRandomNumBetween(minRandom, maxRandom int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(maxRandom-minRandom+1) + minRandom
}
