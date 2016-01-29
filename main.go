package main

import (
  "fmt"
  "io"
  "io/ioutil"
  "log"
  "net/http"
  "strings"
)

func main() {
  fmt.Println("")
  defer fmt.Println()

  c, e := ReadConfig()

  if e!=nil {
    log.Print(e)
    return
  }

  InitializeDatabase(c.GetDatabaseFileLocation())
  //MigrateMongoToSqlLite()

  http.HandleFunc("/", handler)
  http.HandleFunc("/view/", viewHandler)
  http.HandleFunc("/json/", jsonHandler)
  http.HandleFunc("/js/", jsHandler)
  http.HandleFunc("/css/", cssHandler)

  httpPort := fmt.Sprintf(":%d", c.HttpPort)

  log.Printf("Waiting for connections on port %s", httpPort)
  e = http.ListenAndServe(httpPort, nil)

  if e!=nil {
    log.Fatal(e)
  }
}

func handler(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "/view", http.StatusFound)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {

  var routes = map[string]string {
    "/view/": "index.html",
    "/view/addgame": "addgame.html",
    "/view/allgames": "allgames.html",
    "/view/editgame": "editgame.html",
  }

  htmlFilePath, ok := routes[r.URL.Path]

  if ok {
    fileContent, err := ioutil.ReadFile(fmt.Sprintf("views/%s", htmlFilePath))

    if err!=nil {
      log.Print(err)
      return
    }

    fmt.Fprintf(w, "%s", fileContent)
    return
  }
  log.Print(r.URL.Path)
  http.NotFound(w, r)
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {

  requestBody := getRequestBody(r)

  urlToLower := strings.ToLower(r.URL.Path)

  if strings.HasPrefix(urlToLower, "/json/getgame/") {
    urlSplit := strings.Split(urlToLower, "/")
    gameId := urlSplit[len(urlSplit)-1]
    responseString, err := GetGameById(gameId)

    if err!=nil {
      http.Error(w, fmt.Sprintf("%s", err), 500)
      return
    }
    fmt.Fprintf(w, responseString)
    return
  }

  var dispatchFunc func([]byte) (string, error)

  switch urlToLower {
  case "/json/savegame":
    dispatchFunc = SaveGameFromJson
  case "/json/deletegame":
    dispatchFunc = DeleteGame
  case "/json/getplatforms":
    dispatchFunc = GetAllPlatforms
  case "/json/getgenres":
    dispatchFunc = GetAllGenres
  case "/json/gethardwaretypes":
    dispatchFunc = GetAllHardwareTypes
  case "/json/getgames":
    dispatchFunc = GetAllGames
  default:
    http.NotFound(w, r)
    return
  }

  responseString, err := dispatchFunc(requestBody)

  if err!=nil {
    http.Error(w, fmt.Sprintf("%s", err), 500)
    return
  }

  fmt.Fprint(w, responseString)
}

func getRequestBody(r *http.Request) []byte {
  requestBody := make([]byte, r.ContentLength)
  _, err := r.Body.Read(requestBody)

  if err!=nil && err!=io.EOF {
    log.Print(err)
    return requestBody
  }

  return requestBody
}

func jsHandler(w http.ResponseWriter, r *http.Request) {

  rootDir := http.Dir("js/")
  assetHandler(rootDir, w, r)
}

func cssHandler(w http.ResponseWriter, r *http.Request) {

  w.Header().Set("Content-Type", "text/css")
  rootDir := http.Dir("css/")
  assetHandler(rootDir, w, r)
}

func assetHandler(rootDir http.Dir, w http.ResponseWriter, r *http.Request) {
  filePath := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("/%s", rootDir))

  fileToServe, err := rootDir.Open(filePath)

  if err!=nil {
    log.Print(err)
    http.NotFound(w, r)
    return
  }

  defer fileToServe.Close()

  var fileContent []byte

  for {
    buffer := make([]byte, 2048)
    _, err = fileToServe.Read(buffer)

    if err==io.EOF {
      break
    }
    if err!=nil {
      log.Fatal(err)
    }

    for _, b := range(buffer) {
      if b==0 {
        break
      }
      fileContent = append(fileContent, b)
    }
  }

  fmt.Fprintf(w, "%s", fileContent)
}
