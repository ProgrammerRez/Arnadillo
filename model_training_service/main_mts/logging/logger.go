package logging

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

// Creating Logging Levels
const(
	InfoLevel = iota
	WarnLevel
	ErrorLevel
)

// Creating the Logger Struct
type Logger struct{
	infoLogger		*log.Logger
	warnLogger		*log.Logger
	errorLogger		*log.Logger
	Level			int
}

var logger *Logger

const LOG_PATH string = "./logs.txt"


// Init Function

func Init(){

	// Create the Log file

	f, err := CreateLoggingFile(LOG_PATH)

	if err != nil{
		log.Fatal(err.Error())
	}

	log.SetOutput(f)
	flags := log.Ldate | log.Lmicroseconds | log.Lshortfile

	logger = &Logger{
		infoLogger: log.New(f, "INFO: ",flags),
		warnLogger: log.New(f, "WARN: ",flags),
		errorLogger: log.New(f, "ERROR: ",flags),
	}
}


// This function initializes the logging file on the `LOG_PATH`
func CreateLoggingFile(path string) (*os.File, error){

	dir := filepath.Dir(path)

	if dir != "."{
		if err := os.MkdirAll(path, os.ModePerm); err != nil{
			return nil, errors.New("Unable to create log path")
		}
	}

	f, err := os.OpenFile(path, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
		if err != nil{
			return nil, errors.New("Unable to create log file")
		}
	
	return f, nil

}

// Logging Functions

// Set the Level for the Logger
func SetLevel(level int){
	if logger != nil{
		logger.Level = level
	}
}

// Info Logging
func Info(message string){
	if logger.Level <= InfoLevel{
		logger.infoLogger.Println(message)
	}
}

// Warn Logging
func Warn(message string){
	if logger.Level <= WarnLevel{
		logger.warnLogger.Println(message)
	}
}

// Error Logging
func Error(message string){
	if logger.Level <= ErrorLevel{
		logger.errorLogger.Println(message)
	}
}