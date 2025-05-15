package logger

import "log"

type Logger interface {
	Println(args ...any)
	Printf(format string, args ...any)
}

type Silent struct{}

func (l Silent) Println(args ...any)               {}
func (l Silent) Printf(format string, args ...any) {}

type Verbose struct{}

func (v Verbose) Println(args ...any) {
	log.Println(args...)
}
func (v Verbose) Printf(format string, args ...any) {
	log.Printf(format, args...)
}
