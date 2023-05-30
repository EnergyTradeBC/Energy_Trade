#!/bin/bash

# imports
. scripts/utils.sh
# some of the variables exported in it are used inside this sript
. scripts/setEnvVar.sh

# Set environment variables for the peer org
setGlobals() {
    local USING_ORG=""
    if [ -z "$OVERRIDE_ORG" ]; then
        USING_ORG=$1
    else
        USING_ORG="${OVERRIDE_ORG}"
    fi
    infoln "Using organization ${USING_ORG}"
    if [ $USING_ORG -eq 1 ]; then
        export CORE_PEER_LOCALMSPID="Org1MSP"
        export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
        export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
        export CORE_PEER_ADDRESS=localhost:7051

    elif [ $USING_ORG -eq 2 ]; then
        export CORE_PEER_LOCALMSPID="Org2MSP"
        export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG2_CA
        export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
        export CORE_PEER_ADDRESS=localhost:9051

    elif [ $USING_ORG -eq 3 ]; then
        export CORE_PEER_LOCALMSPID="Org3MSP"
        export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG3_CA
        export CORE_PEER_MSPCONFIGPATH=${PWD}/o11051rganizations/peerOrganizations/org3.example.com/users/Admin@org3.example.com/msp
        export CORE_PEER_ADDRESS=localhost:

    elif [ $USING_ORG -eq 9 ]; then
        export CORE_PEER_LOCALMSPID="Org9MSP"
        export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG3_CA
        export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org9.example.com/users/Admin@org9.example.com/msp
        export CORE_PEER_ADDRESS=localhost:
    else
        errorln "ORG Unknown"
    fi

    if [ "$VERBOSE" == "true" ]; then
        env | grep CORE
    fi
}

# Set environment variables for use in the CLI container
setGlobalsCLI() {
  setGlobals $1

  local USING_ORG=""
  if [ -z "$OVERRIDE_ORG" ]; then
    USING_ORG=$1
  else
    USING_ORG="${OVERRIDE_ORG}"
  fi
  if [ $USING_ORG -eq 1 ]; then
    export CORE_PEER_ADDRESS=peer0.org1.example.com:7051
  elif [ $USING_ORG -eq 2 ]; then
    export CORE_PEER_ADDRESS=peer0.org2.example.com:9051
  elif [ $USING_ORG -eq 3 ]; then
    export CORE_PEER_ADDRESS=peer0.org3.example.com:11051
  else
    errorln "ORG Unknown"
  fi
}
