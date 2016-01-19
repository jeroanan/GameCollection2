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

  rawPlatformData, e := ioutil.ReadFile("masterfiles/platforms.json")

  if e!=nil {
    log.Fatal(e)
  }

  var mf PlatformsMasterFile

  e = json.Unmarshal(rawPlatformData, &mf)

  if e!=nil {
    log.Fatal(e)
  }

  return mf
}
