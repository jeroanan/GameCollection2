package main

import (
  "encoding/json"
  "io/ioutil"
  "log"
)

type PlatformsMasterFile struct {
  Platforms []Platform
}

func LoadPlatformsMasterFile() PlatformsMasterFile {

  log.Print("Reading platforms master file")

  rawPlatformData := loadMasterFileString("masterfiles/platforms.json")

  var mf PlatformsMasterFile

  e := json.Unmarshal(rawPlatformData, &mf)

  if e!=nil {
    log.Fatal(e)
  }

  return mf
}

type GenresMasterFile struct {
  Genres []Genre
}
func LoadGenresMasterFile() GenresMasterFile {

  log.Print("Reading genres master file")

  rawGenreData := loadMasterFileString("masterfiles/genres.json")

  var mf GenresMasterFile

  e := json.Unmarshal(rawGenreData, &mf)

  if e!=nil {
    log.Fatal(e)
  }

  return mf
}

type HardwareTypesMasterFile struct {
  HardwareTypes []HardwareType
}
func LoadHardwareTypesMasterFile() HardwareTypesMasterFile {

  log.Print("Reading hardware types master file")

  rawHardwareTypeData := loadMasterFileString("masterfiles/hardwaretypes.json")

  var mf HardwareTypesMasterFile

  e := json.Unmarshal(rawHardwareTypeData, &mf)

  if e!=nil {
    log.Fatal(e)
  }

  return mf
}

func loadMasterFileString(path string) []byte {

  rawData, e := ioutil.ReadFile(path)

  if (e!=nil) {
    log.Fatal(e)
  }

  return rawData
}
