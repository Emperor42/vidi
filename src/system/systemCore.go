package system

import (
	"fmt"
)

type SystemAux struct {
	name string
}

func Init() SystemAux{
	fmt.Println("Hello world")
	return SystemAux{name:"temp"}
}