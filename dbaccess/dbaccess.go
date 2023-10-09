// Package dbaccess provides utility functions to take care of figuring out database table metadata and inserting generated data.
// Currently, this package is heavily tied to Microsoft SQL Server.
package dbaccess

import (
	"database/sql"
	"fmt"
	_ "github.com/microsoft/go-mssqldb"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

// TableSchemaData defines the metadata query columns used to determine the data type and amount to generate.
type TableSchemaData struct {
	ColumnName     string
	SystemDataType string
	MaxLength      int
	IsNullable     bool
	IsIdentity     bool
}

// GetTableSchemaData Connects to a database and queries the system tables to gather metadata about the given table.
// Return a TableSchemaData slice filled with metadata for the given table.
func GetTableSchemaData(tableName string) ([]TableSchemaData, error) {
	db := connectToDb()
	defer db.Close()
	queryStr := `SELECT AC.[name] AS [columnName], TY.[name] AS systemDataType, AC.[max_length] AS maxLength, 
AC.[is_nullable] AS isNullable, CASE WHEN AC.is_identity=1 THEN 1 
WHEN FAC.foreignKeyName IS NOT NULL THEN 1 
ELSE 0 END AS isIdentity
FROM sys.[tables] AS T
INNER JOIN sys.[all_columns] AC ON T.[object_id] = AC.[object_id]
INNER JOIN sys.[types] TY ON AC.[system_type_id] = TY.[system_type_id] AND AC.[user_type_id] = TY.[user_type_id]
OUTER APPLY (
SELECT FAC.[name] AS foreignKeyName
FROM sys.[foreign_keys] FKC 
INNER JOIN sys.[all_columns] FAC ON FAC.[object_id] = FKC.referenced_object_id 
AND FKC.parent_object_id = AC.[object_id] 
AND FAC.[name] = AC.[name]
AND FAC.is_identity = 1
) FAC
WHERE T.[is_ms_shipped] = 0
AND (OBJECT_SCHEMA_NAME(T.[object_id],DB_ID()) = 'dbo')
AND (T.[name] = ?)
ORDER BY T.[name], AC.[column_id]`

	rows, err := db.Query(queryStr, tableName)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	var metaDataRecords []TableSchemaData
	for rows.Next() {
		var field TableSchemaData
		err = rows.Scan(&field.ColumnName, &field.SystemDataType, &field.MaxLength, &field.IsNullable, &field.IsIdentity)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		metaDataRecords = append(metaDataRecords, field)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return metaDataRecords, nil
}

// InsertGeneratedData takes in a query string and a slice of values that will
// get swapped with their respective placeholders in the query string when the query is executed.
func InsertGeneratedData(query string, vals ...any) {
	db := connectToDb()
	defer db.Close()
	response, err := db.Exec(query, vals...)
	if err != nil {
		panic(err)
	}

	rows, err := response.RowsAffected()
	if err != nil {
		log.Printf("Error when getting rows affected: %s", err)
	}
	log.Printf("Rows affected: %d", rows)
}

type dbSetupData struct {
	dbUser, dbPass, dbName, dbHostName string
}

// getDBConnectionSetup Looks up the given Environment Variables to determine the information needed to connect to the database.
// Return dbSetupData filled with the necessary connection information (if the proper environment variables exist).
func getDBConnectionSetup() dbSetupData {
	var dbSetup dbSetupData
	dbCredsEnv, exists := os.LookupEnv("MSSQLServerDbCreds")
	if exists {
		splitCreds := strings.Split(dbCredsEnv, ":")
		dbSetup.dbUser = splitCreds[0]
		dbSetup.dbPass = splitCreds[1]
	}

	dbNameEnv, exists := os.LookupEnv("MSSQLServerDB")
	if exists {
		dbSetup.dbName = dbNameEnv
	}

	dbHostNameEnv, exists := os.LookupEnv("MSSQLServerHost")
	if exists {
		dbSetup.dbHostName = dbHostNameEnv
	}
	return dbSetup
}

// connectToDb will get the appropriate DB connection setup information and connect to the database.
// This returns db which is a reference to the database connection.
func connectToDb() *sql.DB {
	var db *sql.DB
	dbConnect := getDBConnectionSetup()
	query := url.Values{}
	query.Add("database", dbConnect.dbName)
	dbUrl := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(dbConnect.dbUser, dbConnect.dbPass),
		Host:     fmt.Sprintf("%v:%d", dbConnect.dbHostName, 1433),
		RawQuery: query.Encode(),
	}

	db, err := sql.Open("mssql", dbUrl.String())
	if err != nil {
		//parse or initialization error.
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}
