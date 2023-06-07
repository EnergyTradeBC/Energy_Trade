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
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

var assetID = "energy_1"
var orgMSP = "Org1MSP"
var auctionID = "auction_1"
var bidTransactionID = ""
var contract_money *client.Contract
var contract_energy *client.Contract
var contract_auction *client.Contract

type NetworkConfig struct {
	ChannelName          string
	ChaincodeMoneyName   string
	ChaincodeEnergyName  string
	ChaincodeAuctionName string
}

type MQTTConfig struct {
	broker   string
	port     int
	topic    string
	clientID string
	username string
	password string
}

type MQTTmessage struct {
	ProducedEnergy float64
	ConsumedEnergy float64
}

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
	router.GET("/primo_get", getPRIMO)
	router.GET("/secondo_get", getSECONDO)

	router.POST("/primo_post", postPRIMO)
	router.POST("/secondo_post", postSECONDO)

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
	startMoneyChaincodeEventListening(ctx, network, chaincodeMoneyName)
	startEnergyChaincodeEventListening(ctx, network, chaincodeEnergyName)
	startAuctionChaincodeEventListening(ctx, network, chaincodeAuctionName)

	for true {

	}

}

// UTILITY FUNCTIONS

// Reads from a json the channel name and the name of the chaincodes committed to that channel
func readNetworkConfig() (string, string, string, string) {
	content, err := ioutil.ReadFile("./network_config.json")
	if err != nil {
		panic(fmt.Errorf("Error when opening file: %w", err))
	}

	var payload NetworkConfig
	err = json.Unmarshal(content, &payload)
	if err != nil {
		panic(fmt.Errorf("Error during Unmarshal(): %w", err))
	}

	return payload.ChannelName, payload.ChaincodeMoneyName, payload.ChaincodeEnergyName, payload.ChaincodeAuctionName
}

// Reads from a json the broker, port and topic to perform MQTT communications
func readMQTTConfig() (string, int, string, string, string, string) {
	content, err := ioutil.ReadFile("./mqtt_config.json")
	if err != nil {
		panic(fmt.Errorf("Error when opening file: %w", err))
	}

	var payload MQTTConfig
	err = json.Unmarshal(content, &payload)
	if err != nil {
		panic(fmt.Errorf("Error during Unmarshal(): %w", err))
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

// MQTT UTILS HANDLERS

// Receives the message from the smart meter and perform some actions
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	var message MQTTmessage
	err := json.Unmarshal(msg.Payload(), &message)
	if err != nil {
		panic(fmt.Errorf("Error during Unmarshal(): %w", err))
	}

	readEnergyAssetByID(contract_energy)
	// Prendo l'output: se assetID non esiste allora stop e vado avanti
	//					se assetID esiste invece, registro la quantità rimanente e poi chiamo il delete
	// if <asset-esiste> {
	//		salvo la quantità da qualche parte ed eventualmente la invio a chi di dovere
	deleteEnergyAsset(contract_energy)
	// }

	if message.ProducedEnergy > message.ConsumedEnergy {
		production_excess := message.ProducedEnergy - message.ConsumedEnergy
		s_excess := fmt.Sprintf("%v", production_excess)

		createEnergyAsset(contract_energy, s_excess)

		// Prima creo l'asset e poi indico l'asta
	} else if message.ProducedEnergy < message.ConsumedEnergy {
		consumption_excess := message.ConsumedEnergy - message.ProducedEnergy

		// Devo far sapere in qualche modo al listener dell'asta quanta energia ho bisogno di comprare (variabile globale?)
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

// (communicates an error to mobile application?)
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

// REST FUNCTIONS (GET ad POST)

type inputPOSTprimo struct {
}

type inputPOSTsecondo struct {
}

type outputGETprimo struct {
}

type outputGETsecondo struct {
}

// getPRIMO responds with
func getPRIMO(c *gin.Context) {
	var output outputGETprimo

	c.IndentedJSON(http.StatusOK, output)
}

// getSECONDO responds with
func getSECONDO(c *gin.Context) {
	var output outputGETsecondo

	c.IndentedJSON(http.StatusOK, output)
}

// postPRIMO ... from JSON received in the request body.
func postPRIMO(c *gin.Context) {
	var input inputPOSTprimo

	// Call BindJSON to bind the received JSON to input
	if err := c.BindJSON(&input); err != nil {
		return
	}

	// Add the new album to the slice.
	// albums = append(albums, newAlbum)

	c.IndentedJSON(http.StatusCreated, input)
}

// postSECONDO ... from JSON received in the request body.
func postSECONDO(c *gin.Context) {
	var input inputPOSTsecondo

	// Call BindJSON to bind the received JSON to input
	if err := c.BindJSON(&input); err != nil {
		return
	}

	// Add the new album to the slice.
	// albums = append(albums, newAlbum)

	c.IndentedJSON(http.StatusCreated, input)
}
