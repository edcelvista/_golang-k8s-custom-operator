package lib

import (
	"log"
	"os"
	"strings"
)

var debug = false

func Debug(v ...interface{}) {
	if debug {
		log.Println(v...)
	}
}

func Init() {
	if strings.ToLower(os.Getenv("IS_DEBUG")) == "true" {
		debug = true
	}

	log.Printf("💡 Debug enabled: %v", debug)
	// •	%v → Print the values
	// •	%+v → Print field names and values
	// •	%#v → Print Go syntax (main.Person{Name:"Alice", Age:30})
}
