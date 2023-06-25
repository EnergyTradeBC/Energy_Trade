package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// Auction data
type Auction struct {
	AuctionID    string             `json:"auctionID"`
	Type         string             `json:"objectType"`
	ItemSold     string             `json:"item"`
	Seller       string             `json:"seller"`
	Quantity     float32            `json:"quantity"`
	Orgs         []string           `json:"organizations"`
	PrivateBids  map[string]BidHash `json:"privateBids"`
	RevealedBids map[string]FullBid `json:"revealedBids"`
	Winners      []Winners          `json:"winners"`
	Price        float32            `json:"price"`
	Status       string             `json:"status"`
	Auditor      bool               `json:"auditor"`
	StartingBid  float32            `json:"startingBid"`
}

// event provides an organized struct for emitting events
type AuctionEvent struct {
	AuctionID   string    `json:"auctionID"`
	ItemSold    string    `json:"item"`
	StartingBid float32   `json:"startingBid"`
	Seller      string    `json:"seller"`
	Quantity    float32   `json:"quantity"`
	Winners     []Winners `json:"winners"`
	Price       float32   `json:"price"`
	Status      string    `json:"status"`
}

// FullBid is the structure of a revealed bid
// per il momento sono tutti "string" perché ho problemi con la gestione di valori numerici (non gli vanno mai bene alle API per gli smart contract)
type FullBid struct {
	Type     string `json:"objectType"`
	Quantity string `json:"quantity"`
	Price    string `json:"price"`
	Org      string `json:"org"`
	Buyer    string `json:"buyer"`
}

// BidHash is the structure of a private bid
type BidHash struct {
	Org  string `json:"org"`
	Hash string `json:"hash"`
}

// Winners stores the winners of the auction
type Winners struct {
	Buyer    string  `json:"buyer"`
	Quantity float32 `json:"quantity"`
}

var auctionStruct Auction

func queryAuctionByID(contract *client.Contract, auction_ID string) string {
	fmt.Println("\n--> Evaluate Transaction: queryAuctionByID, function returns all the information about the auction with <auctionID>")

	evaluateResult, err := contract.EvaluateTransaction("QueryAuction", auction_ID)

	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)

	return result
}

func queryBidByID(contract *client.Contract, auction_ID string, bidID string) string {
	fmt.Println("\n--> Evaluate Transaction: queryBidByID, function returns all the information about the bid with <BidID>")

	evaluateResult, err := contract.EvaluateTransaction("QueryBid", auction_ID, bidID)

	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)

	return result
}

// Submit a transaction synchronously, blocking until it has been committed to the ledger.
func createAuction(contract *client.Contract, item string, startingBid string, quantity string) {
	fmt.Printf("\n--> Submit Transaction: createAuction, creates a new auction with auctionID, item, prize and quantity arguments \n")

	//_, err := contract.SubmitTransaction("CreateAuction", auctionID, item, price, quantity)
	_, err := contract.Submit("CreateAuction",
		client.WithArguments(auctionID, item, quantity, startingBid))

	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Create a bid and submit it then returns the bid ID if gone correctly.
func makeBid(contract *client.Contract, auction_ID string, quantity string, price string) string {
	fmt.Printf("\n--> Evaluate Transaction: get your client ID \n")

	buyerID, err := contract.EvaluateTransaction("GetSubmittingClientIdentity")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}

	// Build the transient map with the bid information
	privateBid := map[string][]byte{
		"objectType": []byte("bid"),
		"quantity":   []byte(quantity),
		"price":      []byte(price),
		"org":        []byte(orgMSP),
		"buyer":      []byte(string(buyerID)),
	}

	fmt.Printf("\n--> Submit Transaction: makeBid, creates a new bid \n")

	bidID, err := contract.Submit("Bid",
		client.WithArguments(auction_ID),
		client.WithTransient(privateBid),
		client.WithEndorsingOrganizations(orgMSP))

	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully with bid ID %s\n", string(bidID))

	return string(bidID)
}

