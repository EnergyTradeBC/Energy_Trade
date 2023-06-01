. scripts/utils.sh

export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem

actAsOrg() {
  local USING_ORG=""
  if [ -z "$OVERRIDE_ORG" ]; then
    USING_ORG=$1
  else
    USING_ORG="${OVERRIDE_ORG}"
  fi
  infoln "Using organization ${USING_ORG}"
  if [ $USING_ORG -eq 1 ]; then
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
  elif [ $USING_ORG -eq 2 ]; then
    export CORE_PEER_LOCALMSPID="Org2MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/tlsca/tlsca.org2.example.com-cert.pem
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
    export CORE_PEER_ADDRESS=localhost:9051
  else
    export CORE_PEER_LOCALMSPID="Org${USING_ORG}MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org${USING_ORG}.example.com/tlsca/tlsca.org${USING_ORG}.example.com-cert.pem
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org${USING_ORG}.example.com/users/Admin@org${USING_ORG}.example.com/msp
    export CORE_PEER_ADDRESS=localhost:$((${USING_ORG}+8))051
  fi

  if [ "$VERBOSE" == "true" ]; then
    env | grep CORE
  fi
}

# Set environment variables for use in the CLI container
actAsOrgCLI() {
    actAsOrg $1

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
    else
        export CORE_PEER_ADDRESS=peer0.org${ORG}.example.com:$((${USING_ORG}+8))051
  fi
}

# parsePeerConnectionParameters $@
# Helper function that sets the peer connection parameters for a chaincode
# operation
parsePeerConnectionParameters() {
  PEER_CONN_PARMS=()
  PEERS=""
  ORG_ARRAY=$1
  for org in ${ORG_ARRAY[@]}; do
    actAsOrg $org
    PEER="peer0.org$org"
    ## Set peer addresses
    if [ -z "$PEERS" ]
    then
	PEERS="$PEER"
    else
	PEERS="$PEERS $PEER"
    fi
    PEER_CONN_PARMS=("${PEER_CONN_PARMS[@]}" --peerAddresses $CORE_PEER_ADDRESS)
  
    TLSINFO=(--tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE)
    PEER_CONN_PARMS=("${PEER_CONN_PARMS[@]}" "${TLSINFO[@]}")
  done
}

verifyResult() {
  if [ $1 -ne 0 ]; then
    fatalln "$2"
  fi
}