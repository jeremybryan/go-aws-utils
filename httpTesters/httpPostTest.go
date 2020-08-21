package main

import (
	"flag"
	"fmt"
	"net/http"
)

/**
Simple Http post example
 */
func main ()  {
	endpoint := flag.String("endpoint", "", "a string")
	flag.Parse()

	if *endpoint != "" {
		resp, err := http.Post(*endpoint, "application/json", nil)
		if err != nil {
		   fmt.Errorf("Error completing HTTP POST %s\n", err.Error())
		} else {
                   defer resp.Body.Close()	
		   fmt.Printf("HTTP Post to %s has been completed. Response status: %s", *endpoint, resp.Status)
                }
	} else {
		fmt.Print("No endpoint has been defined. Exiting")
	}
}
