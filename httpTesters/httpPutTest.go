package main

import (
	"flag"
	"fmt"
	"net/http"
)

/**
Simple HTTP put example
 */
func main ()  {
	endpoint := flag.String("endpoint", "", "a string")
	flag.Parse()

	if *endpoint != "" {
		req, err := http.NewRequest("PUT", *endpoint, nil)
		if err != nil {
			fmt.Errorf("Error creating HTTP PUT request %s", err.Error())
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Errorf("Error completing HTTP PUT %s", err.Error())
			return
		}
		defer resp.Body.Close()
		fmt.Printf("HTTP Post to %s has been completed.\n Response status: %s\n", *endpoint, resp.Status)
	} else {
		fmt.Print("No endpoint has been defined. Exiting")
	}
}
