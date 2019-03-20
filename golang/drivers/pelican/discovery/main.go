package main

import (
	"fmt"
	"os"

	"github.com/SoftwareDefinedBuildings/bw2-contrib/driver/pelican/storage"
	"github.com/SoftwareDefinedBuildings/bw2-contrib/driver/pelican/types"
	"github.com/immesys/spawnpoint/spawnable"

	_ "github.com/lib/pq"
)

func main() {
	params := spawnable.GetParamsOrExit()
	username := params.MustString("username")
	password := params.MustString("password")
	sitename := params.MustString("sitename")

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
