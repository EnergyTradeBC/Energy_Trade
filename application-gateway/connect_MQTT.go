package main

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTConfig struct {
	broker   string
	port     int
	topic    string
	clientID string
	username string
	password string
}

type MQTTmessage struct {
	Timestamp      string
	ProducedEnergy float64
	ConsumedEnergy float64
}

type RemainingAsset struct {
	Timestamp string
	Quantity  string
}

// Variable used to register the status of the client inside the auction (buyer or seller)
var clientAuctionStatus = "buyer"

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	var message MQTTmessage
	err := json.Unmarshal(msg.Payload(), &message)
	if err != nil {
		panic(fmt.Errorf("error during Unmarshal(): %w", err))
	}

	remainingQuantity := readEnergyAssetByID(contract_energy)

	// If there is some remaining quantity I save the remaining quantity and then delete the asset
	if remainingQuantity != "failed" {

		// If the client in the previous time slot was a seller we can add the remaining assets to the array
		// If the client was a buyer, the quantity that remains cannot be sold to the accumulator or donated so it must deleted only
		if clientAuctionStatus == "seller" {

			remainingAsset := &RemainingAsset{
				Timestamp: time.Now().String(),
				Quantity:  remainingQuantity,
			}

			// COME GESTIIAMO IL NODO DELLE DONAZIONI? E' UN PEER FISSO DI CUI OGNI CLIENT CONOSCE GLI ENDPOINT E
			// RICEVE UNA PERCENTUALE DEGLI ASSET CHE AVANZANO? (CHE PERCENTUALE?)

			remainingAssetsArray.AssetArray = append(remainingAssetsArray.AssetArray, remainingAsset)
		}

		deleteEnergyAsset(contract_energy)
	}

	if message.ProducedEnergy > message.ConsumedEnergy {
		production_excess := message.ProducedEnergy - message.ConsumedEnergy
		p_excess := fmt.Sprintf("%v", production_excess)

		fmt.Printf("Production excess: %s \n", p_excess)
		// I leave the time to the buyers to "activate" the listener
		time.Sleep(5 * time.Second)

		// Set the client status to seller
		clientAuctionStatus = "seller"

		// Call the seller manager
		manageSeller(p_excess)

	} else if message.ProducedEnergy < message.ConsumedEnergy {
		consumption_excess := message.ConsumedEnergy - message.ProducedEnergy
		c_excess = fmt.Sprintf("%v", consumption_excess)

		fmt.Printf("Consumption excess: %s \n", c_excess)

		// Set the client status to buyer
		clientAuctionStatus = "buyer"

		// With the global variable I make the listener aware of the fact that it must catch createAuction events
		// The buyer operations are managed inside the listeners
		createAuctionListener = true
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

// (communicates an error to mobile application?)
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
}

func sub_topic(client mqtt.Client, topic string) {
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s", topic)
}

// func publish(client mqtt.Client, topic string, message string) {
// 	token := client.Publish(topic, 0, false, message)
// 	token.Wait()
// 	time.Sleep(time.Second)
// }
