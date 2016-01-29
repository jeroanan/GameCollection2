package main

import (
  "encoding/json"
  "fmt"
  "log"
)

func GetAllPlatforms(requestBody []byte) (string, error) {

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
    log.Printf("Error while marshalling platform json: %s", err)
    return "", err
  }

  return fmt.Sprintf("%s", j), err
}


func GetAllGenres(requestBody[] byte) (string, error) {

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
    log.Print("Error while marshalling genres: %s", err)
    return "", err
  }

  return fmt.Sprintf("%s", j), err
}

func GetAllHardwareTypes(requestBody[] byte) (string, error) {

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
    log.Printf("Error while marshalling hardware types json: %s", err)
    return "", err
  }

  return fmt.Sprintf("%s", j), err
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

func SaveGameFromJson(requestBody []byte) (string, error) {
  g := Game{}
  err := json.Unmarshal(requestBody, &g)

  if err!=nil {
    return "", err
  }

  return SaveGame(g)
}

func SaveGame(g Game) (string, error) {
  /* At this point I could check to see if an identical game has been saved (e.g. same title/platform) and raise
   * an error if that is the case, but at this point I won't do that. It would be possible to have this situation. For
   * instance there are two different versions of "Tetris" for the NES (by Nintendo and Tengen).
   *
   * Maybe at a later date we could have something in place that will go back and ask the user for confirmation if this
   * is the case.
   */

  j, err := GetGameById(fmt.Sprintf("%d",g.RowId))

  if err!=nil {
    log.Printf("Error while retrieving game %s during save operation: %s", g, err)
    return "", err
  }

  existingGame := Game{}
  err = json.Unmarshal([]byte(j), &existingGame)

  if err!=nil {
    log.Printf("Error unmrashaling json %s: %s", j, err)
    return "", err
  }

  if existingGame.RowId==g.RowId && g.RowId!=0 {
    updateGame(g)
    return "", err
  }

  return addGame(g)
}

func addGame(g Game) (string, error) {

  log.Printf("Adding new game %s", g)

  insertString := "INSERT INTO Game (Title, Genre, Platform, NumberOwned, NumberBoxed, NumberOfManuals, DatePurchased, "
  insertString += "ApproximatePurchaseDate, Notes) "
  insertString += "VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"

  stmt, closerFunc := GetQuery(insertString)
  defer closerFunc()

  approximatePurchaseDate := 0
  if g.ApproximatePurchaseDate {
    approximatePurchaseDate = 1
  }

  _, err := stmt.Exec(g.Title, g.Genre, g.Platform, g.NumberOwned, g.NumberBoxed, g.NumberOfManuals, g.DatePurchased,
    approximatePurchaseDate, g.Notes)

  if err!=nil {
    log.Printf("Error while adding new game %s: %s", g, err)
  }

  return "", err
}

func updateGame(g Game) (string, error) {

  log.Printf("Updating game %s with row id of %d", g, g.RowId)

  updateString := "UPDATE GAME "
  updateString += "SET Title=?, Genre=?, Platform=?, NumberOwned=?, NumberBoxed=?, NumberOfManuals=?, DatePurchased=?, "
  updateString += "ApproximatePurchaseDate=?, Notes=? "
  updateString += "WHERE RowId=?"

  stmt, closerFunc := GetQuery(updateString)
  defer closerFunc()

  approximatePurchaseDate := 0
  if g.ApproximatePurchaseDate {
    approximatePurchaseDate = 1
  }

  _, err := stmt.Exec(g.Title, g.Genre, g.Platform, g.NumberOwned, g.NumberBoxed, g.NumberOfManuals, g.DatePurchased,
    approximatePurchaseDate, g.Notes, g.RowId)

  if err!=nil {
    log.Printf("Error while updating existing game %s: %s", g, err)
  }

  return "", err
}

