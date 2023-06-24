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
		p_excess := fmt.Sprintf("%v", production_excess)

		fmt.Printf("Production excess: %s \n", p_excess)
		// I leave the time to the buyers to "activate" the listener
		time.Sleep(5 * time.Second)

		// Call the seller manager
		manageSeller(p_excess)

	} else if message.ProducedEnergy < message.ConsumedEnergy {
		consumption_excess := message.ConsumedEnergy - message.ProducedEnergy
		c_excess = fmt.Sprintf("%v", consumption_excess)

		// Devo far sapere in qualche modo al listener dell'asta quanta energia ho bisogno di comprare (variabile globale?)

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

func publish(client mqtt.Client, topic string, message string) {
	token := client.Publish(topic, 0, false, message)
	token.Wait()
	time.Sleep(time.Second)
}
