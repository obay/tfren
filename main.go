package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".tf") && !file.IsDir() {
			/*********************************************************************************************************************************************************************/
			f, err := os.Open(file.Name())
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "resource") {
					fileNameParts := strings.Split(line, " ")
					if len(fileNameParts) == 4 {
						newFileName := fileNameParts[0] + "." + strings.Replace(fileNameParts[1], "\"", "", -1) + "." + strings.Replace(fileNameParts[2], "\"", "", -1) + ".tf"
						if file.Name() != newFileName {
							e := os.Rename(file.Name(), newFileName)
							if e != nil {
								log.Fatal(e)
							}
							fmt.Println("Renamed file from " + file.Name() + " to " + newFileName)
						}
					}
				}
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
			/*********************************************************************************************************************************************************************/
		}
	}
}
