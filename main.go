package main

import (
	"DBSeedr/dataGenerator"
	"DBSeedr/dbaccess"
	"database/sql"
	"fmt"
	"strings"
)

func main() {
	const tableToFillWithRandomData = "Addr"
	//const tableToFillWithRandomData = "Person"

	// Add any known Foreign Key name and value pairs here to ensure they get inserted
	hardCodedFKs := make(map[string]int)
	//hardCodedFKs["addressID"] = 1

	dbMetaData, err := dbaccess.GetTableSchemaData(tableToFillWithRandomData)
	if err != nil {
		panic(err)
	}

	var sbQueryUpper, sbQueryLower, fullQuery strings.Builder
	sbQueryUpper.WriteString("INSERT INTO dbo.")
	sbQueryUpper.WriteString(tableToFillWithRandomData)
	sbQueryUpper.WriteString("(")
	sbQueryLower.WriteString(" VALUES (")
	metaDataSliceLength := len(dbMetaData) - 1
	var generatedVals []any
	for idx, dbField := range dbMetaData {
		// if the dbField is an identity field, skip it as it should get auto generated when a record is inserted
		if !dbField.IsIdentity {
			sbQueryUpper.WriteString("[" + dbField.ColumnName + "]")
			sbQueryLower.WriteString("?")
			// if we are not done looping through metadata, add a comma before the next field and placeholder
			if idx < metaDataSliceLength {
				sbQueryUpper.WriteString(", ")
				sbQueryLower.WriteString(", ")
			}
			generated, err := dataGenerator.GenerateRandomData(dbField.SystemDataType, dbField.MaxLength)
			// Data could not be generated for the given type, append the NULL string to the generatedVals slice
			if err != nil {
				fmt.Println(err)
				generatedVals = append(generatedVals, sql.NullString{})
			}
			// Data was generated, append it to the generatedVals slice
			if err == nil {
				generatedVals = append(generatedVals, generated)
			}
		}
	}
	// Add any Foreign Keys to the insert query
	if totalFKs := len(hardCodedFKs); totalFKs > 0 {
		insertLen := 0
		for key, value := range hardCodedFKs {
			if insertLen < totalFKs {
				sbQueryUpper.WriteString(", ")
				sbQueryLower.WriteString(", ")
			}
			sbQueryUpper.WriteString("[" + key + "]")
			sbQueryLower.WriteString("?")
			generatedVals = append(generatedVals, value)
			insertLen++
		}
	}

	sbQueryUpper.WriteString(")")
	sbQueryLower.WriteString(")")
	fullQuery.WriteString(sbQueryUpper.String())
	fullQuery.WriteString(sbQueryLower.String())
	dbaccess.InsertGeneratedData(fullQuery.String(), generatedVals...)
}
