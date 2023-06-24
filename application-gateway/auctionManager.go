package main

import (
	"context"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func manageSeller(quantity string) {

	// Prima creo l'asset e poi indico l'asta

	createEnergyAsset(contract_energy, quantity)
	createAuction(contract_auction, "energy", moneySettings.startingBid, quantity)

	// I wait 30 seconds before closing the auction (the other members have 30 seconds to make their bid)
	time.Sleep(30 * time.Second)
	closeAuction(contract_auction)

	// I wait another 30 seconds before ending the auction (the other members have 30 seconds to reveal their bid)
	time.Sleep(30 * time.Second)
	endAuction(contract_auction)
	// When I call endAuction that function automatically copies the action in the global variable "Auction"
	// therefore we can access it to take the informations needed

	// Directly send the assets to the winners
	fmt.Printf("Transfering the assets to the winners\n")
	for _, winner := range auctionStruct.Winners {

		transferEnergyAssetAsync(contract_energy, winner.Buyer, fmt.Sprintf("%v", winner.Quantity))
	}

}

func managerBuyer(ctx context.Context, network *client.Network, chaincodeName string) {

}
