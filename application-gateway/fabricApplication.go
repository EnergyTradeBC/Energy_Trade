/*
Copyright 2021 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type NetworkConfig struct {
	ChannelName         string
	ChaincodeMoneyName  string
	ChaincodeEnergyName string
}

type MQTTConfig struct {
	broker   string
	port     int
	topic    string
	clientID string
	username string
	password string
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
	gw, err := client.Connect(
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
	defer gw.Close()

	// Read network information from a json file
	channelName, chaincodeMoneyName, chaincodeEnergyName := readNetworkConfig()

	// Get network and smart contract objects
	network := gw.GetNetwork(channelName)
	contract_money := network.GetContract(chaincodeMoneyName)
	contract_energy := network.GetContract(chaincodeEnergyName)

	for true {

	}

}

// UTILITY FUNCTIONS

// Reads from a json the channel name and the name of the chaincodes committed to that channel
func readNetworkConfig() (string, string, string) {
	content, err := ioutil.ReadFile("./network_config.json")
	if err != nil {
		panic(fmt.Errorf("Error when opening file: %w", err))
	}

	var payload NetworkConfig
	err = json.Unmarshal(content, &payload)
	if err != nil {
		panic(fmt.Errorf("Error during Unmarshal(): %w", err))
	}

	return payload.ChannelName, payload.ChaincodeMoneyName, payload.ChaincodeEnergyName
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

// FUNCTIONS TO MANAGE THE ENERGY CONTRACT

// Evaluate a transaction to query ledger state.
func getAllEnergyAssets(contract *client.Contract) {
	fmt.Println("\n--> Evaluate Transaction: GetAllEnergyAssets, function returns all the current energy assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

// Evaluate a transaction by assetID to query ledger state.
func readEnergyAssetByID(contract *client.Contract, asset_ID string) {
	fmt.Printf("\n--> Evaluate Transaction: ReadAsset, function returns asset attributes\n")

	evaluateResult, err := contract.EvaluateTransaction("ReadAsset", asset_ID)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

// Submit a transaction synchronously, blocking until it has been committed to the ledger.
func createEnergyAsset(contract *client.Contract, asset_ID string, quantity string) {
	fmt.Printf("\n--> Submit Transaction: CreateEnergyAsset, creates new energy asset with asset_ID and quantity arguments \n")

	_, err := contract.SubmitTransaction("CreateAsset", asset_ID, quantity)
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Submit transaction asynchronously, blocking until the transaction has been sent to the orderer, and allowing
// this thread to process the chaincode response (e.g. update a UI) without waiting for the commit notification
func transferAssetAsync(contract *client.Contract, asset_ID string, newOwner_ID string, transfer_quantity string) {
	fmt.Printf("\n--> Async Submit Transaction: TransferAsset, transfer part or the entire energy asset to a new owner")

	submitResult, commit, err := contract.SubmitAsync("TransferAsset", client.WithArguments(asset_ID, newOwner_ID, transfer_quantity))
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction asynchronously: %w", err))
	}

	fmt.Printf("\n*** Successfully submitted transaction to transfer ownership from %s to Mark. \n", string(submitResult))
	fmt.Println("*** Waiting for transaction commit.")

	if commitStatus, err := commit.Status(); err != nil {
		panic(fmt.Errorf("failed to get commit status: %w", err))
	} else if !commitStatus.Successful {
		panic(fmt.Errorf("transaction %s failed to commit with status: %d", commitStatus.TransactionID, int32(commitStatus.Code)))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Submit a transaction synchronously, blocking until it has been committed to the ledger.
func deleteEnergyAsset(contract *client.Contract, asset_ID string) {
	fmt.Printf("\n--> Submit Transaction: DeleteEnergyAsset, deletes an energy asset using its asset_ID \n")

	_, err := contract.SubmitTransaction("DeleteAsset", asset_ID)
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
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

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
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
