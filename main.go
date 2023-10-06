package main

import (
	"DBSeedr/dataGenerator"
	"DBSeedr/dbaccess"
	"fmt"
	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	dbdata, err := dbaccess.GetTableSchemaData("Person")
	if err != nil {
		panic(err)
	}

	for _, record := range dbdata {
		fmt.Println(record.SystemDataType, record.MaxLength)
	}

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
