package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type Entry struct {
	Kategory     string
	Voicechannel string
	Textchannel  string
	Creator      string
}

var (
	databasepath string
)

func Init(path string) {
	databasepath = path
}

func CreateDatabase() {
	os.Remove(databasepath) // I delete the file to avoid duplicated records. SQLite is a file based database.

	log.Println("Creating sqlite-database.db...")
	file, err := os.Create(databasepath) // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite-database.db created")

	sqliteDatabase, _ := sql.Open("sqlite3", "./"+databasepath) // Open the created SQLite File
	defer sqliteDatabase.Close()                                // Defer Closing the database
	createTable(sqliteDatabase)                                 // Create Database Tables
}

func createTable(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE entry (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"kategorie" TEXT,
		"textchannel" TEXT,
		"voicechannel" TEXT,
		"creator" TEXT		
		);` // SQL Statement for Create Table

	statement, err := db.Prepare(createStudentTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
}

// We are passing db reference connection from main to our method with other parameters
func AddEntry(kategorie string, textchannel string, voicechannel string, creator string) {
	sqliteDatabase, _ := sql.Open("sqlite3", "./"+databasepath)                                                                      // Open the created SQLite File
	statement, err := sqliteDatabase.Prepare(`INSERT INTO entry(kategorie ,textchannel, voicechannel, creator) VALUES (?, ?, ?, ?)`) // Prepare statement. This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("%s, %s, %s, %s", kategorie, textchannel, voicechannel, creator)
	queryStatus, err := statement.Exec(kategorie, textchannel, voicechannel, creator)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(queryStatus.RowsAffected())
}

func GetAllEntrys() []Entry {
	sqliteDatabase, _ := sql.Open("sqlite3", "./"+databasepath) // Open the created SQLite File
	entries := []Entry{}
	row, err := sqliteDatabase.Query("SELECT * FROM entry")
	if err != nil {
		log.Fatal(err)
	}
	//defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id string
		var kategorie string
		var textchannel string
		var voicechannel string
		var creator string
		err = row.Scan(&id, &kategorie, &textchannel, &voicechannel, &creator)
		fmt.Println(err)
		entries = append(entries, Entry{kategorie, textchannel, voicechannel, creator})

	}

	return entries

}
