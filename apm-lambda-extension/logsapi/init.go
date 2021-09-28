// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0

package logsapi

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

const DefaultHttpListenerPort = "1234"

// Init initializes the configuration for the Logs API and subscribes to the Logs API for HTTP
func Subscribe(extensionID string, eventTypes []EventType) error {
	extensions_api_address, ok := os.LookupEnv("AWS_LAMBDA_RUNTIME_API")
	if !ok {
		return errors.New("AWS_LAMBDA_RUNTIME_API is not set")
	}

	logsAPIBaseUrl := fmt.Sprintf("http://%s", extensions_api_address)

	logsAPIClient, err := NewClient(logsAPIBaseUrl)
	if err != nil {
		return err
	}

	bufferingCfg := BufferingCfg{
		MaxItems:  10000,
		MaxBytes:  262144,
		TimeoutMS: 1000,
	}
	if err != nil {
		return err
	}
	address := ListenOnAddress()
	if !strings.HasPrefix(address, "http://") {
		address = "http://" + address
	}
	destination := Destination{
		Protocol:   HttpProto,
		URI:        URI(address),
		HttpMethod: HttpPost,
		Encoding:   JSON,
	}

	_, err = logsAPIClient.Subscribe(eventTypes, bufferingCfg, destination, extensionID)
	return err
}
