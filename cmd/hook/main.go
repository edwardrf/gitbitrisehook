package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	bitrise "github.com/edwardrf/gitbitrisehook"
)

const config = "bitrise.json"

func main() {
	cmd := filepath.Base(os.Args[0])

	conf, err := os.Open("bitrise.json")
	if err != nil {
		log.Printf("Failed to open config file %v due to error: %v", config, err)
	}

	params := make(map[string]string)
	err = json.NewDecoder(conf).Decode(&params)
	if err != nil {
		log.Printf("Failed to decode bitrise config file with error: %v", err)
	}

	br := bitrise.New(params["api_slug"], params["api_token"])

	switch cmd {
	case "post-update":
		if len(os.Args) < 2 {
			log.Printf("Not enough arguments for post-update")
		}
		err := br.Trigger(os.Args[1], "")
		if err != nil {
			log.Printf("Failed to process %v trigger with error: %v", cmd, err)
		}
	case "post-receive":
		lines := bufio.NewScanner(os.Stdin)
		for lines.Scan() {
			line := lines.Text()
			parts := strings.Split(line, " ")
			if len(parts) < 3 {
				log.Printf("Invalid input for post-receive hook: [%s]", line)
			}
			newHash := parts[1]
			ref := parts[2]
			err := br.Trigger(ref, newHash)
			if err != nil {
				log.Printf("Failed to send request: %v", err)
			}
		}
	default:
		log.Printf("hook %v not implemented yet", cmd)
	}

}
