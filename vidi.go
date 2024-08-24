package vidi

import "fmt"

type VidiContext struct {
	systemTest string
}

func (v *VidiContext) ProcessBody() {
	fmt.Println("temp")
}
