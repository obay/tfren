package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
)

func main() {
	currentdirectory, err := os.Open(".")
	if err != nil {
		log.Fatal(err)
	}

	defer currentdirectory.Close()

	files, err := currentdirectory.Readdir(-1)
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

				// If the line is a comment, skip it
				if strings.HasPrefix(line, "#") {
					continue
				}

				if strings.HasPrefix(line, "resource") {
					fileNameParts := strings.Split(line, " ")
					if len(fileNameParts) == 4 {
						newFileName := fileNameParts[0] + "." + strings.Replace(fileNameParts[1], "\"", "", -1) + "." + strings.Replace(fileNameParts[2], "\"", "", -1) + ".tf"
						/* Skip if the file is named correctly already ***********************************************************************************************************/
						if newFileName == file.Name() {
							PrintSuccess("Named correctly: " + file.Name() + ". Skipping...")
							break
						}
						/* Make sure file with the same name doesn't already exist ***********************************************************************************************/
						if _, err := os.Stat(newFileName); err == nil {
							// path/to/whatever exists
							PrintError("A file with the name \"" + newFileName + "\" exists already. Unable to rename \"" + file.Name() + "\". Skipping...")
							break
						} else if errors.Is(err, os.ErrNotExist) {
							// path/to/whatever does *not* exist
							// Rename file
							if file.Name() != newFileName {
								e := os.Rename(file.Name(), newFileName)
								if e != nil {
									log.Fatal(e)
								}
								PrintWarning("Renamed file from " + file.Name() + " to " + newFileName)
							}
						} else {
							// File may or may not exist. See err for details.
							// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
							log.Fatal(err)
						}
					}
					break // skip rest of the file once first resource is found
				}
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
			/*********************************************************************************************************************************************************************/
		}
	}
}
