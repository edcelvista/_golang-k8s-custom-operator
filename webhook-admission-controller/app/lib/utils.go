package lib

import (
	"log"
	"os"
	"strings"
)

var debug = false

func Debug(v string) {
	if debug {
		log.Printf("🚨 [DEBUG] %v", v)
	}
}

func DefaultIfEmpty(s, def string) string {
	if s == "" {
		return def
	}
	return s
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
