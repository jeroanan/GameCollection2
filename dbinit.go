package main

import (
  "database/sql"
  "fmt"
  _ "github.com/mattn/go-sqlite3"
  "log"
  "strings"
)

func InitializeDatabase(fileName string) {
  db, err := sql.Open("sqlite3", fileName)

  if err!=nil {
    fmt.Println(err)
    return
  }

  defer func() {
    db.Close()
  }()

  createTables(db)
  populateMasterTables(db)
}

func createTables(db *sql.DB) {
  createGenreTable(db)
  createHardwareTypeTable(db)
  createPlatformTable(db)
  createHardwareTable(db)
  createGameTable(db)
}
func createGenreTable(db *sql.DB) {

  stmt := "CREATE TABLE Genre (Name TEXT NOT NULL PRIMARY KEY, Description TEXT);"
  createTableIfNotExist(db, "Genre", stmt)
}

func createHardwareTypeTable(db *sql.DB) {
  stmt := "CREATE TABLE HardwareType(Name TEXT NOT NULL PRIMARY KEY, Description TEXT);"
  createTableIfNotExist(db, "HardwareType", stmt)
}

func createPlatformTable(db *sql.DB) {
  stmt := "CREATE TABLE Platform(Name TEXT NOT NULL PRIMARY KEY, Description TEXT);"
  createTableIfNotExist(db, "Platform", stmt)
}

func createHardwareTable(db *sql.DB) {

  cols := [...]string{
    "Name TEXT NOT NULL",
    "HardwareType INTEGER References HardwareType(id)",
    "Platform TEXT References Platform(id)",
    "NumberOwned INTEGER",
    "NumberBoxed INTEGER",
    "Notes TEXT"}
  colString := strings.Join(cols[:], ", ")
  stmt := fmt.Sprintf("CREATE TABLE Hardware(%s)", colString)

  createTableIfNotExist(db, "Hardware", stmt)
}

func createGameTable(db *sql.DB) {
  cols := [...]string{
    "Title TEXT NOT NULL",
    "Genre TEXT References Genre(Name)",
    "Platform TEXT References Platform(Name)",
    "NumberOwned INTEGER",
    "NumberBoxed INTEGER",
    "NumberOfManuals INTEGER",
    "DatePurchased TEXT",
    "ApproximatePurchaseDate INTEGER",
    "Notes TEXT"}
  colString := strings.Join(cols[:], ", ")
  stmt := fmt.Sprintf("CREATE TABLE Game(%s)", colString)

  createTableIfNotExist(db, "Game", stmt)
}

func createTableIfNotExist(db *sql.DB, tableName string, createStatement string) {

  if tableExists(tableName, db) {
    log.Printf("%s table already exists -- no need to recreate", tableName)
    return
  }

  log.Printf("Going to create %s table", tableName)

  _, err := db.Exec(createStatement)

  if err!=nil {
    log.Fatal(err)
  }

  log.Printf("%s table created successfully", tableName)
}

func tableExists(tableName string, db *sql.DB) bool {

  log.Printf("Checking for existence of %s table", tableName)

  tx := getDbTransaction(db)

  stmt := prepareQuery(tx, "SELECT COUNT(name) FROM sqlite_master WHERE type='table' AND name=?")

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

func populateMasterTables(db *sql.DB) {
  populatePlatformsTable(db)
  populateGenresTable(db)
  populateHardwareTypeTable(db)
}

func populatePlatformsTable(db *sql.DB) {

  log.Print("Populating Platform table from masterfiles/platforms.json")

  tx := getDbTransaction(db)

  existQuery := prepareQuery(tx, "SELECT COUNT(Name) FROM Platform WHERE Name=?")
  defer existQuery.Close()

  platformData := LoadPlatformsMasterFile()

  var platforms []NameDescription

  for _, v := range(platformData.Platforms) {
    platforms = append(platforms, v)
  }

  platformsToAdd := getNonExistentNameDescriptionRows(existQuery, platforms)

  insertStmt := prepareQuery(tx, "INSERT INTO Platform (Name, Description) VALUES (?, ?)")
  defer insertStmt.Close()
  insertNameDescriptionTable(insertStmt, platformsToAdd, "platform")

  tx.Commit()
}

func populateGenresTable(db *sql.DB) {

  log.Print("Populating Genres table from masterfiles/genres.json")

  tx := getDbTransaction(db)

  existQuery := prepareQuery(tx, "SELECT COUNT(Name) FROM Genre WHERE Name=?")
  defer existQuery.Close()

  genreData := LoadGenresMasterFile()

  var genres []NameDescription

  for _, v := range(genreData.Genres) {
    genres = append(genres, v)
  }

  genresToAdd := getNonExistentNameDescriptionRows(existQuery, genres)

  insertStmt := prepareQuery(tx, "INSERT INTO Genre (Name, Description) VALUES (?, ?)")
  defer insertStmt.Close()
  insertNameDescriptionTable(insertStmt, genresToAdd, "genre")

  tx.Commit()
}

func populateHardwareTypeTable(db *sql.DB) {

  log.Print("Populating HardwareType table from masterfiles/hardwaretypes.json")

  tx := getDbTransaction(db)

  existQuery := prepareQuery(tx, "SELECT COUNT(Name) FROM HardwareType WHERE Name=?")
  defer existQuery.Close()

  hardwareTypeData := LoadHardwareTypesMasterFile()

  var hardwareTypes []NameDescription

  for _, v := range(hardwareTypeData.HardwareTypes) {
    hardwareTypes = append(hardwareTypes, v)
  }

  hardwareTypesToAdd := getNonExistentNameDescriptionRows(existQuery, hardwareTypes)

  insertStmt := prepareQuery(tx, "INSERT INTO HardwareType (Name, Description) VALUES (?, ?)")
  defer insertStmt.Close()
  insertNameDescriptionTable(insertStmt, hardwareTypesToAdd, "hardware type")

  tx.Commit()
}

func getNonExistentNameDescriptionRows(existQuery *sql.Stmt, items []NameDescription) []NameDescription {

  var out []NameDescription

  for _, item := range(items) {

    rs, err := existQuery.Query(item.GetName())

    if err!=nil {
      log.Fatal(err)
    }

    var count int
    for rs.Next() {
      rs.Scan(&count)
    }

    if count==0 {
      out = append(out, item)
    }
  }

  return out
}

func insertNameDescriptionTable(insertStmt *sql.Stmt, items []NameDescription, itemName string) {

  for _, item := range(items) {
    log.Printf("Inserting %s %s", itemName, item)
    _, err := insertStmt.Exec(item.GetName(), item.GetDescription())

    if err!=nil {
      log.Fatal(err)
    }
  }
}

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
