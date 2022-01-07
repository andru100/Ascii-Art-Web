package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Map containing ascii interpretation of characters
var asciiTable = map[int][]string{}

// Var holds the path to the desired banner file
var banner string

func main() {
	args := os.Args
	// Return expected input if incorrect amount of args given
	if len(args) != 3 {
		fmt.Println("Usage: go run . [STRING] [BANNER]\n\nEX: go run . something standard")
		return
	}

	//	Set the banner file path depending on the argument given
	if err := setBannerFile(args[2]); err != nil {
		fmt.Println("ERROR: " + err.Error())
		return
	}

	// Attempt to fill the asciiTable map
	if err := createAsciiTable(); err != nil {
		fmt.Println("ERROR: " + err.Error())
		return
	}

	// Split input text into lines
	lines := strings.Split(args[1], "\\n")
	// Removing excess nil strings from slice
	if len(lines) > 0 {
		isEmpty := true
		for _, line := range lines {
			if line != "" {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			lines = lines[:len(lines)-1]
		}
	}

	// Print ascii interpretation of input
	for _, line := range lines {
		// If the line is empty then print empty line and continue to next word
		if len([]rune(line)) == 0 {
			fmt.Println()
			continue
		}
		// Iterate over each line of the ascii characters
		for asciiLine := 0; asciiLine < 8; asciiLine++ {
			for asciiChar := 0; asciiChar < len([]rune(line)); asciiChar++ {
				// Print each character in string line by line
				fmt.Print(asciiTable[int(line[asciiChar])][asciiLine])
			}
			fmt.Println()
		}
	}
}

// Read banner file and add contents to asciiTable
func createAsciiTable() error {
	// Attempt to open banner file
	file, err := os.Open(banner)
	if err != nil {
		return err
	}
	// Mark banner file to be closed after function has finished
	defer file.Close()

	// Read contents of file
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	// Add to the asciiTable map using contents of file
	parseBanner(contents)

	return nil
}

// Add ascii character interpretations to asciiTable map
func parseBanner(banner []byte) {
	// Remove the erroneous "\r" character (present in thinkertoy banner)
	correctBanner := strings.ReplaceAll(string(banner), "\r", "")
	// Split banner by lines
	lines := strings.Split(correctBanner, "\n")

	// First character is space which is decimal value of 32
	letter := 32
	for i := 1; i < len(lines); i++ {
		// Once 8 lines of the character has been read, go to next character
		if len(asciiTable[letter]) == 8 {
			letter++
			continue
		}
		// Create an empty slice of strings at map index of letter it doesnt exist
		if _, ok := asciiTable[letter]; !ok {
			asciiTable[letter] = []string{}
		}
		// Add the line to the letter element of the map
		asciiTable[letter] = append(asciiTable[letter], lines[i])
	}
}

// Set banner file from argument given
func setBannerFile(arg string) error {
	switch strings.ToLower(arg) {
	// If one of the valid banner names then change banner to corresponding path
	case "shadow", "standard", "thinkertoy":
		banner = "./ascii-art/banners/" + arg + ".txt"
		return nil
	}
	return errors.New("invalid banner")
}
