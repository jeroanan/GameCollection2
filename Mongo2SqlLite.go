package main

import (
  "fmt"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "log"
  "strconv"
)

type MongoGame struct {
  _Game__id string

  _Game__title string
}

func MigrateMongoToSqlLite() {
  session, err := mgo.Dial("127.0.0.1")

  log.Print("Migrating from old MongoDb datasource to Sqlite3...")

  if err != nil {
    log.Fatal(err)
  }

  defer session.Close()

  MigratePlatforms(session)
  MigrateGenres(session)
  MigrateHardwareTypes(session)
  MigrateGames(session)
}

func MigratePlatforms(session *mgo.Session) {
  result := mongoGetCollectionContents("platforms", session)

  var ps []Platform
  for _, doc := range(result) {
    var p Platform
    p.Name = doc["_Platform__name"]
    p.Description = doc["_Platform__description"]
    ps = append(ps, p)
  }

  for _, p := range(ps) {
    r, err := GetPlatformByName(p.Name)

    if err!=nil {
      log.Fatalf("Error while retrieving platform %s: %s", p, err)
    }

    if r.Name!=p.Name {
      log.Printf("Could not find platform %s; adding it.", p)
      err = AddPlatform(p)

      if err!=nil {
        log.Printf("Error while adding platform %s: %s", p, err)
      }
    }
  }
}

func MigrateGenres(session *mgo.Session) {
  result := mongoGetCollectionContents("genres", session)

  var gs []Genre

  for _, doc := range(result) {
    g := Genre{}
    g.Name = doc["_Genre__name"]
    g.Description = doc["_Genre__description"]
    gs = append(gs, g)
  }

  for _, g := range(gs) {
    r, err := GetGenreByName(g.Name)

    if err!=nil {
      log.Fatalf("Error while retrieving genre %s: %s", g, err)
    }

    if r.Name!=g.Name {
      log.Printf("Could not find genre %s; adding it.", g)
      err = AddGenre(g)

      if err!=nil {
        log.Printf("Error while adding genre %s: %s", g, err)
      }
    }
  }
}

func MigrateHardwareTypes(session *mgo.Session) {
  result := mongoGetCollectionContents("hardware_types", session)

  var hts []HardwareType

  for _, doc := range(result) {
    ht := HardwareType{}
    ht.Name = doc["_HardwareType__name"]
    ht.Description = doc["_HardwareType__description"]
    hts = append(hts, ht)
  }

  for _, ht := range(hts) {
    r, err := GetHardwareTypeByName(ht.Name)

    if err!=nil {
      log.Fatalf("Error while retrieving hardware type %s: %s", ht, err)
    }

    if r.Name!=ht.Name {
      log.Printf("Could not find hardware type %s; adding it.", ht)
      err = AddHardwareType(ht)

      if err!=nil {
        log.Printf("Error while adding hardware type %s: %s", ht, err)
      }
    }
  }
}

func MigrateGames(session *mgo.Session) {
  result := mongoGetCollectionContents("games", session)

  var gs []Game

  for _, doc := range(result) {
    g := Game{}
    g.Title = doc["_Game__title"]

    mongoGenreName := doc["_Game__genre"]
    genre, err := GetGenreByName(mongoGenreName)

    if err!=nil {
      log.Printf("Error while retrieving genre with name %s. Aborting import of %s: %s", mongoGenreName, g, err)
      continue
    }

    g.Genre = fmt.Sprintf("%d", genre.RowId)

    mongoPlatformName := doc["_Game__platform"]
    platform, err := GetPlatformByName(mongoPlatformName)

    if err!= nil {
      log.Printf("Error while retrieving platform with name %s. Aborting import of %s: %s", mongoPlatformName, g, err)
      continue
    }

    g.Platform = fmt.Sprintf("%d", platform.RowId)

    mongoNumberOfCopies := doc["_Game__num_copies"]
    g.NumberOwned, err = strconv.Atoi(mongoNumberOfCopies)

    if err!=nil {
      log.Printf("Error while parsing number of copies %s. Aborting import of %s: %s", mongoNumberOfCopies, g, err)
      continue
    }

    mongoNumBoxed := doc["_Game__num_boxed"]
    g.NumberBoxed, err = strconv.Atoi(mongoNumBoxed)

    if err!=nil {
      log.Printf("Error while parsing number boxed %s. Aborting import of %s: %s", mongoNumBoxed, g, err)
      continue
    }

    mongoNumberOfManuals := doc["_Game__num_manuals"]
    g.NumberOfManuals, err = strconv.Atoi(mongoNumberOfManuals)

    if err!=nil {
      log.Printf("Error while parsing number of manuals %s. Aborting import of %s: %s", mongoNumberOfManuals, g, err)
    }

    g.DatePurchased = doc["_Game__date_purchased"]

    mongoApproximatePurchaseDate := doc["_Game__approximate_date_purchased"]

    if mongoApproximatePurchaseDate=="" {
      mongoApproximatePurchaseDate = "False"
    }

    g.ApproximatePurchaseDate, err = strconv.ParseBool(mongoApproximatePurchaseDate)
    
    if err!=nil {
      log.Printf("Error while parsing approximate purchase date %s. Aborting import of %s: %s",
        mongoApproximatePurchaseDate, g, err)
      continue
    }

    g.Notes = doc["_Game__notes"]

    gs = append(gs, g)
  }

  for _, g := range(gs) {
    existingGames, err := GetGamesByNameAndPlatform(g)

    if err!=nil {
      log.Printf("Error while checking if %s exists. Aborting import of %s : %s", g, g, err)
      continue
    }

    if len(existingGames)==0 {
      log.Printf("%s not found; adding.", g)
      _, err := SaveGame(g)

      if err!=nil {
        log.Printf("Error while import game %s. Aborting import of %s: %s", g, g, err)
        continue
      }
    }
  }
}

func mongoGetCollectionContents(collectionName string, session *mgo.Session) []map[string]string {
  c := session.DB("GamesCollection").C(collectionName)

  var result []map[string]string
  err := c.Find(bson.M{}).All(&result)

  if err!=nil {
    log.Fatal(err)
  }

  return result
}