func GetAllGames(requestBody []byte) (string, error) {

  queryString := "SELECT g.RowId, g.Title, g.Genre, COALESCE(p.Name, ''), g.NumberOwned, g.NumberBoxed, "
  queryString += "g.NumberOfManuals, g.DatePurchased, g.ApproximatePurchaseDate, g.Notes "
  queryString += "FROM Game g "
  queryString += "LEFT JOIN Platform p ON g.Platform=p.RowId "
  queryString += "LEFT JOIN Genre gen on g.Genre=gen.RowId "

  var gameList GameList
  var err error

  stmt, closerFunc := GetQuery(queryString)
  defer closerFunc()

  rs, err := stmt.Query()

  if err!=nil {
    log.Printf("Error while running sql query to retrieve games: \n%s\n%s", queryString, err)
    return "", err
  }

  for rs.Next() {
    g := Game{}

    err = rs.Scan(&g.RowId, &g.Title, &g.Genre, &g.Platform, &g.NumberOwned, &g.NumberBoxed, &g.NumberOfManuals,
      &g.DatePurchased, &g.ApproximatePurchaseDate, &g.Notes)

    if (err!=nil) {
      log.Printf("Error while scanning games resultset: %s", err)
      return "", err
    }

    gameList.Games = append(gameList.Games, g)
  }

  j, err := json.Marshal(gameList)

  if err!=nil {
    log.Printf("Error while unmarshalling games json: %s", err)
  }

  return string(j), err
}

func GetGamesByNameAndPlatform(g Game) ([]Game, error) {
  queryString := "SELECT g.RowId, g.Title, g.Genre, COALESCE(p.Name, ''), g.NumberOwned, g.NumberBoxed, "
  queryString += "g.NumberOfManuals, g.DatePurchased, g.ApproximatePurchaseDate, g.Notes "
  queryString += "FROM Game g "
  queryString += "LEFT JOIN Platform p ON g.Platform=p.RowId "
  queryString += "LEFT JOIN Genre gen on g.Genre=gen.RowId "
  queryString += "WHERE g.Title = ? AND g.Platform = ?"

  stmt, closerFunc := GetQuery(queryString)
  defer closerFunc()

  var gs []Game

  rs, err := stmt.Query(g.Title, g.Platform)

  if err!=nil {
    return gs, err
  }

  for rs.Next() {
    g := Game{}

    err = rs.Scan(&g.RowId, &g.Title, &g.Genre, &g.Platform, &g.NumberOwned, &g.NumberBoxed, &g.NumberOfManuals,
      &g.DatePurchased, &g.ApproximatePurchaseDate, &g.Notes)

    if err!=nil {
      return gs, err
    }

    gs = append(gs, g)
  }
  return gs, err
}

type GameResult struct {
  Game
  PlatformId string
  GenreId string
}

func GetGameById(gameId string) (string, error) {

  var err error

  queryString := "SELECT g.RowId, g.Title, g.Genre, COALESCE(gen.Name, ''), g.Platform, COALESCE(p.Name, ''), "
  queryString += "g.NumberOwned, g.NumberBoxed, g.NumberOfManuals, g.DatePurchased, g.ApproximatePurchaseDate, g.Notes "
  queryString += "FROM Game g "
  queryString += "LEFT JOIN Platform p ON g.Platform=p.RowId "
  queryString += "LEFT JOIN Genre gen on g.Genre=gen.RowId "
  queryString += "WHERE g.RowId=? "

  stmt, closerFunc := GetQuery(queryString)
  defer closerFunc()

  if err!=nil {
    return "", err
  }

  rs, err := stmt.Query(gameId)

  if err!=nil {
    return "", err
  }

  g := GameResult{}
  for rs.Next() {
    err = rs.Scan(&g.RowId, &g.Title, &g.GenreId, &g.Genre, &g.PlatformId, &g.Platform, &g.NumberOwned,
      &g.NumberBoxed, &g.NumberOfManuals, &g.DatePurchased, &g.ApproximatePurchaseDate, &g.Notes)

    if err!=nil {
      return "", err
    }
  }

  j, err := json.Marshal(g)
  return string(j), err
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
    err = rs.Scan(&ht.RowId, &ht.Name, &ht.Description)

    if err!=nil {
      return ht, err
    }
  }

  return ht, err
}

func AddHardwareType(ht HardwareType) error {
  var err error

  queryString := "INSERT INTO HardwareType (Name, Description) VALUES (?, ?)"
  stmt, closerFunc := GetQuery(queryString)
  defer closerFunc()

  _, err = stmt.Exec(ht.Name, ht.Description)
  return err
}
