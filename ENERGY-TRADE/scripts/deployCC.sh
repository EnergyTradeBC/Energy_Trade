#!/bin/bash

. scripts/utils.sh
. scripts/deployCCHelp.sh
. scripts/envVar.sh
. scripts/actAsOrg.sh

export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/

# Extract the MODE as the first argument and shift to the next one. 
# If there is no argument calls the general help
if [[ $# -lt 1 ]] ; then
    deployCCHelp
    exit 0
else
    MODE=$1
    shift

    if [ "$MODE" = "-h" ]; then
        deployCCHelp
        exit 0

    elif [ "$MODE" = "installChaincode" ]; then
        if [[ $# -lt 1 ]] ; then
            deployCCHelp $MODE
            exit 0
        else
            ORG=$1
            shift
        fi
    
    elif [ "$MODE" = "queryInstalled" ]; then
        if [[ $# -lt 1 ]] ; then
            deployCCHelp $MODE
            exit 0
        else
            ORG=$1
            shift
        fi
    
    elif [ "$MODE" = "approveForMyOrg" ]; then
        if [[ $# -lt 1 ]] ; then
            deployCCHelp $MODE
            exit 0
        else
            ORG=$1
            shift
        fi
    
    elif [ "$MODE" = "checkCommitReadiness" ]; then
        if [[ $# -lt 1 ]] ; then
            deployCCHelp $MODE
            exit 0
        else
            ORG=$1
            shift
        fi
    
    elif [ "$MODE" = "commitChaincodeDefinition" ]; then
        if [[ $# -lt 1 ]] ; then
            deployCCHelp $MODE
            exit 0
        else
            ORG_ARRAY=()
            while [[ $1 =~ ^-?[0-9]+$ ]]; do
              ORG_ARRAY+=($1)
              shift
            done
        fi
    
    elif [ "$MODE" = "queryCommitted" ]; then
        if [[ $# -lt 1 ]] ; then
            deployCCHelp $MODE
            exit 0
        else
            ORG=$1
            shift
        fi
    
    elif [ "$MODE" = "chaincodeInvokeInit" ]; then
        if [[ $# -lt 1 ]] ; then
            deployCCHelp $MODE
            exit 0
        else
            ORG_ARRAY=$1
            shift
        fi
    fi
fi

# Set the default values
CHANNEL_NAME="cer"
CC_NAME=""
CC_SRC_PATH=""
CC_VERSION="1.0"
CC_SEQUENCE="1"
CC_INIT_FCN="NA"
CC_END_POLICY="NA"
CC_COLL_CONFIG="NA"
DELAY="3"
MAX_RETRY="5"
VERBOSE="false"
CC_RUNTIME_LANGUAGE=golang

FABRIC_CFG_PATH=$PWD/../config/

while [[ $# -ge 1 ]] ; do
    key="$1"
    shift
    case $key in
    -h)
        deployCCHelp $MODE
        exit 0
        ;;
    -c)
        if [[ $# -ge 1 ]] ; then
          CHANNEL_NAME="$1"
          shift
        else
          fatalln "No channel name provided after flag '-c'"
        fi
        ;;
    -ccn)
        if [[ $# -ge 1 ]] ; then
          CC_NAME="$1"
          shift
        else
          fatalln "No chaincode name provided after flag '-ccn'"
        fi
        ;;
    -ccp)
        if [[ $# -ge 1 ]] ; then
          CC_SRC_PATH="$1"
          shift
        else
          fatalln "No chaincode source path provided after flag '-ccp'"
        fi
        ;;
    -ccv)
        if [[ $# -ge 1 ]] ; then
          CC_VERSION="$1"
          shift
        else
          fatalln "No chaincode version provided after flag '-ccv'"
        fi
        ;;
    -ccs)
        if [[ $# -ge 1 ]] ; then
          CC_SEQUENCE="$1"
          shift
        else
          fatalln "No chaincode sequence provided after flag '-ccs'"
        fi
        ;;
    -cci)
        if [[ $# -ge 1 ]] ; then
          CC_INIT_FCN="$1"
          shift
        else
          fatalln "No chaincode init function provided after flag '-cci'"
        fi
        ;;
    -ccep)
        if [[ $# -ge 1 ]] ; then
          CC_END_POLICY="$1"
          shift
        else
          fatalln "No chaincode endorsement policy provided after flag '-ccep'"
        fi
        ;;
    -cccg)
        if [[ $# -ge 1 ]] ; then
          CC_COLL_CONFIG="$1"
          shift
        else
          fatalln "No chaincode coll config provided after flag '-cccg'"
        fi
        ;;
    -d)
        if [[ $# -ge 1 ]] ; then
          DELAY="$1"
          shift
        else
          fatalln "No delay provided after flag '-d'"
        fi
        ;;
    -r)
        if [[ $# -ge 1 ]] ; then
          MAX_RETRY="$1"
          shift
        else
          fatalln "No max retry value provided after flag '-d'"
        fi
        ;;
    -verbose)
        VERBOSE=true
        ;;
    *)
        errorln "Unknown flag: $key"
        deployCCHelp
        exit 1
        ;;
    esac
done

#User has not provided a name
if [ -z "$CC_NAME" ] || [ "$CC_NAME" = "NA" ]; then
    fatalln "No chaincode name was provided. Check for the correct parameters in the help"
    deployCCHelp
    exit 1
fi

println "executing with the following"
println "- CHANNEL_NAME: ${C_GREEN}${CHANNEL_NAME}${C_RESET}"
println "- CC_NAME: ${C_GREEN}${CC_NAME}${C_RESET}"
println "- CC_SRC_PATH: ${C_GREEN}${CC_SRC_PATH}${C_RESET}"
println "- CC_VERSION: ${C_GREEN}${CC_VERSION}${C_RESET}"
println "- CC_SEQUENCE: ${C_GREEN}${CC_SEQUENCE}${C_RESET}"
println "- CC_END_POLICY: ${C_GREEN}${CC_END_POLICY}${C_RESET}"
println "- CC_COLL_CONFIG: ${C_GREEN}${CC_COLL_CONFIG}${C_RESET}"
println "- CC_INIT_FCN: ${C_GREEN}${CC_INIT_FCN}${C_RESET}"
println "- DELAY: ${C_GREEN}${DELAY}${C_RESET}"
println "- MAX_RETRY: ${C_GREEN}${MAX_RETRY}${C_RESET}"
println "- VERBOSE: ${C_GREEN}${VERBOSE}${C_RESET}"



INIT_REQUIRED="--init-required"
# check if the init fcn should be called
if [ "$CC_INIT_FCN" = "NA" ]; then
  INIT_REQUIRED=""
fi

if [ "$CC_END_POLICY" = "NA" ]; then
  CC_END_POLICY=""
else
  CC_END_POLICY="--signature-policy $CC_END_POLICY"
fi

if [ "$CC_COLL_CONFIG" = "NA" ]; then
  CC_COLL_CONFIG=""
else
  CC_COLL_CONFIG="--collections-config $CC_COLL_CONFIG"
fi



# import utils
. scripts/envVar.sh
. scripts/deployCCUtils.sh

function checkPrereqs() {
    jq --version > /dev/null 2>&1

    if [[ $? -ne 0 ]]; then
      errorln "jq command not found..."
      errorln
      errorln "Follow the instructions in the Fabric docs to install the prereqs"
      errorln "https://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html"
      exit 1
    fi
}

function vendoringGoDep() {
    infoln "Vendoring Go dependencies at $CC_SRC_PATH"
    pushd $CC_SRC_PATH
    GO111MODULE=on go mod vendor
    popd
    successln "Finished vendoring Go dependencies"
}

# Checking correctness of the mode
if [ "$MODE" == "packageChaincode" ]; then
    # Check if the user has provided a path for the chaincode
    if [ -z "$CC_SRC_PATH" ] || [ "$CC_SRC_PATH" = "NA" ]; then
        fatalln "No chaincode path was provided."
        exit 1
    ## Make sure that the path to the chaincode exists
    elif [ ! -d "$CC_SRC_PATH" ] && [ ! -f "$CC_SRC_PATH" ]; then
        fatalln "Path to chaincode does not exist. Please provide different path."
        exit 1
    fi
    # Vendoring go dependencies at the chaincode source path
    vendoringGoDep
    # Check for prerequisites
    checkPrereqs

    infoln "Packaging chaincode '${CC_NAME}'"
    packageChaincode

elif [ "$MODE" == "installChaincode" ]; then
    infoln "Installing chaincode..."
    installChaincode 

elif [ "$MODE" == "queryInstalled" ]; then
    infoln "Querying whether the chaincode is installed"
    queryInstalled

elif [ "$MODE" == "approveForMyOrg" ]; then
    infoln "Approving the definition of the chaincode"
    approveForMyOrg 

elif [ "$MODE" == "checkCommitReadiness" ]; then
    infoln "Checking whether the chaincode definition is ready to be committed"

    # Dovremmo passargli in modo dinamico la policy strategy per controllare se è verificata o meno? comunque non ho ancora ben capito come 
    # funzioni checkCommitReadiness, suppongo che quello che gli venga passato sia la policy da controllare e a noi dovrebbe interessare 
    # controllare se la policy imposta dal canale sia verificata (?) 
    checkCommitReadiness "\"Org1MSP\": true" "\"Org2MSP\": true" "\"Org3MSP\": true"

elif [ "$MODE" == "commitChaincodeDefinition" ]; then
    infoln "Committing the chaincode definition"
    commitChaincodeDefinition $ORG_ARRAY

elif [ "$MODE" == "queryCommitted" ]; then
    infoln "Checking for successful definition commit"
    queryCommitted

elif [ "$MODE" == "chaincodeInvokeInit" ]; then
    if [ "$CC_INIT_FCN" = "NA" ]; then
        fatalln "No chaincode init function was provided."
    else
        chaincodeInvokeInit $ORG_ARRAY
    fi

elif [ "$MODE" == "chaincodeQuery" ]; then
  true

fi

exit 0
