package main

import (
	"DBSeedr/dataGenerator"
	"DBSeedr/dbaccess"
	"database/sql"
	"fmt"
	"strings"
	"sync"
)

func main() {
	const tableToFillWithRandomData = "Addr"
	//const tableToFillWithRandomData = "Person"
	dbMetaData := getDbMetaData(tableToFillWithRandomData)

	// Add any DB field name and value pairs that you want to fill with static data
	// Example: non-nullable Foreign Keys will need to be in here, otherwise the insert will fail
	staticFields := make(map[string]any)
	//staticFields["state"] = "MN"
	//staticFields["addressID"] = 1

	var filteredMetaData []dbaccess.TableSchemaData
	for idx, dbField := range dbMetaData {
		if _, exists := staticFields[dbField.ColumnName]; !exists {
			filteredMetaData = append(filteredMetaData, dbMetaData[idx])
		}
	}

	var waitGroup sync.WaitGroup
	for idx := 0; idx < 5; idx++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			seedDb(tableToFillWithRandomData, staticFields, filteredMetaData)
		}()
	}
	waitGroup.Wait()
}

// getDbMetaData uses the dbaccess package to get the database metadata and return the result.
func getDbMetaData(tableToFillWithRandomData string) []dbaccess.TableSchemaData {
	dbMetaData, err := dbaccess.GetTableSchemaData(tableToFillWithRandomData)
	if err != nil {
		panic(err)
	}
	return dbMetaData
}

// seedDb will build and execute an insert query for the given table, database metadata, and any fields the user has chosen to fill on their own.
func seedDb(tableToFillWithRandomData string, staticFields map[string]any, dbMetaData []dbaccess.TableSchemaData) {
	var sbQueryUpper, sbQueryLower, fullQuery strings.Builder
	sbQueryUpper.WriteString("INSERT INTO dbo.")
	sbQueryUpper.WriteString(tableToFillWithRandomData)
	sbQueryUpper.WriteString("(")
	sbQueryLower.WriteString(" VALUES (")
	metaDataSliceLength := len(dbMetaData) - 1
	var generatedVals []any
	for idx, dbField := range dbMetaData {
		_, fieldIsStaticallyFilled := staticFields[dbField.ColumnName]
		// If the dbField is an identity field, skip it as it should get auto generated when a record is inserted
		// If the dbField is in the list of staticFields, skip auto generating it
		if !dbField.IsIdentity && !fieldIsStaticallyFilled {
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
	// Add any staticFields records to the insert query
	if totalStaticFields := len(staticFields); totalStaticFields > 0 {
		insertLen := 0
		for key, value := range staticFields {
			if insertLen < totalStaticFields && len(generatedVals) > 0 {
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
