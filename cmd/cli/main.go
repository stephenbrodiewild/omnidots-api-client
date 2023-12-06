package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/stephenbrodiewild/omnidots-api-client/pkg/client"
)

func main() {
	// Define and parse command-line flags
	tokenFlag := flag.String("token", "", "Omnidots API token")
	commandFlag := flag.String("command", "", "Command to execute ('list-sensors')")
	flag.Parse()

	// Check for required token
	if *tokenFlag == "" {
		fmt.Println("Error: API token is required")
		os.Exit(1)
	}

	// Create client
	ctx := context.Background()
	server := "http://honeycomb.omnidots.com/api/v1"
	omnidotsClient, err := client.NewClientWithResponses(server, *tokenFlag)
	if err != nil {
		fmt.Printf("Failed to initialize client: %v\n", err)
		os.Exit(1)
	}

	// Execute command
	switch *commandFlag {
	case "list-sensors":
		listSensors(ctx, omnidotsClient)
	default:
		fmt.Printf("Unknown command: %s\n", *commandFlag)
		os.Exit(1)
	}
}

// listSensors lists the sensors
func listSensors(ctx context.Context, omnidotsClient *client.ClientWithResponses) {
	sensorsResponse, err := omnidotsClient.ListSensorsWithResponse(ctx)
	if err != nil {
		fmt.Printf("Failed to list sensors: %v\n", err)
		os.Exit(1)
	}

	if sensorsResponse.JSON200 != nil {
		bytes, err := json.Marshal(*sensorsResponse.JSON200.Sensors)
		if err != nil {
			fmt.Println("Failed to marshal the sensors")
			os.Exit(1)
		}
		fmt.Println(string(bytes))
	} else {
		fmt.Println("No sensors found or an error occurred in fetching the data.")
	}
}
