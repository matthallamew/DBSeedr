# DBSeedr
Seed a database table(s) with random data for UI testing

Given a database table name, get the metadata (data type, length, and nullability) for all of the fields. With this information, generate random data for insertion into the table. Use the maximum length of string/varchar type fields to determine how long the generated data should be.
