function org_files_ccp {
    # it takes file $4 and gives as output the same file after substituting '${ORG}' with $1, '${ADDRESS}' with $2 and '${CHAINCODE_ADDRESS}' with $3
    # (note: as each line of the function substitutes only the first occurence of the keyword in each line of the input file, some are repeated -
    # '${ORG}' appears trice in the same line in some files and '${ADDRESS}' does twice)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${ORG}/$1/" \
        -e "s/\${ORG}/$1/" \
        -e "s/\${ADDRESS}/$2/" \
        -e "s/\${ADDRESS}/$2/" \
        -e "s/\${CHAINCODE_ADDRESS}/$3/" \
        $4 | sed -e $'s/\\\\n/\\\n          /g'
}

function addOrgFiles () {
    # Creates all the files needed for an org to later be inserted into the network.
    # All the files refer to a template version of themselves stored inside of orgTemplate,
    # moreover the folder structure follows the one of orgTemplate. For a generic orgX (the folder orgX is stored inside of ENERGY-TRADE):
    #
    # orgX/
    # |
    # |--compose/
    # |  |
    # |  |--docker/
    # |  |  |
    # |  |  |--peercfg/
    # |  |  |  |
    # |  |  |  |--core.yaml
    # |  |  |
    # |  |  |--docker-compose-orgX.yaml
    # |  |
    # |  |--compose-orgX.yaml
    # |
    # |--configtx.yaml
    # |--orgX-crypto.yaml
    if [ ! -d "org${ORG}" ]; then
        mkdir org${ORG}
    fi
    ADDRESS=$PORT_1
    CHAINCODE_ADDRESS=$PORT_2
    echo "$(org_files_ccp $ORG $ADDRESS $CHAINCODE_ADDRESS orgTemplate/org-crypto.yaml)" > org${ORG}/org${ORG}-crypto.yaml
    echo "$(org_files_ccp $ORG $ADDRESS $CHAINCODE_ADDRESS orgTemplate/configtx.yaml)" > org${ORG}/configtx.yaml
    if [ ! -d "org${ORG}/compose" ]; then
        mkdir org${ORG}/compose
    fi
    echo "$(org_files_ccp $ORG $ADDRESS $CHAINCODE_ADDRESS orgTemplate/compose/compose-org.yaml)" > org${ORG}/compose/compose-org${ORG}.yaml
    if [ ! -d "org${ORG}/compose/docker" ]; then
        mkdir org${ORG}/compose/docker
    fi
    echo "$(org_files_ccp $ORG $ADDRESS $CHAINCODE_ADDRESS orgTemplate/compose/docker/docker-compose-org.yaml)" > org${ORG}/compose/docker/docker-compose-org${ORG}.yaml
    if [ ! -d "org${ORG}/compose/docker/peercfg" ]; then
        mkdir org${ORG}/compose/docker/peercfg
    fi
    cp orgTemplate/compose/docker/peercfg/core.yaml org${ORG}/compose/docker/peercfg
}

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function json_ccp {
    # See org_files_cpp()
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        organizations/ccp-template.json
}

function yaml_ccp {
    # See org_files_cpp()
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        organizations/ccp-template.yaml | sed -e $'s/\\\\n/\\\n          /g'
}

function generateCPP() {
    # Creates the files connection-orgX.json (from template ccp-template.json inside organizations/)
    # and connection-orgX.yaml (from template ccp-template.yaml inside organizations/) inside  organizations/peerOrganizations/orgX.example.com/
    P0PORT=$PORT_1
    CAPORT=$PORT_2
    PEERPEM=organizations/peerOrganizations/org${ORG}.example.com/tlsca/tlsca.org${ORG}.example.com-cert.pem
    CAPEM=organizations/peerOrganizations/org${ORG}.example.com/ca/ca.org${ORG}.example.com-cert.pem

    echo "$(json_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/org${ORG}.example.com/connection-org${ORG}.json
    echo "$(yaml_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/org${ORG}.example.com/connection-org${ORG}.yaml
}

function generateOrg() {    
    # Create Organziation crypto material using cryptogen.
    # The config file (orgX-crypto.yaml) is found in orgX/, the output goes inside organizations/peerOrganizations/.
    # Inside here it creates the folder orgX.example.com and the 5 subfolders ca/, msp/, peers/, tlsca/, users/ and all their contents (including CA certificates)
    export PATH=${PWD}/../bin:$PATH
    which cryptogen
    if [ "$?" -ne 0 ]; then
        fatalln "cryptogen tool not found. exiting"
    fi
    infoln "Generating certificates using cryptogen tool"

    infoln "Creating Org${ORG} Identities"

    set -x
    cryptogen generate --config=org${ORG}/org${ORG}-crypto.yaml --output="organizations"
    res=$?d
    { set +x; } 2>/dev/null
    if [ $res -ne 0 ]; then
        fatalln "Failed to generate certificates..."
    fi

    infoln "Generating CCP files for Org${ORG}"
    generateCPP
}

# Generate channel configuration transaction
function generateOrgDefinition() {
    # Create organization definition (orgX.json) inside organizations/peerOrganizations/orgX.example.com/
    which configtxgen
    if [ "$?" -ne 0 ]; then
        fatalln "configtxgen tool not found. exiting"
    fi
    infoln "Generating Org${ORG} organization definition"
    export FABRIC_CFG_PATH=$PWD/org${ORG}/
    set -x
    configtxgen -printOrg Org${ORG}MSP > organizations/peerOrganizations/org${ORG}.example.com/org${ORG}.json
    res=$?
    { set +x; } 2>/dev/null
    if [ $res -ne 0 ]; then
        fatalln "Failed to generate Org${ORG} organization definition..."
    fi
}

function OrgUp () {
    # Start org nodes creating their docker volume and running them
    if [ "${DATABASE}" == "couchdb" ]; then
        DOCKER_SOCK=${DOCKER_SOCK} ${CONTAINER_CLI_COMPOSE} -f ${COMPOSE_FILE_BASE} -f $COMPOSE_FILE_ORG -f ${COMPOSE_FILE_COUCH_BASE} -f $COMPOSE_FILE_COUCH_ORG up -d 2>&1
    else
        DOCKER_SOCK=${DOCKER_SOCK} ${CONTAINER_CLI_COMPOSE} -f ${COMPOSE_FILE_BASE} -f $COMPOSE_FILE_ORG up -d 2>&1
    fi
    if [ $? -ne 0 ]; then
        fatalln "Unable to start Org network"
    fi
}

function addOrg () {
    # If the test network is not up, abort
    if [ ! -d organizations/ordererOrganizations ]; then
        errorln "Network not up"
    fi

    # generate artifacts if they don't exist
    if [ ! -d "organizations/peerOrganizations/org${ORG}.example.com" ]; then
        generateOrg
        generateOrgDefinition
    fi

    infoln "Bringing up Org${ORG} peer"
    OrgUp
}