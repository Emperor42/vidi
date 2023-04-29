package vidi

import (
	"fmt"
)

func TempFlag(){
	fmt.Println("Hello world")
}

type Database interface {
	Connect() Database
	Access(GID) string
}

type GID interface {
	GenerateString() string
	Access(uint64) uint64
	Check(GID) bool
}