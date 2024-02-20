package logging

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var stdout = os.Stdout

type printer struct {
	print  func(prefix string, a ...any)
	printf func(prefix, format string, a ...any)
}

var (
	instance *printer
	lock     sync.Mutex
)

func init() {
	log.SetOutput(stdout)
	log.SetFlags(0)
	DisableGitHubFormat()
}

func EnableGitHubFormat() {
	lock.Lock()
	defer lock.Unlock()
	instance.print = func(prefix string, a ...any) {
		log.Println(append([]any{prefix}, a...)...)
	}
	instance.printf = func(prefix, format string, a ...any) {
		log.Printf(fmt.Sprintf("%s%s", prefix, format), a...)
	}
}

func DisableGitHubFormat() {
	lock.Lock()
	defer lock.Unlock()
	instance = defaultPrinter()
}

func defaultPrinter() *printer {
	return &printer{
		print: func(prefix string, a ...any) {
			log.Println(a...)
		},
		printf: func(prefix, format string, a ...any) {
			log.Printf(format, a...)
		},
	}
}
