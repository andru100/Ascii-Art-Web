package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
)

// This struct holds relavent info needed for the ascii-art function
type params struct {
	Text   string
	Banner string
}

// This struct holds the ascii-art return to be returned from the POST request
type returnValue struct {
	Value string
}

// Map of allowed banners
var bannerMap = map[string]bool{"standard": true, "shadow": true, "thinkertoy": true}

var mostRecentAscii string

var result []byte

func main() {
	StartServer()
}

// Run ascii-web-server
func StartServer() {
	// Serve the index.html file to "/"
	http.Handle("/", http.FileServer(http.Dir("./")))

	// Handle the "/ascii-art" request
	http.HandleFunc("/ascii-art", getAscii)
	
	// Handle download request for txt file
	http.HandleFunc("/txt", txt)

	// Handle download request for html file
	http.HandleFunc("/html", html)

	// Print server start message and attempt to listen at port 8080
	fmt.Println("Starting server at web port 4000")
	if err := http.ListenAndServe(":4000", nil); err != nil {
		panic(err)
	}
}

// Receive params, run ascii-art program with params as args and respond with output
func getAscii(w http.ResponseWriter, r *http.Request) {
	// Set CORS policies
	setCORS(w)

	// Create variable used to check if http request is Bad
	requestOk := true

	// Read the contents of the body of the request
	contents, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// Request cannot be read so respond with a http error code 400, bad request
		http.Error(w, "Bad request", 400)
		requestOk = false
	}
	args := params{}
	// Attempt to read the contents into a params struct
	err = json.Unmarshal(contents, &args)
	if err != nil {
		panic(err)
	}

	// Check if the request is ok and that the banner requested is available
	if requestOk && bannerMap[args.Banner] {

		/*	Iterate over chars in the Text adding the char to formattedText, if a newline char
			is found add "\n" instead*/
		formattedText := ""
		for _, char := range args.Text {
			if char == 10 {
				formattedText += "\\n"
			} else {
				formattedText += string(char)
			}
		}

		// Run the ascii-art program with params passed as ards
		result, err = exec.Command(
			"go", "run", "ascii-art/ascii-art.go", formattedText, args.Banner,
		).Output()
		if err != nil {
			// Error running the program so notify http client with error code 500, Internal server error
			http.Error(w, "Internal Server Error", 500)
			panic(err)
		}

		// Format the result of the ascii-art program to a json struct to be returned to caller of request
		response, err := json.Marshal(returnValue{Value: string(result)})
		if err != nil {
			// Notify client of error
			http.Error(w, "Internal Server Error", 500)
			panic(err)
		}

		// Store the returned ascii in a var
		mostRecentAscii = string(result)

		// Return the response
		w.Write(response)

	} else {
		// In the case where a banner hasnt been found. Respond with an error code 404.
		http.Error(w, "Not Found", 404)
	}
}

// Return the ascii-art text as a file to be exported
func txt(w http.ResponseWriter, r *http.Request) {
	// Set CORS policies
	setCORS(w)

	// create response and send
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len([]byte(mostRecentAscii))))
	w.Header().Set("Content-Disposition", "attachment; filename=ascii-art.txt")
	w.Write([]byte(mostRecentAscii))
}

// Return a html version of ascii art
func html(w http.ResponseWriter, r *http.Request) {
	// Set CORS policies
	setCORS(w)

	// Create the html
	html := "<!DOCTYPE html> <html lang='en'> <head> <meta charset='UTF-8'> <meta http-equiv='X-UA-Compatible' content='IE=edge'> <meta name='viewport' content='width=device-width, initial-scale=1.0'> <title>Document</title> </head> <body> <div> <pre>"
	html += string(mostRecentAscii)
	html += "</pre> <div  style='clear:both'></div> </div> </body> </html>"
	
	// create response and send
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Length", strconv.Itoa(len(html)))
	w.Header().Set("Content-Disposition", "attachment; filename=ascii.html")
	w.Write([]byte(html))
}


// This func sets the CORS rules
func setCORS(w http.ResponseWriter) {
	header := w.Header()
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Access-Control-Allow-Methods", "POST, GET")
	header.Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
}

func check(err error) { // handle errors
    if err != nil {
        panic(err)
    }
}