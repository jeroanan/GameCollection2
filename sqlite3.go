package main

import (
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
  "log"
)

func getDbTransaction(db *sql.DB) *sql.Tx {
  tx,err := db.Begin()

  if err!=nil {
    log.Fatal(err)
  }

  return tx
}

func prepareQuery(tx *sql.Tx, query string) *sql.Stmt {

  q, err := tx.Prepare(query)

  if err!=nil {
    log.Fatal(err)
  }

  return q
}

func getDbConnection() *sql.DB {
  db, err := sql.Open("sqlite3", getDatabaseFileLocation())

  if err!=nil {
    log.Fatal(err)
  }

  return db
}

func getDatabaseFileLocation() string {
  c, e := ReadConfig()

  if e!= nil {
    log.Fatal(e)
  }

  return c.GetDatabaseFileLocation()
}
