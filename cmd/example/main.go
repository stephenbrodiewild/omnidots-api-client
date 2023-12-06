package main

import (
	"context"
	"fmt"
	"os"

	"github.com/stephenbrodiewild/omnidots-api-client/pkg/client"
)

func main() {
	ctx := context.Background()
	server := "http://honeycomb.omnidots.com/api/v1"
	token := os.Getenv("OMNIDOTS_TOKEN")

	client, err := client.NewClientWithResponses(server, token)
	if err != nil {
		panic(fmt.Errorf("failed to initialise client: %s", err))

	}

	sensorsResponse, err := client.ListSensorsWithResponse(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to list sensors: %s", err))
	}

	if sensorsResponse.JSON200 != nil {
		for _, sensor := range *sensorsResponse.JSON200.Sensors {
			fmt.Println(*sensor.Name)
		}
	}
}
