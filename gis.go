package gis

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var root = os.Getenv("GOROOT")
var pattern = regexp.MustCompile("^type [A-Z][A-Za-z]* interface ?{")

func List() {
	fmt.Println("looking in", root)

	i := 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if pattern.Match(scanner.Bytes()) {
				i += 1
				fmt.Println(path, scanner.Text())
			}
		}

		err = scanner.Err()
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(i)
}
