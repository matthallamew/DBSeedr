# DBSeedr
Seed a database table(s) with random data for UI testing

Given a database table name, get the metadata (example: data type, maximum length, etc) for all the fields in the table. 
With this information, generate random data for insertion into the table. The maximum length derived from the metadata
is used to determine how long the generated data should be so each field can be as full as possible.

Currently, this is program is tightly coupled to Microsoft SQL Server.