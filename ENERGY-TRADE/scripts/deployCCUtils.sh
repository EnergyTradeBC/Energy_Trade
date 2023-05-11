#!/bin/bash

# OGNI FUNZIONE USA IL COMANDO PEER: PRIMA ERA PRESENTE SET GLOBALS IN MODO DA AGIRE COME LO SPECIFICO PEER RICHIESTO,
# MA DAL MOMENTO IN CUI QUESTO CODICE OPERA SU UN SINGOLO PEER I PARAMETRI GLOBALI DOVREBBERO ESSERE GIA' PRESENTI
# COME VARIABILI DI ENVIRONMENT


# Package a chaincode and write the package to a file
function packageChaincode() {
    set -x
    peer lifecycle chaincode package ${CC_NAME}.tar.gz --path ${CC_SRC_PATH} --lang ${CC_RUNTIME_LANGUAGE} --label ${CC_NAME}_${CC_VERSION} >&log.txt
    res=$?
    PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid ${CC_NAME}.tar.gz)
    { set +x; } 2>/dev/null
    cat log.txt
    verifyResult $res "Chaincode packaging has failed"
    successln "Chaincode is packaged"
}

# Install a chaincode 
function installChaincode() {
    # The variables must be created each time since each time we call deployCC we launch/create a new instance of it 
    PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid ${CC_NAME}.tar.gz)

    set -x
    peer lifecycle chaincode queryinstalled --output json | jq -r 'try (.installed_chaincodes[].package_id)' | grep ^${PACKAGE_ID}$ >&log.txt
    if test $? -ne 0; then
        peer lifecycle chaincode install ${CC_NAME}.tar.gz >&log.txt
        res=$?
    fi
    { set +x; } 2>/dev/null
    cat log.txt
    verifyResult $res "Chaincode installation has failed"
    successln "Chaincode is installed"
}

# Verify if a chaincode is installed
# (è praticamente uguale alla precedente però non c'è l'if che permette di installare la chaincode se non è presente)
function queryInstalled() {
    # The variables must be created each time since each time we call deployCC we launch/create a new instance of it 
    PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid ${CC_NAME}.tar.gz)

    set -x
    peer lifecycle chaincode queryinstalled --output json | jq -r 'try (.installed_chaincodes[].package_id)' | grep ^${PACKAGE_ID}$ >&log.txt
    res=$?
    { set +x; } 2>/dev/null
    cat log.txt
    verifyResult $res "Query installed has failed"
    successln "Query installed successful"
}

# Approve a chaincode definition on a channel
# MODIFICARE LA CALL HARDCODED ALL'ORDERER (-o localhost:7050 e --ordererTLSHostnameOverride)
function approveForMyOrg() {
    # The variables must be created each time since each time we call deployCC we launch/create a new instance of it 
    PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid ${CC_NAME}.tar.gz)

    set -x
    peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${CC_VERSION} --package-id ${PACKAGE_ID} --sequence ${CC_SEQUENCE} ${INIT_REQUIRED} ${CC_END_POLICY} ${CC_COLL_CONFIG} >&log.txt
    res=$?
    { set +x; } 2>/dev/null
    cat log.txt
    verifyResult $res "Chaincode definition approved on channel '$CHANNEL_NAME' failed"
    successln "Chaincode definition approved on channel '$CHANNEL_NAME'"
}

# Check whether a chaincode definition is ready to be committed on a channel
function checkCommitReadiness() {
    infoln "Checking the commit readiness of the chaincode definition on channel '$CHANNEL_NAME'..."
    local rc=1
    local COUNTER=1
    # continue to poll
    # we either get a successful response, or reach MAX RETRY
    while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ]; do
        sleep $DELAY
        infoln "Attempting to check the commit readiness of the chaincode definition. Retry after $DELAY seconds."
        set -x
        peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${CC_VERSION} --sequence ${CC_SEQUENCE} ${INIT_REQUIRED} ${CC_END_POLICY} ${CC_COLL_CONFIG} --output json >&log.txt
        res=$?
        { set +x; } 2>/dev/null
        let rc=0
        for var in "$@"; do
            grep "$var" log.txt &>/dev/null || let rc=1
        done
        COUNTER=$(expr $COUNTER + 1)
    done
    cat log.txt
    if test $rc -eq 0; then
        infoln "Checking the commit readiness of the chaincode definition successful on channel '$CHANNEL_NAME'"
    else
        fatalln "After $MAX_RETRY attempts, Check commit readiness result is INVALID!"
    fi
}

# Nei parametri di input della seguente funzione ci sono entrambe le organizzazioni => modificare per fare in modo che possa essere 
# utilizzata da un singolo peer => teoricamente dovrebbe funzionare come le funzioni precedenti, quindi se utilizzata da 1 peer
# per eseguire azioni su se stesso non dovrebbe necessitare dell'argomento PEER_CONN_PARMS, derivato da envVar.sh 
# Per riferimento ai parametri utilizzabili in "peer lifecycle chaincode" guardare https://hyperledger-fabric.readthedocs.io/en/release-2.5/commands/peerlifecycle.html

