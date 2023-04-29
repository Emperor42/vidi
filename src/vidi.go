package main

import (
    "log"
    "net/http"
    "vidi/database/vidi"
    "vidi/system"
)

func main() {

    system := system.Init()

    core:= vidi.Init()

    log.Println(system)

    log.Println(core)

    http.Handle("/", http.FileServer(http.Dir("./server")))

    log.Fatal(http.ListenAndServe(":8082", nil))

}
