package main

import (
	"fmt"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/svg"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	iconFileName := "icons.js"

	_, errIconFileName := os.Stat(iconFileName)
	if errIconFileName != nil {
		if os.IsExist(errIconFileName) {
			log.Fatalf(errIconFileName.Error())
		}
	}

	f, errCreateFile := os.Create(iconFileName)
	if errCreateFile != nil {
		log.Fatal(errCreateFile)
		return
	}

	_, err := f.WriteString("const Icons = {\n")
	if err != nil {
		return
	}

	defer func(f *os.File, s string) {
		_, err := f.WriteString(s)
		if err != nil {
			return
		}
	}(f, "\n}\n\nexport default Icons")

	var pathName string
	pathName = "icons/"

	_, errStat := os.Stat(pathName)

	if errStat != nil {
		log.Fatalf("%v", errStat)
	}

	readDir, errReadDir := os.ReadDir(pathName)

	if errReadDir != nil {
		log.Fatalf("%v", errReadDir)
	}

	countFor := 0

	for _, element := range readDir {
		countFor++
		isValid, _ := regexp.MatchString("\\.svg", element.Name())
		if !isValid {
			fmt.Printf("Bu dosya SVG değildir: %v \n", element.Name())
			continue
		}

		content, err := ioutil.ReadFile(fmt.Sprintf("%v%v", pathName, element.Name()))

		if err != nil {
			log.Fatal(err)
		}

		iconName := element.Name()[0:(len(element.Name()) - 4)] // Icon name
		iconContent := string(content)                          // Icon Content

		m := minify.New()
		m.AddFunc("image/svg+xml", svg.Minify)

		var s string

		s, err = m.String("image/svg+xml", iconContent)
		if err != nil {
			panic(err)
		}

		var hasComma string

		if len(readDir) == countFor {
			hasComma = ""
		}else{
			hasComma = ",\n"
		}

		_, err = f.WriteString(fmt.Sprintf("\t\"%v\": \"%v\"%v", iconName, strings.ReplaceAll(s, "\"", "\\\""), hasComma))
		if err != nil {
			return
		}
	}
}
