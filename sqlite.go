package main

import (
  "database/sql"
  "fmt"
  _ "github.com/mattn/go-sqlite3"
  "log"
)


func InitializeDatabase(fileName string) {
  db, err := sql.Open("sqlite3", fileName)

  if err!=nil {
    fmt.Println(err)
    return
  }

  createGenreTable(db)
  createHardwareTypeTable(db)

  defer func() {
    db.Close()
  }()
}

func createGenreTable(db *sql.DB) {

  stmt := "CREATE TABLE Genre (Name TEXT NOT NULL PRIMARY KEY, Description TEXT);"
  createTableIfNotExist(db, "Genre", stmt)
}

func createHardwareTypeTable(db *sql.DB) {
    stmt := "CREATE TABLE HardwareType(Name TEXT NOT NULL PRIMARY KEY, Description TEXT);"
    createTableIfNotExist(db, "HardwareType", stmt)
}

func createTableIfNotExist(db *sql.DB, tableName string, createStatement string) {

  if !tableExists(tableName, db) {
    log.Printf("Going to create %s table", tableName)

    _, err := db.Exec(createStatement)

    if err!=nil {
      log.Fatal(err)
    }

    log.Printf("%s table created successfully", tableName)
    return
  }

  log.Printf("%s table already exists -- no need to recreate", tableName)
}

func tableExists(tableName string, db *sql.DB) bool {

  log.Printf("Checking for existence of %s table", tableName)

  tx, err := db.Begin()

  if (err!=nil) {
    log.Fatal(err)
  }

  stmt, err := tx.Prepare("SELECT COUNT(name) FROM sqlite_master WHERE type='table' AND name=?")

  if (err!=nil) {
    log.Fatal(err)
  }

  defer stmt.Close()

  rs, err := stmt.Query(tableName)

  if err!= nil {
    log.Fatal(err)
  }

  var count int
  for rs.Next() {
    rs.Scan(&count)
  }

  tx.Commit()
  return count > 0
}
