package main

import (
  "encoding/json"
  "fmt"
  "log"
)

func GetAllPlatforms() string {

  queryString := "SELECT RowId, Name, Description FROM Platform"

  ndrs := GetAllNameDescriptionRows(queryString)

  var pmf PlatformsMasterFile

  for _, v := range(ndrs) {
    var p Platform
    p.RowId = v.RowId
    p.Name = v.Name
    p.Description = v.Description

    pmf.Platforms = append(pmf.Platforms, p)
  }

  j, err := json.Marshal(pmf)

  if (err!=nil) {
    log.Fatal(err)
  }

  return fmt.Sprintf("%s", j)
}


func GetAllGenres() string {

  queryString := "SELECT RowId, Name, Description FROM Genre"

  ndrs := GetAllNameDescriptionRows(queryString)

  var gmf GenresMasterFile

  for _, v := range(ndrs) {
    var genre Genre
    genre.RowId = v.RowId
    genre.Name = v.Name
    genre.Description = v.Description

    gmf.Genres = append(gmf.Genres, genre)
  }

  j, err := json.Marshal(gmf)

  if err!=nil {
    log.Fatal(err)
  }

  return fmt.Sprintf("%s", j)
}

func GetAllHardwareTypes() string {

  queryString := "SELECT RowId, Name, Description FROM HardwareType"

  ndrs := GetAllNameDescriptionRows(queryString)

  var htmf HardwareTypesMasterFile

  for _, v := range(ndrs) {
    var ht HardwareType
    ht.RowId = v.RowId
    ht.Name = v.Name
    ht.Description = v.Description

    htmf.HardwareTypes = append(htmf.HardwareTypes, ht)
  }

  j, err := json.Marshal(htmf)

  if err!=nil {
    log.Fatal(err)
  }

  return fmt.Sprintf("%s", j)
}

func GetAllNameDescriptionRows(queryString string) []NameDescriptionTable {

  db := getDbConnection()
  defer db.Close()

  tx := getDbTransaction(db)
  defer tx.Commit()

  stmt := prepareQuery(tx, queryString)
  defer stmt.Close()

  rs, err := stmt.Query()

  if err!=nil {
    log.Fatal(err)
  }

  var items []NameDescriptionTable

  for rs.Next() {
    var item NameDescriptionTable

    rs.Scan(&item.RowId, &item.Name, &item.Description)
    items = append(items, item)
  }

  return items
}

func SaveGame(requestBody []byte) (string, error) {
  var g Game
  err := json.Unmarshal(requestBody, &g)

  if err!=nil {
    return "", err
  }

  /* At this point I could check to see if an identical game has been saved (e.g. same title/platform) and raise
   * an error if that is the case, but at this point I won't do that. It would be possible to have this situation. For
   * instance there are two different versions of "Tetris" for the NES (by Nintendo and Tengen).
   *
   * Maybe at a later date we could have something in place that will go back and ask the user for confirmation if this
   * is the case.
   */

  insertString := "INSERT INTO Game (Title, Genre, Platform, NumberOwned, NumberBoxed, NumberOfManuals, DatePurchased, "
  insertString += "ApproximatePurchaseDate, Notes) "
  insertString += "VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"

  db := getDbConnection()
  defer db.Close()

  tx := getDbTransaction(db)
  defer tx.Commit()

  stmt := prepareQuery(tx, insertString)
  defer stmt.Close()

  stmt.Exec(g.Title, g.Genre, g.Platform, g.NumberOwned, g.NumberBoxed, g.NumberOfManuals, g.DatePurchased,
    g.ApproximatePurchaseDate, g.Notes)
  return "", err
}

func GetAllGames() (string) {
  db := getDbConnection()
  queryString := "SELECT g.RowId, g.Title, g.Genre, COALESCE(p.Name, ''), g.NumberOwned, g.NumberBoxed, "
  queryString += "g.NumberOfManuals, g.DatePurchased, g.ApproximatePurchaseDate, g.Notes "
  queryString += "FROM Game g "
  queryString += "LEFT JOIN Platform p ON g.Platform=p.RowId "
  queryString += "LEFT JOIN Genre gen on g.Genre=gen.RowId "

  var gameList GameList
  var err error

  rs, err := db.Query(queryString)
  if err!=nil {
    log.Print(err)
    return ""
  }

  for rs.Next() {
    var g Game

    err = rs.Scan(&g.RowId, &g.Title, &g.Genre, &g.Platform, &g.NumberOwned, &g.NumberBoxed, &g.NumberOfManuals,
      &g.DatePurchased, &g.ApproximatePurchaseDate, &g.Notes)

    if (err!=nil) {
      log.Print(err)
    }

    gameList.Games = append(gameList.Games, g)
  }

  j, _ := json.Marshal(gameList)
  return string(j)
}

func DeleteGame(requestBody []byte) (string, error) {
  var err error
  var g Game

  err = json.Unmarshal(requestBody, &g)

  if err!=nil {
    return "", err
  }

  deleteString := "DELETE FROM GAME WHERE RowId=?"
  stmt, closerFunc := GetQuery(deleteString)
  defer closerFunc()

  _, err = stmt.Exec(g.RowId)
  return "", err
}

func GetPlatformByName(name string) (Platform, error) {
  var err error
  var p Platform

  queryString := "SELECT RowId, * FROM Platform WHERE Name=?"
  stmt, closerFunc := GetQuery(queryString)
  defer closerFunc()

  rs, err := stmt.Query(name)

  if err!=nil {
    return p, err
  }

  for rs.Next() {
    err = rs.Scan(&p.RowId, &p.Name, &p.Description)

    if err!=nil {
      return p, err
    }
  }

  return p, err
}

func AddPlatform(platform Platform) error {
  var err error

  insertString := "INSERT INTO Platform (Name, Description) VALUES (?, ?)"
  stmt, closerFunc := GetQuery(insertString)
  defer closerFunc()

  _, err = stmt.Exec(platform.Name, platform.Description)
  return err
}


func GetGenreByName(name string) (Genre, error) {
  var err error
  var g Genre

  queryString := "SELECT RowId, * FROM Genre WHERE Name=?"
  stmt, closerFunc := GetQuery(queryString)
  defer closerFunc()

  rs, err := stmt.Query(name)

  if err!=nil {
    return g, err
  }

  for rs.Next() {
    err = rs.Scan(&g.RowId, &g.Name, &g.Description)

    if err!=nil {
      return g, err
    }
  }
  return g, err
}

func AddGenre(genre Genre) error {
  var err error

  queryString := "INSERT INTO Genre (Name, Description) VALUES (?, ?)"
  stmt, closerFunc := GetQuery(queryString)
  defer closerFunc()

  _, err = stmt.Exec(genre.Name, genre.Description)
  return err
}

func GetHardwareTypeByName(name string) (HardwareType, error) {
  var err error
  var ht HardwareType

  queryString := "SELECT RowId, * From HardwareType WHERE Name=?"
  stmt, closerFunc := GetQuery(queryString)
  defer closerFunc()

  rs, err := stmt.Query(name)

  if err!=nil {
    return ht, err
  }

  for rs.Next() {
    err = rs.Scan(&ht.Name, ht.Description)

    if err!=nil {
      return ht, err
    }
  }

  return ht, err
}
