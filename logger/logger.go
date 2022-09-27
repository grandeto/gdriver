package logger

import "log"

func Info(args ...interface{}) {
	log.Printf("INFO: %v", args)
}

func Error(args ...interface{}) {
	log.Printf("ERROR: %v", args)
}

func Fatal(args ...interface{}) {
	log.Fatalf("FATAL: %v", args)
}
