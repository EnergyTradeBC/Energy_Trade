/*
Copyright 2021 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type NetworkConfig struct {
	ChannelName          string
	ChaincodeMoneyName   string
	ChaincodeEnergyName  string
	ChaincodeAuctionName string
}

type moneySettingsStruct struct {
	startingBid     float32
	maxPrizeAuction float32
}

const orgMSP = "Org1MSP"

const assetID = "energy"
const auctionID = "auction_1"

const moneyName = ""
const moneySymbol = ""
const moneyDecimals = ""

// var bidTransactionID = ""
var moneySettings = &moneySettingsStruct{
	startingBid:     0.0,
	maxPrizeAuction: 0.0,
}

var contract_money *client.Contract
var contract_energy *client.Contract
var contract_auction *client.Contract

var c_excess = ""

func main() {
	// SETUP THE MQTT CONNECTION TO THE SMART METER
	broker, port, topic, clientID, username, password := readMQTTConfig()

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	clientMQTT := mqtt.NewClient(opts)
	if token := clientMQTT.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Subscribe to the smart meter's topic
	sub_topic(clientMQTT, topic)

	// SETUP THE REST INTERFACE
	router := gin.Default()
	router.GET("/moneySettings", getMoneySettings)
	router.POST("/moneySettings", postMoneySettings)

	router.GET("/currentBalance", getCurrentBalance)
	router.GET("/moneyTransactions", getMoneyTransactions)

	router.GET("/currentEnergyAsset", getCurrentEnergyAsset)
	router.GET("/energyTransactions", getEnergyTransactions)
	router.GET("/remainingAssets", getRemainingEnergy)

	router.Run("localhost:8080")

	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	gateway, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gateway.Close()

	// Read network information from a json file
	channelName, chaincodeMoneyName, chaincodeEnergyName, chaincodeAuctionName := readNetworkConfig()

	// Get network and smart contract objects
	network := gateway.GetNetwork(channelName)
	// Update the global variable defined in the package scope
	contract_money = network.GetContract(chaincodeMoneyName)
	contract_energy = network.GetContract(chaincodeEnergyName)
	contract_auction = network.GetContract(chaincodeAuctionName)

	// Context used for event listening
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen for events emitted by subsequent transactions
	startEnergyChaincodeEventListening(ctx, network, chaincodeEnergyName)
	startAuctionChaincodeEventListening(ctx, network, chaincodeAuctionName)

	for {

	}

}

// UTILITY FUNCTIONS

// Reads from a json the channel name and the name of the chaincodes committed to that channel
func readNetworkConfig() (string, string, string, string) {
	content, err := ioutil.ReadFile("./network_config.json")
	if err != nil {
		panic(fmt.Errorf("error when opening file: %w", err))
	}

	var payload NetworkConfig
	err = json.Unmarshal(content, &payload)
	if err != nil {
		panic(fmt.Errorf("error during Unmarshal(): %w", err))
	}

	return payload.ChannelName, payload.ChaincodeMoneyName, payload.ChaincodeEnergyName, payload.ChaincodeAuctionName
}

// Reads from a json the broker, port and topic to perform MQTT communications
func readMQTTConfig() (string, int, string, string, string, string) {
	content, err := ioutil.ReadFile("./mqtt_config.json")
	if err != nil {
		panic(fmt.Errorf("error when opening file: %w", err))
	}

	var payload MQTTConfig
	err = json.Unmarshal(content, &payload)
	if err != nil {
		panic(fmt.Errorf("error during Unmarshal(): %w", err))
	}

	return payload.broker, payload.port, payload.topic, payload.clientID, payload.username, payload.password
}

// Format JSON data
func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}
