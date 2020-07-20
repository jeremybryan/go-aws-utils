package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main ()  {
	endpoint := flag.String("endpoint", "", "a string")
	flag.Parse()

	if *endpoint != "" {
		resp, err := http.Post(*endpoint, "application/json", nil)
		if err != nil {
			fmt.Errorf("Error completing HTTP POST %s", err.Error())
		}
		defer resp.Body.Close()
		fmt.Printf("HTTP Post to %s has been completedl Response status: %s", endpoint, resp.Status)
	} else {
		fmt.Print("No endpoint has been defined. Exiting")
	}
}
