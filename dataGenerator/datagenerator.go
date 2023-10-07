// Package dataGenerator implements utility functions to generate random data that can be utilized to
// seed a database table(s) with random data for UI testing.
package dataGenerator

import (
	"errors"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// GenerateRandomData takes in a string dataType that indicates the database field's data type
// and an int amountToGenerate that indicates how many total characters should be generated.
func GenerateRandomData(dataType string, amountToGenerate int) (any, error) {
	dataTypeLower := strings.ToLower(dataType)
	switch {
	case isStringType(dataTypeLower):
		return generateRandomString(amountToGenerate), nil
	case isIntType(dataTypeLower):
		return generateRandomNumBetween(1, amountToGenerate, "integer"), nil
	case isDecimalType(dataTypeLower):
		return generateRandomNumBetween(1, amountToGenerate, "decimal"), nil
	case isBoolType(dataTypeLower):
		return generateRandomNumBetween(0, amountToGenerate, "integer"), nil
	default:
		return "", errors.New("Cannot generate for data type " + dataType)
	}
}

// isStringType checks the given typeStr to see if it fits within a String database type.
// It will return true if it does or false if it does not.
func isStringType(typeStr string) bool {
	switch typeStr {
	case "varchar", "char", "text":
		return true
	case "nvarchar", "nchar", "ntext":
		return true
	default:
		return false
	}
}

// isIntType checks the given typeStr to see if it fits within an Integer database type.
// It will return true if it does or false if it does not.
func isIntType(typeStr string) bool {
	switch typeStr {
	case "bigint", "int", "smallint", "tinyint":
		return true
	default:
		return false
	}
}

// isDecimalType checks the given typeStr to see if it fits within a Decimal/Float database type.
// It will return true if it does or false if it does not.
func isDecimalType(typeStr string) bool {
	switch typeStr {
	case "decimal", "numeric":
		return true
	case "float", "real":
		return true
	case "money", "smallmoney":
		return true
	default:
		return false
	}
}

// isBoolType checks the given typeStr to see if it fits within a Boolean database type.
// It will return true if it does or false if it does not.
func isBoolType(typeStr string) bool {
	switch typeStr {
	case "bit":
		return true
	default:
		return false
	}
}

// generateRandomString will randomly choose a number of characters up to the amountToGenerate.
// Spaces are added in between each group of random characters (between 1 and 9 characters) to form a mock sentence.
// The generatedLettersAndSpacesSlice is returned as a string.
func generateRandomString(amountToGenerate int) string {
	sliceSize := amountToGenerate
	// Limit data types that allow "unlimited" data to a more reasonable amount
	if sliceSize == -1 {
		sliceSize = 30000
	}
	const lowerAndUpperCaseLettersEng = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const lettersLengthMax = len(lowerAndUpperCaseLettersEng) - 1
	generatedLettersAndSpacesSlice := make([]byte, sliceSize)
	stopLen := sliceSize
	savedLetterLen, remainderLen := 0, 0
	for savedLetterLen < stopLen {
		randBetween1And9 := generateRandomNumBetween(1, 9, "integer").(int)
		if (savedLetterLen + randBetween1And9) < stopLen {
			// Add as many random letters as is equal to randBetween1And9
			for idx := 0; idx < randBetween1And9; idx++ {
				randomLetterIdx := generateRandomNumBetween(0, lettersLengthMax, "integer").(int)
				generatedLettersAndSpacesSlice[savedLetterLen] = lowerAndUpperCaseLettersEng[randomLetterIdx]
				savedLetterLen++
			}
			// Add a space after we add random letters to separate out the random letter groups
			if savedLetterLen > 0 {
				generatedLettersAndSpacesSlice[savedLetterLen] = ' '
			}
		}
		savedLetterLen++

		// Bail out of the for loop if savedLetterLen + randBetween1And9
		// is greater than the total length of the slice
		if (savedLetterLen + randBetween1And9) >= stopLen {
			remainderLen = stopLen - savedLetterLen
			savedLetterLen = stopLen
		}
	}

	// Fill up the remaining space in the array with characters
	lastToAddLen := stopLen - remainderLen
	if remainderLen > 0 {
		for idx := 0; idx < remainderLen; idx++ {
			randomLetterIdx := generateRandomNumBetween(0, lettersLengthMax, "integer").(int)
			generatedLettersAndSpacesSlice[lastToAddLen] = lowerAndUpperCaseLettersEng[randomLetterIdx]
			lastToAddLen++
		}
	}
	return string(generatedLettersAndSpacesSlice)
}

// generateRandomNumBetween takes in two integers, one to indicate the minimum and one to indicate the maximum.
// A randomly generated integer between the minimum the maximum (inclusive) is returned.
func generateRandomNumBetween(minRandom, maxRandom int, returnType string) any {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	switch returnType {
	case "integer":
		return rand.Intn(maxRandom-minRandom+1) + minRandom
	case "decimal":
		numOfZeros := math.Ceil(float64(maxRandom) / 3)
		multiplyStr := "1"
		for idx := 0; idx < int(numOfZeros); idx++ {
			multiplyStr += "0"
		}
		multiFloat, err := strconv.ParseFloat(multiplyStr, 64)
		if err != nil {
			multiFloat = 100
		}
		return rand.Float64() * ((float64(maxRandom) - float64(minRandom) + 1) + float64(minRandom)) * multiFloat
	default:
		return 0
	}
}
