package main

import (
    "log"
    "net/http"
    "vidi/database/vidi"
    "vidi/system"
)

func main() {

    http.Handle("/", http.FileServer(http.Dir("./server")))

    log.Fatal(http.ListenAndServe(":8082", nil))

    system.TempFlag()

    vidi.TempFlag()

}
