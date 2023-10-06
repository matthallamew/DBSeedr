package dbaccess

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

type TableSchemaData struct {
	ColumnName     string
	SystemDataType string
	MaxLength      int
	IsNullable     bool
	IsIdentity     bool
}

// GetTableSchemaData Connects to a database and queries the system tables to gather metadata about the given table.
// Return a TableSchemaData slice filled with metadata about the given table.
func GetTableSchemaData(tableName string) ([]TableSchemaData, error) {
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
	defer db.Close()
	if err != nil {
		//parse or initialization error.
		log.Fatal(err)
	}

	queryStr := `SELECT AC.[name] AS [ColumnName], TY.[name] AS SystemDataType, AC.[max_length] AS MaxLength,  AC.[is_nullable] AS IsNullable, AC.is_identity AS IsIdentity
FROM sys.[tables] AS T
INNER JOIN sys.[all_columns] AC ON T.[object_id] = AC.[object_id]
INNER JOIN sys.[types] TY ON AC.[system_type_id] = TY.[system_type_id] AND AC.[user_type_id] = TY.[user_type_id]
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
	var fieldData []TableSchemaData
	for rows.Next() {
		var fields TableSchemaData
		err = rows.Scan(&fields.ColumnName, &fields.SystemDataType, &fields.MaxLength, &fields.IsNullable, &fields.IsIdentity)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		fieldData = append(fieldData, fields)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return fieldData, nil
}

type dbSetupData struct {
	dbUser, dbPass, dbName, dbHostName string
}

// getDBConnectionSetup Looks up the given Environment Variables to determine the information needed to connect to the database.
// Return a string containing the username and a string containing the password, if they exist.
// Otherwise, return two empty strings.
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
