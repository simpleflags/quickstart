package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	sdk "github.com/simpleflags/golang-server-sdk"
)

const sdkKey = "PUT your SDK key"

const featureFlagKey = "bool-flag"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	target := map[string]interface{}{
		"identifier": "enver",
	}

	err := sdk.Initialize(sdkKey)
	if err != nil {
		log.Printf("could not connect to SF servers %v", err)
	}

	defer func() {
		if err := sdk.Close(); err != nil {
			log.Printf("error while closing client err: %v", err)
		}
	}()
	sdk.WaitForInitialization()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				showFeature := sdk.Evaluate(featureFlagKey, target).Bool(false)

				fmt.Printf("KeyFeature flag '%s' is %t for this user\n", featureFlagKey, showFeature)
				time.Sleep(10 * time.Second)
			}
		}
	}()

	<-ctx.Done()
}
