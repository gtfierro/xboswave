package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/john-b-yang/xboswave/golang/drivers/pelican/storage"
	"github.com/john-b-yang/xboswave/golang/drivers/pelican/types"

	_ "github.com/lib/pq"
)

func main() {
	var configBytes []byte
	configBytes, readErr := ioutil.ReadFile("../config.json")
	if readErr != nil {
		fmt.Printf("Failed to read config.json file properly: %s\n", readErr)
	}
	var configData map[string]interface{}
	if unmarshalErr := json.Unmarshal(configBytes, &configData); unmarshalErr != nil {
		fmt.Printf("Failed to unmarshal config.json file properly: %s\n", configData)
	}
	username := configData["username"].(string)
	password := configData["password"].(string)
	sitename := configData["sitename"].(string)

	pelicans, err := types.DiscoverPelicans(username, password, sitename)
	if err != nil {
		fmt.Printf("Failed to discover Pelican thermostats: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Discovered %d Pelican(s), writing to remote DB...\n", len(pelicans))
	if err = storage.WritePelicans(pelicans, sitename); err != nil {
		fmt.Printf("Failed to write to database: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Success!")
}