# Commit the chaincode definition on the channel.
function commitChaincodeDefinition() {
    # ==> MODIFICARE LA CALL HARDCODED ALL'ORDERER (-o localhost:7050) E CAPIRE COSA è --ordererTLSHostnameOverride

    # while 'peer chaincode' command can get the orderer endpoint from the
    # peer (if join was successful), let's supply it directly as we know
    # it using the "-o" option 
    set -x
    peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${CC_VERSION} --sequence ${CC_SEQUENCE} ${INIT_REQUIRED} ${CC_END_POLICY} ${CC_COLL_CONFIG} >&log.txt
    res=$?
    { set +x; } 2>/dev/null
    cat log.txt
    verifyResult $res "Chaincode definition commit failed on peer0.org${ORG} on channel '$CHANNEL_NAME' failed"
    successln "Chaincode definition committed on channel '$CHANNEL_NAME'"
}

# queryCommitted ORG
function queryCommitted() {
    EXPECTED_RESULT="Version: ${CC_VERSION}, Sequence: ${CC_SEQUENCE}, Endorsement Plugin: escc, Validation Plugin: vscc"
    infoln "Querying chaincode definition on channel '$CHANNEL_NAME'..."
    local rc=1
    local COUNTER=1
    # continue to poll
    # we either get a successful response, or reach MAX RETRY
    while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ]; do
        sleep $DELAY
        infoln "Attempting to Query committed status. Retry after $DELAY seconds."
        set -x
        peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name ${CC_NAME} >&log.txt
        res=$?
        { set +x; } 2>/dev/null
        test $res -eq 0 && VALUE=$(cat log.txt | grep -o '^Version: '$CC_VERSION', Sequence: [0-9]*, Endorsement Plugin: escc, Validation Plugin: vscc')
        test "$VALUE" = "$EXPECTED_RESULT" && let rc=0
        COUNTER=$(expr $COUNTER + 1)
    done
    cat log.txt
    if test $rc -eq 0; then
        successln "Query chaincode definition successful on channel '$CHANNEL_NAME'"
    else
        fatalln "After $MAX_RETRY attempts, Query chaincode definition result is INVALID!"
    fi
}

# Fare riferimento a https://hyperledger-fabric.readthedocs.io/en/release-2.5/commands/peerchaincode.html#peer-chaincode-invoke
# per capire il funzionamento di --isInit (per quanto riguarda invece l'utilizzo della funzione "chaincode invoke" per un solo peer
# invece che per una lista di peers vale lo stesso di commitChaincodeDefinition)

# The chaincode must have the initLedger method defined

function chaincodeInvokeInit() {
    # ==> MODIFICARE LA CALL HARDCODED ALL'ORDERER (-o localhost:7050) E CAPIRE COSA è --ordererTLSHostnameOverride

    # while 'peer chaincode' command can get the orderer endpoint from the
    # peer (if join was successful), let's supply it directly as we know
    # it using the "-o" option
    set -x
    fcn_call='{"function":"'${CC_INIT_FCN}'","Args":[]}'
    infoln "invoke fcn call:${fcn_call}"
    peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" -C $CHANNEL_NAME -n ${CC_NAME} "${PEER_CONN_PARMS[@]}" --isInit -c ${fcn_call} >&log.txt
    res=$?
    { set +x; } 2>/dev/null
    cat log.txt
    verifyResult $res "Invoke execution failed "
    successln "Invoke transaction successful on channel '$CHANNEL_NAME'"
}





# Questa funzione esegue una query alla chaincode; un esempio di argomento (Args) è '{"Args":["ReadAsset","asset6"]}' che permette di 
# leggere l'asset numero 6 dalla chaincode => è necessario inserire una funzione del genere all'interno dello script per eseguire il 
# deploy della chaincode????

# ORG è sottintesa essere una variabile di environment presettata e disponibile su ogni peer

function chaincodeQuery() {
    infoln "Querying on channel '$CHANNEL_NAME'..."
    local rc=1
    local COUNTER=1
    # continue to poll
    # we either get a successful response, or reach MAX RETRY
    while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ]; do
        sleep $DELAY
        infoln "Attempting to Query. Retry after $DELAY seconds."
        set -x
        peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"Args":["org.hyperledger.fabric:GetMetadata"]}' >&log.txt
        res=$?
        { set +x; } 2>/dev/null
        let rc=$res
        COUNTER=$(expr $COUNTER + 1)
    done
    cat log.txt
    if test $rc -eq 0; then
        successln "Query successful on peer0.org${ORG} on channel '$CHANNEL_NAME'"
    else
        fatalln "After $MAX_RETRY attempts, Query result on peer0.org${ORG} is INVALID!"
    fi
}