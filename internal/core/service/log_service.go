package service

import (
	"log"
	"os"
	"regexp"
	"strings"
)

func ViewLogs() {
	data, err := os.ReadFile(logFile)
	if err != nil {
		log.Println("Error reading logs:", err)
		return
	}
	log.Println("LOG FILE CONTENT:\n", string(data))
}

func FilterLogs(filter string) {
	data, err := os.ReadFile(logFile)
	if err != nil {
		log.Println("Error reading logs:", err)
		return
	}

	lines := strings.Split(string(data), "\n")
	regex, err := regexp.Compile(filter)
	if err != nil {
		log.Println("Invalid regex:", err)
		return
	}

	for _, line := range lines {
		if regex.MatchString(line) {
			log.Println(line)
		}
	}
}
