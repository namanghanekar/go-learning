package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var Logs []string
var mu sync.Mutex
var logPathOnce sync.Once
var logPath string

func Log(event string) {
	msg := fmt.Sprintf("[%s] %s",
		time.Now().Format("15:04:05"),
		event,
	)

	log.Println(msg)

	mu.Lock()
	Logs = append(Logs, msg)
	appendToFile(msg)
	mu.Unlock()
}

func GetLogs() []string {
	mu.Lock()
	defer mu.Unlock()

	if fileLogs, err := os.ReadFile(getLogPath()); err == nil && len(fileLogs) > 0 {
		lines := splitLines(string(fileLogs))
		if len(lines) > 0 {
			return lines
		}
	}

	cp := make([]string, len(Logs))
	copy(cp, Logs)
	return cp
}

func appendToFile(msg string) {
	path := getLogPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		log.Println("log directory error:", err)
		return
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Println("log file error:", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(msg + "\n"); err != nil {
		log.Println("log write error:", err)
	}
}

func getLogPath() string {
	logPathOnce.Do(func() {
		logPath = filepath.Join(findProjectRoot(), "doc_logs", "ticket-system.log")
	})
	return logPath
}

func findProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "shared", "logger", "logger.go")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return dir
		}
		dir = parent
	}
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, r := range s {
		if r == '\n' {
			line := s[start:i]
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			if line != "" {
				lines = append(lines, line)
			}
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
