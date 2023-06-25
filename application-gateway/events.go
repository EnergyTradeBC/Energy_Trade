package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

var createAuctionListener = false
var closeAuctionListener = false
var endAuctionListener = false

var transferEnergyListener = false

// we can use this variable to check if it is received the correct quantity of asset
// var transferEnergyRequired = 0.0

// func startBlockEventListening(ctx context.Context, network *client.Network) {
// 	fmt.Println("\n*** Start block event listening")

// 	events, err := network.BlockEvents(ctx)
// 	if err != nil {
// 		panic(fmt.Errorf("failed to start chaincode event listening: %w", err))
// 	}

// 	go func() {
// 		for event := range events {

// 			// Come si gestiscono le variabili di tipo [][]byte e cosa c'è all'interno? sono una cazzo di matriosca di blocchi
// 			// https://github.com/hyperledger/fabric-protos-go-apiv2/blob/v0.3.0/common/common.pb.go#L970

// 			asset := formatJSON(event.Data.Data)
// 			fmt.Printf("\n<-- Energy chaincode event received: %s - %s\n", event.EventName, asset)

// 			// per TESTING (?) e vedere come sono realmente i blocchi che vengono scambiati nella blockchain

// 		}
// 	}()
// }

func startEnergyChaincodeEventListening(ctx context.Context, network *client.Network, chaincodeName string) {
	fmt.Println("\n*** Start energy chaincode event listening")

	events, err := network.ChaincodeEvents(ctx, chaincodeName)
	if err != nil {
		panic(fmt.Errorf("failed to start chaincode event listening: %w", err))
	}

	// close(events)
	// CAN WE CLOSE EVENTS OR ONCE THEY ARE OPEN THEY CANNOT BE CLOSED?
	// If we cannot close them we must manage the reception of events in a different way (global variables? 0 listen, 1 not listen)

	go func() {
		for event := range events {

			clientID := getAccountID(contract_auction)

			var energyEvent EnergyEvent
			err = json.Unmarshal(event.Payload, &energyEvent)
			if err != nil {
				panic(fmt.Errorf("failed to start chaincode event listening: %w", err))
			}
			fmt.Printf("\n<-- Energy chaincode event received: %s - %s\n", event.EventName, event.Payload)

			if event.EventName == "TransferAsset" && transferEnergyListener {

				if energyEvent.New_OwnerID == clientID {

					// CONTROLLARE CHE LA QUANTITA' CHE E' STATA INVIATA SIA EFFETTIVAMENTE QUELLA CHE CI SI ASPETTAVA

					// INVECE CHE GESTIRE QUESTA OPERAZIONE CON GLI EVENTI POTREI FARLO COME HO FATTO PER IL DENARO
					// NEL CASO DEL SELLER, OVVERO FARE DELLE QUERY ALL'ASSET FINTANTO CHE NON DIVENTA QUELLO CHE MI ASPETTAVO

					var energyQuantity []float64
					energyQuantity = append(energyQuantity, float64(energyEvent.Quantity))

					// Create the energy transaction struct and register the transactions in the array for the GET interface
					energyQuantities_struct := &EnergyTransactions_struct{
						Transactions: energyQuantity,
						Timestamp:    time.Now().String(),
					}

					energyTransactionsArray.TransactionsArray = append(energyTransactionsArray.TransactionsArray, energyQuantities_struct)

					transferEnergyListener = false
				}
			}
		}
	}()
}

func startAuctionChaincodeEventListening(ctx context.Context, network *client.Network, chaincodeName string) {
	fmt.Println("\n*** Start auction chaincode event listening")

	events, err := network.ChaincodeEvents(ctx, chaincodeName)
	if err != nil {
		panic(fmt.Errorf("failed to start chaincode event listening: %w", err))
	}

	go func() {
		// Create variables to store the bid ID and list of available auctions
		// type transientAuction struct {
		// 	ID       string
		// 	Quantity float32
		// }
		var bidID string
		// var auctionList []*transientAuction

		for event := range events {
			var auctionEvent AuctionEvent
			err = json.Unmarshal(event.Payload, &auctionEvent)
			if err != nil {
				panic(fmt.Errorf("failed to start chaincode event listening: %w", err))
			}
			fmt.Printf("\n<-- Auction chaincode event received: %s - %s\n", event.EventName, event.Payload)

			if event.EventName == "createAuction" && createAuctionListener {

				// WE SHOULD WAIT FOR MORE THAN 1 AUCTION AND THEN DECIDE TO MAKE BID IN THE BEST ONE
				// FOR NOW THE CODE IS MADE FOR JOINING THE FIRST AUCTION THAT THE LISTENER SENSE

				// if auctionEvent.StartingBid < moneySettings.maxPrizeAuction {

				// 	temp := &transientAuction{
				// 		ID:       auctionEvent.AuctionID,
				// 		Quantity: auctionEvent.Quantity,
				// 	}
				// 	auctionList = append(auctionList, temp)
				// }

				// As it is now we wait for at least 2 auctions
				// or we wait at least X seconds (?????)
				// if len(auctionList) == 2 {

				// }

				if auctionEvent.StartingBid < moneySettings.maxPrizeAuction {

					createAuctionListener = false

					// Make and then submits a bid to the chosen auction using the AuctionID and save the bid ID for the reveal
					bidID = makeBid(contract_auction, auctionEvent.AuctionID, fmt.Sprintf("%v", auctionEvent.Quantity), fmt.Sprintf("%v", moneySettings.maxPrizeAuction))
					submitBid(contract_auction, auctionEvent.AuctionID, bidID)

					// QUERY DELL'ASTA PER CONFERMARE LA BID?

					closeAuctionListener = true
				}

			} else if event.EventName == "closeAuction" && closeAuctionListener {

				closeAuctionListener = false

				// Reveal the bid created previously using its ID
				revealBid(contract_auction, auctionEvent.AuctionID, bidID)

				endAuctionListener = true

			} else if event.EventName == "endAuction" && endAuctionListener {

				endAuctionListener = true

				// Retrieve the ID of the client to check if it is the winner
				buyerID, err := contract_auction.EvaluateTransaction("GetSubmittingClientIdentity")
				if err != nil {
					panic(fmt.Errorf("failed to evaluate transaction: %w", err))
				}

				// Variable used as a flag for client victory
				quantity := float32(0)

				for _, winner := range auctionEvent.Winners {

					if winner.Buyer == string(buyerID) {
						quantity = winner.Quantity
					}
				}

				if quantity != 0 {

					time.Sleep(3 * time.Second)

					moneyAmount := auctionEvent.Price * quantity
					transferMoney(contract_auction, auctionEvent.Seller, fmt.Sprintf("%v", moneyAmount))

					var moneyAmountArray []float64
					moneyAmountArray = append(moneyAmountArray, float64(-moneyAmount))

					// Create the money transaction struct and register the transactions in the array for the GET interface
					moneyTransactionsStruct := &MoneyTransactions_struct{
						Transactions: moneyAmountArray,
						Price:        auctionEvent.Price,
						Timestamp:    time.Now().String(),
					}

					moneyTransactionsArray.TransactionsArray = append(moneyTransactionsArray.TransactionsArray, moneyTransactionsStruct)

					// Activate the listener for energy transfers
					// transferEnergyRequired = float64(quantity)
					transferEnergyListener = true

				} else {

					// If the client is not a winner activate again the <createAuctionListener> to join new auctions
					createAuctionListener = true
				}

			}

		}
	}()
}
