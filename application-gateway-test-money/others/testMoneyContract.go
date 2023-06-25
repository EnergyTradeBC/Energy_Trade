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

const (
	mspID         = "Org1MSP" //"Org2MSP"
	cryptoPath    = "../ENERGY-TRADE/organizations/peerOrganizations/org1.example.com"
	certPath      = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem"
	keyPath       = cryptoPath + "/users/User1@org1.example.com/msp/keystore/"
	tlsCertPath   = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint  = "localhost:7051"
	gatewayPeer   = "peer0.org1.example.com"
	moneyName     = "Token"
	moneySymbol   = "[--]"
	moneyDecimals = "0"
	id2           = "eDUwOTo6Q049VXNlcjFAb3JnMi5leGFtcGxlLmNvbSxPVT1jbGllbnQsTD1TYW4gRnJhbmNpc2NvLFNUPUNhbGlmb3JuaWEsQz1VUzo6Q049Y2Eub3JnMi5leGFtcGxlLmNvbSxPPW9yZzIuZXhhbXBsZS5jb20sTD1TYW4gRnJhbmNpc2NvLFNUPUNhbGlmb3JuaWEsQz1VUw=="
	id3           = "eDUwOTo6Q049VXNlcjFAb3JnMy5leGFtcGxlLmNvbSxPVT1jbGllbnQsTD1TYW4gRnJhbmNpc2NvLFNUPUNhbGlmb3JuaWEsQz1VUzo6Q049Y2Eub3JnMy5leGFtcGxlLmNvbSxPPW9yZzMuZXhhbXBsZS5jb20sTD1TYW4gRnJhbmNpc2NvLFNUPUNhbGlmb3JuaWEsQz1VUw=="
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
	chaincodeName := "moneyAsset"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	channelName := "channeltest"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	//initializeMoneyContract(contract)
	getAccountID(contract)
	//getAccountBalance(contract)
	mint(contract, id2, "200")
	mint(contract, id3, "500")
	//getAccountBalance(contract)
	//transferMoney(contract, id3, "100")
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

// Initializes the money smart contract passing 3 settings (name, symbol and decimals).
// Can be used only by the org that will behave as the central banker
func initializeMoneyContract(contract *client.Contract) {
	fmt.Printf("\n--> Submit Transaction: initializeMoneyContract, initializes the money contract with 'name', 'symbol' and 'decimals' \n")

	_, err := contract.Submit("Initialize",
		client.WithArguments(moneyName, moneySymbol, moneyDecimals))

	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Evaluate a transaction requesting the client ID of the caller.
func getAccountID(contract *client.Contract) string {
	fmt.Printf("\n--> Evaluate Transaction: getAccountID, retrieves the client ID of the caller \n")

	evaluateResult, err := contract.EvaluateTransaction("ClientAccountID")

	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	// Rimane però codificato BASE-64, va bene?
	result := string(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)

	return result
}

// Evaluate a transaction requesting the money balance of the caller.
func getAccountBalance(contract *client.Contract) string {
	fmt.Printf("\n--> Evaluate Transaction: getAccountID, retrieves the client ID of the caller \n")

	evaluateResult, err := contract.EvaluateTransaction("ClientAccountBalance")

	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := string(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)

	return result
}

// Creates <amount> of new tokens.
// Can be used only by the org that will behave as the central banker
func mint(contract *client.Contract, recipientID string, amount string) {
	fmt.Printf("\n--> Submit Transaction: mint, creates new tokens for the central banker \n")

	_, err := contract.Submit("Mint",
		client.WithArguments(recipientID, amount))

	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Transfers <amount> of tokens to the new recipient (clientID of the recipient).
func transferMoney(contract *client.Contract, recipientID string, amount string) {
	fmt.Printf("\n--> Submit Transaction: transferMoney, transfer amount of tokens from the caller to a new recipient \n")

	_, err := contract.Submit("Transfer",
		client.WithArguments(recipientID, amount))

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
