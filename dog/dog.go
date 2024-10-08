package dog

import (
	"log"
	"os"
)

type Dog struct {
	Active bool
	file   os.File
	loggy  log.Logger
	num    int
}

func Init(d *Dog) (ok bool) {
	file, err := os.Create("dogfile")
	if err != nil {
		return false
	}
	d.file = *file
	d.loggy.SetOutput(&d.file)
	log.SetOutput(&d.file)
	return true
}

func (d *Dog) Log(in ...any) {
	if d.Active {
		d.loggy.Print("[", d.num, "]")
		d.loggy.Println(in)
		d.num += 1
	}

}
