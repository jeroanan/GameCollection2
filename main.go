package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
)

func main() {
  fmt.Println("")
  defer fmt.Println()

  c, e := ReadConfig()

  if e!=nil {
    log.Print(e)
    return
  }

  InitializeDatabase("./games.db")

  http.HandleFunc("/", handler)

  httpPort := fmt.Sprintf(":%d", c.HttpPort)

  log.Printf("Waiting for connections on port %s", httpPort)
  e = http.ListenAndServe(httpPort, nil)

  if e!=nil {
    log.Fatal(e)
  }
}

func handler(w http.ResponseWriter, r *http.Request) {

  urlMap := map[string]string{
    "/": "index.html",
  }

  viewFile, ok := urlMap[string(r.URL.Path)]

  if ok {
    fileContent, err := ioutil.ReadFile(fmt.Sprintf("views/%s", viewFile))

    if err !=nil {
      log.Printf("Error reading %s: %s\n", viewFile, err)
      return
    }

    fmt.Fprintf(w, "%s", fileContent)
  }
}
