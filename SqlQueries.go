package main

import (
  "encoding/json"
  "fmt"
  "log"
)

func GetAllPlatforms() string {

  queryString := "SELECT Name, Description FROM Platform"

  ndrs := GetAllNameDescriptionRows(queryString)

  var pmf PlatformsMasterFile

  for _, v := range(ndrs) {
    var p Platform
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

  queryString := "SELECT Name, Description FROM Genre"

  ndrs := GetAllNameDescriptionRows(queryString)

  var gmf GenresMasterFile

  for _, v := range(ndrs) {
    var genre Genre
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

  queryString := "SELECT Name, Description FROM HardwareType"

  ndrs := GetAllNameDescriptionRows(queryString)

  var htmf HardwareTypesMasterFile

  for _, v := range(ndrs) {
    var ht HardwareType
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

    rs.Scan(&item.Name, &item.Description)
    items = append(items, item)
  }

  return items
}

func SaveGame(requestBody string) string {
  return ""
}
