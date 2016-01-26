package main

import (
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "log"
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
      log.Fatal(err)
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
      log.Fatal(err)
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
      log.Fatal(err)
    }

    if r.Name!=ht.Name {
      log.Printf("Could not find hardware type %s; adding it.", ht)
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
