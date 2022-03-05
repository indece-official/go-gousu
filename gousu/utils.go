package gousu

import (
	"log"
	"os"

	"github.com/indece-official/go-gousu/v2/gousu/logger"
)

var logFatalf = log.Fatalf

// ContainsString checks if a slice of strings contains a string
func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// CheckError checks for an error and exits the process with result code 1 if err is set
func CheckError(err error) {
	if err != nil {
		log := logger.GetLogger("utils.error")
		log.Fatalf("Error: %s", err)
		os.Exit(1)
	}
}
