package main

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
)

func init() {
	apiKey := os.Getenv("UNIPDF_API_KEY")
	err := license.SetMeteredKey(apiKey)
	if err != nil {
		fmt.Printf("ERROR: Failed to set metered key: %v\n", err)
		fmt.Printf("Make sure to get a valid key from https://cloud.unidoc.io\n")
		panic(err)
	}
}

func main() {
	lk := license.GetLicenseKey()
	if lk == nil {
		fmt.Printf("Failed retrieving license key")
		return
	}
	fmt.Printf("License: %s\n", lk.ToString())

	// GetMeteredState freshly checks the state, contacting the licensing server.
	state, err := license.GetMeteredState()
	if err != nil {
		fmt.Printf("ERROR getting metered state: %+v\n", err)
		panic(err)
	}
	fmt.Printf("State: %+v\n", state)
	if state.OK {
		fmt.Printf("State is OK\n")
	} else {
		fmt.Printf("State is not OK\n")
	}
	fmt.Printf("Credits: %v\n", state.Credits)
	fmt.Printf("Used credits: %v\n", state.Used)
}
