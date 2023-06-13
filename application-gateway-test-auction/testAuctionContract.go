/*
Copyright 2021 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Auction data
type Auction struct {
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

const (
	mspID        = "Org1MSP"
	cryptoPath   = "../ENERGY-TRADE/organizations/peerOrganizations/org1.example.com"
	certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/cert.pem"
	keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore/"
	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint = "localhost:7051"
	gatewayPeer  = "peer0.org1.example.com"
	assetID      = "energy_1"
	auctionID    = "auction_1"
)

func main() {
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

	// Override default values for chaincode and channel name as they may differ in testing contexts.
	chaincodeName := "chaincodetest"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	channelName := "channeltest"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	createAuction(contract, "energy", "5", "15")
	queryAuctionByID(contract, auctionID)
	// bidID := makeBid(contract, auctionID, "15", "10")
	// submitBid(contract, auctionID, bidID)
	// queryBidByID(contract, auctionID, bidID)
	closeAuction(contract)
	// revealBid(contract, auctionID, bidID)
	endAuction(contract)
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection() *grpc.ClientConn {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity() *identity.X509Identity {
	certificate, err := loadCertificate(certPath)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign() identity.Sign {
	files, err := os.ReadDir(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}
	privateKeyPEM, err := os.ReadFile(path.Join(keyPath, files[0].Name()))

	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

// Evaluate a transaction to query ledger state.
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
		client.WithArguments(auctionID, item, startingBid, quantity))

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
		"org":        []byte(mspID),
		"buyer":      []byte(string(buyerID)),
	}

	fmt.Printf("\n--> Submit Transaction: makeBid, creates a new bid \n")

	bidID, err := contract.Submit("CreateAuction",
		client.WithArguments(auction_ID),
		client.WithTransient(privateBid),
		client.WithEndorsingOrganizations(mspID))

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

	var auctionStruct Auction
	err = json.Unmarshal(auctionInfo, &auctionStruct)
	if err != nil {
		panic(fmt.Errorf("Error during Unmarshal(): %w", err))
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

	var auctionStruct Auction
	err = json.Unmarshal(auctionInfo, &auctionStruct)
	if err != nil {
		panic(fmt.Errorf("Error during Unmarshal(): %w", err))
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

	var auctionStruct Auction
	err = json.Unmarshal(auctionInfo, &auctionStruct)
	if err != nil {
		panic(fmt.Errorf("Error during Unmarshal(): %w", err))
	}

	fmt.Printf("\n--> Evaluate Transaction: get bid information\n")

	bidInfo, err := contract.EvaluateTransaction("QueryBid", bidID)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}

	var bidStruct FullBid
	err = json.Unmarshal(bidInfo, &bidStruct)
	if err != nil {
		panic(fmt.Errorf("Error during Unmarshal(): %w", err))
	}

	// Build the transient map with the bid information
	privateBid := map[string][]byte{
		"objectType": []byte("bid"),
		"quantity":   []byte(bidStruct.Quantity),
		"price":      []byte(bidStruct.Price),
		"org":        []byte(mspID),
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

	var auctionStruct Auction
	err = json.Unmarshal(auctionInfo, &auctionStruct)
	if err != nil {
		panic(fmt.Errorf("Error during Unmarshal(): %w", err))
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

// Format JSON data
func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}
