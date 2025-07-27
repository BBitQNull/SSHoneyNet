package pathconfig

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

type PathFlag string

const (
	JSON PathFlag = "json"
	LOG  PathFlag = "log"
)

var (
	jsonPath = flag.String("jsonPath", "", "Path to log JSON file")
	logPath  = flag.String("logPath", "", "Path to save SSHoneyNet_log.jsonl")
)

func GetPath(config PathFlag) string {
	flag.Parse()
	switch config {
	case JSON:
		if *jsonPath != "" {
			absPath, err := filepath.Abs(*jsonPath)
			if err != nil {
				log.Fatalf("Invalid custom jsonPath: %v", err)
			}
			return absPath
		}
		dir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get working dir: %v", err)
		}
		return filepath.Join(dir, "pkg", "model", "fs", "file_tree.json")
	case LOG:
		if *logPath != "" {
			absPath, err := filepath.Abs(*logPath)
			if err != nil {
				log.Fatalf("Invalid custom logPath: %v", err)
			}
			return absPath
		}
		dir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get working dir: %v", err)
		}
		return filepath.Join(dir, "pkg", "model", "log", "SSHoneyNet_log.jsonl")
	}
	return ""
}
