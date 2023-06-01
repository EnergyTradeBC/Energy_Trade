/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"asset-transfer-energy/chaincode"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	assetChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating asset-transfer-energy chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting asset-transfer-energy chaincode: %v", err)
	}
}
