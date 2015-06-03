package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type device struct {
	Data []data
}

type data struct {
	ID       string
	Type     string
	Name     string
	Promoted bool
}

func main() {
	methodsSupported := []string{"type", "name", "promoted"}
	var result device
	var deviceData []data
	var boolValue bool

	// check for valid command line arguments, print usage details
	if len(os.Args) < 3 {
		fmt.Println("usage: sphere-thing-deleter [method] [value]")
		fmt.Println("Supported methods:", strings.Join(methodsSupported, ", "))
		fmt.Println("Examples")
		fmt.Println("To delete all non-promoted things, use:                 ... promoted false")
		fmt.Println("To delete all things with type 'light', use:            ... type light")
		fmt.Println("To delete all things with names containing 'jim', use:  ... name jim")
		os.Exit(1)
	}
	// check second command-line argument, method
	method := os.Args[1]
	value := os.Args[2]
	if !isStringInSlice(method, methodsSupported) {
		fmt.Println("Invalid method. Supported methods:", strings.Join(methodsSupported, ", "))
		os.Exit(1)
	}

	// convert true/false string if needed
	if method == "promoted" {
		var err error
		boolValue, err = strconv.ParseBool(value)
		if err != nil {
			fmt.Println("Invalid value. Must be true or false")
			os.Exit(1)
		}
	}
	// handle spaces in command line arguments (for names)
	if len(os.Args) > 3 {
		value = strings.Join(os.Args[2:], " ")
	}

	b := getThingsJSON()
	// convert bytes (string) to device data type with fields
	err := json.Unmarshal(b, &result)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	deviceData = result.Data
	// loop through all devices, delete those that match parameters
	for i := 0; i < len(deviceData); i++ {
		switch method {
		case "type":
			if deviceData[i].Type == value {
				fmt.Printf("Deleting %s %s with ID: %s\n", deviceData[i].Type, deviceData[i].Name, deviceData[i].ID)
				deleteThing(deviceData[i].ID)
			}
		case "name":
			if strings.Contains(deviceData[i].Name, value) {
				fmt.Printf("Deleting %s %s with ID: %s\n", deviceData[i].Type, deviceData[i].Name, deviceData[i].ID)
				deleteThing(deviceData[i].ID)
			}
		case "promoted":
			if deviceData[i].Promoted == boolValue {
				fmt.Printf("Deleting %s %s with ID: %s\n", deviceData[i].Type, deviceData[i].Name, deviceData[i].ID)
				deleteThing(deviceData[i].ID)
			}
		default:
			fmt.Println("Invalid method. Use 'type' or 'name'")
		}
	}
}

func getThingsJSON() []byte {
	resp, err := http.Get("http://ninjasphere.local:8000/rest/v1/things")
	if err != nil {
		fmt.Printf("Error. %v\n", err)
		return []byte{}
	}
	defer resp.Body.Close()
	dataGet, err := ioutil.ReadAll(resp.Body)

	return dataGet

	// curl way:
	//	var (
	//		cmdOut []byte
	//		err    error
	//	)
	//	cmdName := "curl"
	//	cmdArgs := []string{"-s", "http://ninjasphere.local:8000/rest/v1/things"}
	//
	//	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
	//		fmt.Fprintln(os.Stderr, "There was an error running command: ", err)
	//		os.Exit(1)
	//	}
	//
	//	return cmdOut
}

func deleteThing(id string) {
	cmdName := "curl"
	cmdArgs := strings.Fields("-X DELETE http://ninjasphere.local:8000/rest/v1/things/" + id)
	if _, err := exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running command: ", err)
		os.Exit(1)
	}
}

func isStringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
