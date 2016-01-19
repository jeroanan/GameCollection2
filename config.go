package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
)

type Config struct {
  HttpPort int
}

func ReadConfig() (Config, error) {
  var c Config
  var e error

  rawConfigData, e := ioutil.ReadFile("config.json")

  if e!=nil {
    e = nil

    e = writeDefaultConfig()

    if e != nil {
      errMsg := fmt.Sprintf("config.json was not found. Addtionally the following error was encountered while trying to generate a default config.json: %s", e)
      log.Fatal(errMsg)
    }

    log.Fatal("config.json not found. A default one has been generated. Please updated its values appropriately.")

  }

  e = json.Unmarshal(rawConfigData, &c)

  if (e!=nil) {
    log.Fatal(e)
  }

  e = verifyConfig(c)
  return c, e
}

func writeDefaultConfig() error {
  var c Config
  var e error

  c.HttpPort = -1

  j, e := json.Marshal(&c)

  if e != nil {
    return e
  }

  ioutil.WriteFile("config.json", j, 0664)
  return e
}

func verifyConfig(c Config) error {

  var e error

  if c.HttpPort == -1 {
    log.Fatal(fmt.Sprintf("Change the defaults supplied in config.json (HttpPort is still -1)"))
  }

  return e
}
