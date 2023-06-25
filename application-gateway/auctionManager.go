package main

import (
	"fmt"
	"strconv"
	"time"
)

func manageSeller(quantity string) {

	// Prima creo l'asset e poi indico l'asta

	createEnergyAsset(contract_energy, quantity)
	createAuction(contract_auction, "energy", fmt.Sprintf("%v", moneySettings.startingBid), quantity)

	// I wait 30 seconds before closing the auction (the other members have 30 seconds to make their bid)
	time.Sleep(30 * time.Second)
	closeAuction(contract_auction)

	// I wait another 30 seconds before ending the auction (the other members have 30 seconds to reveal their bid)
	time.Sleep(30 * time.Second)
	endAuction(contract_auction)
	// When I call endAuction that function automatically copies the action in the global variable "Auction"
	// therefore we can access it to take the informations needed

	// Wait for the winners to receive the event and make their operations
	time.Sleep(5 * time.Second)

	// CHECK IF THERE ARE SOME WINNERS OR NOT

	// IF THERE ARE NO WINNERS
	// 	- make this function as a loop that we can start again if there are no winners (and we should end if there are winners)
	//	- start the <createAuctionListener>

	// IF THERE ARE WINNERS GO ON

	// Variable used to manage the asset transfer
	var winnersEnergyQuantities []float64

	// Variables used to manage the money transfer
	totalMoney := 0.0
	var winnersMoneyQuantities []float64

	// Directly send the assets to the winners
	fmt.Printf("Transfering the assets to the winners\n")
	for _, winner := range auctionStruct.Winners {

		transferEnergyAssetAsync(contract_energy, winner.Buyer, fmt.Sprintf("%v", winner.Quantity))
		winnersEnergyQuantities = append(winnersEnergyQuantities, float64(-winner.Quantity))

		// I create the array of money transactions that should happen
		winnersMoneyQuantities = append(winnersMoneyQuantities, (float64(auctionStruct.Price) * float64(winner.Quantity)))
		totalMoney += float64(auctionStruct.Price) * float64(winner.Quantity)
	}

	// Create the energy transaction struct and register the transactions in the array for the GET interface
	energyQuantities_struct := &EnergyTransactions_struct{
		Transactions: winnersEnergyQuantities,
		Timestamp:    time.Now().String(),
	}

	energyTransactionsArray.TransactionsArray = append(energyTransactionsArray.TransactionsArray, energyQuantities_struct)

	// Cicles and continues until the current money balance is equal to the starting balance + the total money
	// that the seller should receive
	currentBalance, _ := strconv.ParseFloat(getAccountBalance(contract_money), 64)
	objectiveBalance := currentBalance + totalMoney
	for currentBalance != objectiveBalance {

		currentBalance, _ = strconv.ParseFloat(getAccountBalance(contract_money), 64)
	}

	// Create the money transaction struct and register the transactions in the array for the GET interface
	moneyTransactionsStruct := &MoneyTransactions_struct{
		Transactions: winnersMoneyQuantities,
		Price:        auctionStruct.Price,
		Timestamp:    time.Now().String(),
	}

	moneyTransactionsArray.TransactionsArray = append(moneyTransactionsArray.TransactionsArray, moneyTransactionsStruct)
}
