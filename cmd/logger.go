package cmd

import (
	"log"
	"os"

	"github.com/fatih/color"
)

var Red = GetColorRed().PrintfFunc()
var Green = GetColorGreen().PrintfFunc()
var Blue = GetColorBlue().PrintfFunc()
var Yellow = GetColorYellow().PrintfFunc()

func GetColorRed() *color.Color {
	col := color.New(color.FgRed)
	return col
}

func GetColorGreen() *color.Color {
	col := color.New(color.FgGreen)
	return col
}

func GetColorBlue() *color.Color {
	col := color.New(color.FgBlue)
	return col
}

func GetColorYellow() *color.Color {
	col := color.New(color.FgYellow)
	return col
}

func Logger() {
	logfile, _ := os.Create("appinfo.log")
	log.SetOutput(logfile)
	log.SetPrefix("Log: ")
	log.SetFlags(log.Ldate | log.Ltime)
}