// Submit the bid created previously adding to the endorsing orgs only the ones already present in the auction.
func submitBid(contract *client.Contract, auction_ID string, bidID string) {
	fmt.Printf("\n--> Evaluate Transaction: get auction information \n")

	auctionInfo, err := contract.EvaluateTransaction("QueryAuction", auction_ID)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}

	err = json.Unmarshal(auctionInfo, &auctionStruct)
	if err != nil {
		panic(fmt.Errorf("error during Unmarshal(): %w", err))
	}

	fmt.Printf("\n--> Submit Transaction: submitBid, submits the bid created previously with bidID \n")

	_, err = contract.Submit("SubmitBid",
		client.WithArguments(auction_ID, bidID),
		client.WithEndorsingOrganizations(auctionStruct.Orgs...))
	// I 3 punti dopo una lista fanno in modo di spacchettare gli elementi all'interno e passarli uno per volta invece che come una lista
	// in questo modo dico al comando che l'endorsement deve essere fatto solo dalle org che hanno preso parte all'asta

	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Close the auction blocking new bids and adding to the endorsing peers only those who partecipated to the auction.
func closeAuction(contract *client.Contract) {
	fmt.Printf("\n--> Evaluate Transaction: get auction information \n")

	auctionInfo, err := contract.EvaluateTransaction("QueryAuction", auctionID)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}

	err = json.Unmarshal(auctionInfo, &auctionStruct)
	if err != nil {
		panic(fmt.Errorf("error during Unmarshal(): %w", err))
	}

	fmt.Printf("\n--> Submit Transaction: closeAuction, closes the auction blocking bids \n")

	_, err = contract.Submit("CloseAuction",
		client.WithArguments(auctionID),
		client.WithEndorsingOrganizations(auctionStruct.Orgs...))
	// I 3 punti dopo una lista fanno in modo di spacchettare gli elementi all'interno e passarli uno per volta invece che come una lista
	// in questo modo dico al comando che l'endorsement deve essere fatto solo dalle org che hanno preso parte all'asta

	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Reveals the bid with ID <bidID> adding as endorsement orgs the ones partecipating the auction.
func revealBid(contract *client.Contract, auction_ID string, bidID string) {
	fmt.Printf("\n--> Evaluate Transaction: get auction information\n")

	auctionInfo, err := contract.EvaluateTransaction("QueryAuction", auction_ID)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}

	err = json.Unmarshal(auctionInfo, &auctionStruct)
	if err != nil {
		panic(fmt.Errorf("error during Unmarshal(): %w", err))
	}

	fmt.Printf("\n--> Evaluate Transaction: get bid information\n")

	bidInfo, err := contract.EvaluateTransaction("QueryBid", bidID)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}

	var bidStruct FullBid
	err = json.Unmarshal(bidInfo, &bidStruct)
	if err != nil {
		panic(fmt.Errorf("error during Unmarshal(): %w", err))
	}

	// Build the transient map with the bid information
	privateBid := map[string][]byte{
		"objectType": []byte("bid"),
		"quantity":   []byte(bidStruct.Quantity),
		"price":      []byte(bidStruct.Price),
		"org":        []byte(orgMSP),
		"buyer":      []byte(bidStruct.Buyer),
	}

	fmt.Printf("\n--> Submit Transaction: revealBid, reveals the bid submitted previously with bidID \n")

	_, err = contract.Submit("SubmitBid",
		client.WithArguments(auction_ID, bidID),
		client.WithTransient(privateBid),
		client.WithEndorsingOrganizations(auctionStruct.Orgs...))

	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Ends the auction electing the winners and adding as endorsement orgs only those who partecipated.
func endAuction(contract *client.Contract) {
	fmt.Printf("\n--> Evaluate Transaction: get auction information \n")

	auctionInfo, err := contract.EvaluateTransaction("QueryAuction", auctionID)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}

	err = json.Unmarshal(auctionInfo, &auctionStruct)
	if err != nil {
		panic(fmt.Errorf("error during Unmarshal(): %w", err))
	}

	fmt.Printf("\n--> Submit Transaction: endAuction, ends the auction and elect the winners \n")

	_, err = contract.Submit("EndAuction",
		client.WithArguments(auctionID),
		client.WithEndorsingOrganizations(auctionStruct.Orgs...))

	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}
