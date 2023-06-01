package main

import (
	"context"
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func startBlockEventListening(ctx context.Context, network *client.Network) {
	fmt.Println("\n*** Start block event listening")

	events, err := network.BlockEvents(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to start chaincode event listening: %w", err))
	}

	go func() {
		for event := range events {

			// Come si gestiscono le variabili di tipo [][]byte e cosa c'è all'interno? sono una cazzo di matriosca di blocchi
			// https://github.com/hyperledger/fabric-protos-go-apiv2/blob/v0.3.0/common/common.pb.go#L970

			asset := formatJSON(event.Data.Data)
			fmt.Printf("\n<-- Energy chaincode event received: %s - %s\n", event.EventName, asset)

			// per TESTING (?) e vedere come sono realmente i blocchi che vengono scambiati nella blockchain

		}
	}()
}

func startEnergyChaincodeEventListening(ctx context.Context, network *client.Network, chaincodeName string) {
	fmt.Println("\n*** Start energy chaincode event listening")

	events, err := network.ChaincodeEvents(ctx, chaincodeName)
	if err != nil {
		panic(fmt.Errorf("failed to start chaincode event listening: %w", err))
	}

	go func() {
		for event := range events {
			asset := formatJSON(event.Payload)
			fmt.Printf("\n<-- Energy chaincode event received: %s - %s\n", event.EventName, asset)

			// per TESTING o per triggerare l'invio del denaro (se decidiamo che prima viene trasferito l'asset e poi il denaro)

			if event.EventName == "TransferAsset" {
				// Saranno necessari altri controlli ma prima bisogna vedere com'è fatto "event" e "event.Payload"

				// LOGICA PER IL TRASFERIMENTO DEL DENARO AL TERMINE DELL'ASTA (appena ricevo l'asset che era stato stabilito triggero la
				// funzione dello smart contract del denaro => trasferisco il denaro al peer da cui ho ricevuto l'asset)
			}
		}
	}()
}

func startMoneyChaincodeEventListening(ctx context.Context, network *client.Network, chaincodeName string) {
	fmt.Println("\n*** Start monry chaincode event listening")

	events, err := network.ChaincodeEvents(ctx, chaincodeName)
	if err != nil {
		panic(fmt.Errorf("failed to start chaincode event listening: %w", err))
	}

	go func() {
		for event := range events {
			asset := formatJSON(event.Payload)
			fmt.Printf("\n<-- Money chaincode event received: %s - %s\n", event.EventName, asset)

			// per TESTING o per triggerare l'invio dell'energia (se decidiamo che prima viene trasferito il denaro e poi l'energia)

			if event.EventName == "Transfer" {
				// Saranno necessari altri controlli ma prima bisogna vedere com'è fatto "event" e "event.Payload"
				// A maggior ragione nel caso del denaro perché le transazioni dovrebbero essere private => come funzia?
				// Il proprietario del denaro o comunque entrambi i due lati nel trasferimento dovrebbero poterlo vedere in chiaro invece (?)

				// LOGICA PER IL TRASFERIMENTO DELL'ASSET AL TERMINE DELL'ASTA (appena ricevo il denaro che era stato stabilito triggero la
				// funzione dello smart contract dell'energia  => trasferisco l'energia al peer da cui ho ricevuto il denaro)
			}
		}
	}()
}

func startAuctionChaincodeEventListening(ctx context.Context, network *client.Network, chaincodeName string) {
	fmt.Println("\n*** Start chaincode event listening")

	events, err := network.ChaincodeEvents(ctx, chaincodeName)
	if err != nil {
		panic(fmt.Errorf("failed to start chaincode event listening: %w", err))
	}

	go func() {
		for event := range events {
			asset := formatJSON(event.Payload)
			fmt.Printf("\n<-- Auction chaincode event received: %s - %s\n", event.EventName, asset)
		}
	}()
}
