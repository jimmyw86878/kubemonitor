package logger

import (
	"log"
	"os"
)

var (
	//Info define
	Info *log.Logger
	//Warning define
	Warning *log.Logger
	//Error define
	Error *log.Logger
)

//NewLogger is to generate logger
func NewLogger() {
	Info = log.New(os.Stdout,
		"INFO: ",
		log.LstdFlags|log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(os.Stdout,
		"WARNING: ",
		log.LstdFlags|log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(os.Stdout,
		"ERROR: ",
		log.LstdFlags|log.Ldate|log.Ltime|log.Lshortfile)
}
