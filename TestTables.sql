

-- Random, simple tables that can be used for testing the generator

IF NOT EXISTS(SELECT * FROM sys.tables WHERE [name]='Addr')
BEGIN
	CREATE TABLE Addr(
	 addressID INT IDENTITY(1,1) NOT NULL,
	 city VARCHAR(100) NULL,
	[state] VARCHAR(10) NULL,
	[zip] VARCHAR(20) NULL
	CONSTRAINT [PK_Addr] PRIMARY KEY CLUSTERED([addressID] ASC));

END


IF NOT EXISTS(SELECT * FROM sys.tables WHERE [name]='Person')
BEGIN
	CREATE TABLE Person(
	personID INT IDENTITY(1,1) NOT NULL,
	CONSTRAINT [PK_Person] PRIMARY KEY CLUSTERED([personID] ASC),
	addressID INT NULL,
	firstName VARCHAR (100) NULL,
	lastName VARCHAR (100)  NULL,
	age  INT                NULL,
	comment  VARCHAR(8000) NULL,
	lifeStory  VARCHAR(MAX) NULL,
	active BIT NOT NULL CONSTRAINT DF_Person_active DEFAULT 1,
	salary   DECIMAL (18, 2),
	CONSTRAINT [FK_Person_Address] FOREIGN KEY ([addressID]) REFERENCES Addr ([addressID]) ON DELETE SET NULL)
END
