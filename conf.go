package hera

import (
	"bufio"
	"os"
	"strings"
)

var ENV = make(map[string]string)

func NewEnv(filename string) {
	if filename == "" {
		panic("config file is empty")
	}
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		panic("open conf panic")
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		env := strings.Split(line, "=")
		if len(env) != 2 {
			panic("conf env error")
		}
		ENV[env[0]] = env[1]
	}
	return
}
