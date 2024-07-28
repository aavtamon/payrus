// PizTec Corporation, 2024. All Rights Reserved.

package main

import "fmt"

var Logger *logger = &logger{}

type logger struct {
}

func (log *logger) Error(format string, params ...any) {
	log.logMessage("ERROR", format, params...)
}

func (log *logger) Info(format string, params ...any) {
	log.logMessage("INFO", format, params...)
}

func (log *logger) Warning(format string, params ...any) {
	log.logMessage("WARN", format, params...)
}

func (log *logger) Debug(format string, params ...any) {
	log.logMessage("DEBUG", format, params...)
}

func (log *logger) logMessage(level string, format string, params ...any) {
	fmt.Printf(level+": "+format+"\n", params...)
}
