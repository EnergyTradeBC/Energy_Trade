#!/bin/bash

C_RESET='\033[0m'
C_RED='\033[0;31m'
C_GREEN='\033[0;32m'
C_BLUE='\033[0;34m'
C_YELLOW='\033[1;33m'

# Print the usage message
function deployCCHelp() {
    USAGE="$1"
    if [ "$USAGE" == "packageChaincode" ]; then
        println "Usage: "
        println "  deployCC.sh ${C_GREEN}packageChaincode${C_RESET} [Flags]"
        println
        println "    Flags:"
        println "      -ccn <name> - Chaincode name"
        println "      -ccp <path> - File path to the chaincode"
        println "      -ccv <version> - Chaincode version (default to \"1.0\")"
        println "      -verbose - Verbose mode"
        println
        println "      -h - Print this message"
    elif [ "$USAGE" == "installChaincode" ]; then
        println "Usage: "
        println "  deployCC.sh ${C_GREEN}installChaincode${C_RESET} [Flags]"
        println
        println "    Flags:"
        println "      -ccn <name> - Chaincode name"
        println "      -verbose - Verbose mode"
        println
        println "      -h - Print this message"
    elif [ "$USAGE" == "queryInstalled" ]; then
        println "Usage: "
        println "  deployCC.sh ${C_GREEN}queryInstalled${C_RESET} [Flags]"
        println
        println "    Flags:"
        println "      -ccn <name> - Chaincode name"
        println "      -verbose - Verbose mode"
        println
        println "      -h - Print this message"
    elif [ "$USAGE" == "approveForMyOrg" ]; then
        println "Usage: "
        println "  deployCC.sh ${C_GREEN}approveForMyOrg${C_RESET} [Flags]"
        println
        println "    Flags:"
        println "      -c <channel name> - Name of the channel in which the chaincode will be approved (default to \"CER\")"
        println "      -ccn <name> - Chaincode name"
        println "      -ccv <version> - Chaincode version (default to \"1.0\")"
        println "      -ccs <sequence> - Chaincode definition sequence. Must be an integer (default to \"1\")"
        println "      -cci <fcn name> - Name of chaincode initialization function (default to \"NA\")"
        println "      -ccep <policy> - Chaincode endorsement policy using signature policy syntax (default to \"NA\")"
        println "      -cccg <collection-config> - File path to private data collections configuration file (default to \"NA\")"
        println "      -verbose - Verbose mode"
        println
        println "      -h - Print this message"
    elif [ "$USAGE" == "checkCommitReadiness" ]; then
        println "Usage: "
        println "  deployCC.sh ${C_GREEN}checkCommitReadiness${C_RESET} [Flags]"
        println
        println "    Flags:"
        println "      -c <channel name> - Name of the channel in which the chaincode definition will be checked (default to \"CER\")"
        println "      -d <delay> - CLI delays for a certain number of seconds (defaults to 3)"
        println "      -r <max retry> - CLI times out after certain number of attempts (defaults to 5)"
        println "      -ccn <name> - Chaincode name"
        println "      -ccv <version> - Chaincode version (default to \"1.0\")"
        println "      -ccs <sequence> - Chaincode definition sequence. Must be an integer (default to \"1\")"
        println "      -cci <fcn name> - Name of chaincode initialization function (default to \"NA\")"
        println "      -ccep <policy> - Chaincode endorsement policy using signature policy syntax (default to \"NA\")"
        println "      -cccg <collection-config> - File path to private data collections configuration file (default to \"NA\")"
        println "      -verbose - Verbose mode"
        println
        println "      -h - Print this message"
    elif [ "$USAGE" == "commitChaincodeDefinition" ]; then
        println "Usage: "
        println "  deployCC.sh ${C_GREEN}commitChaincodeDefinition${C_RESET} [Flags]"
        println
        println "    Flags:"
        println "      -c <channel name> - Name of the channel in which the chaincode definition will be committed (default to \"CER\")"
        println "      -ccn <name> - Chaincode name"
        println "      -ccv <version> - Chaincode version (default to \"1.0\")"
        println "      -ccs <sequence> - Chaincode definition sequence. Must be an integer (default to \"1\")"
        println "      -cci <fcn name> - Name of chaincode initialization function (default to \"NA\")"
        println "      -ccep <policy> - Chaincode endorsement policy using signature policy syntax (default to \"NA\")"
        println "      -cccg <collection-config> - File path to private data collections configuration file (default to \"NA\")"
        println "      -verbose - Verbose mode"
        println
        println "      -h - Print this message"
    elif [ "$USAGE" == "queryCommitted" ]; then
        println "Usage: "
        println "  deployCC.sh ${C_GREEN}queryCommitted${C_RESET} [Flags]"
        println
        println "    Flags:"
        println "      -c <channel name> - Name of the channel in which will be checked the commit of the chaincode (default to \"CER\")"
        println "      -d <delay> - CLI delays for a certain number of seconds (defaults to 3)"
        println "      -r <max retry> - CLI times out after certain number of attempts (defaults to 5)"
        println "      -ccn <name> - Chaincode name"
        println "      -ccv <version> - Chaincode version (default to \"1.0\")"
        println "      -ccs <sequence> - Chaincode definition sequence. Must be an integer (default to \"1\")"
        println "      -verbose - Verbose mode"
        println
        println "      -h - Print this message"
    elif [ "$USAGE" == "chaincodeInvokeInit" ]; then
        println "Usage: "
        println "  deployCC.sh ${C_GREEN}chaincodeInvokeInit${C_RESET} [Flags]"
        println
        println "    Flags:"
        println "      -c <channel name> - Name of the channel in which the chaincode will be invoked (default to \"CER\")"
        println "      -ccn <name> - Chaincode name"
        println "      -cci <fcn name> - Name of chaincode initialization function (default to \"NA\")"
        println "      -verbose - Verbose mode"
        println
        println "      -h - Print this message"

    # elif [ "$USAGE" == "chaincodeQuery" ]; then
    #     println "Usage: "
    #     println "  deployCC.sh ${C_GREEN}chaincodeQuery${C_RESET} [Flags]"
    #     println
    #     println "    Flags:"
    #     println "      -c <channel name> - Name of the channel in which the chaincode will be invoked (default to \"CER\")"
    #     println "      -ccn <name> - Chaincode name"
    #     println "      -d <delay> - CLI delays for a certain number of seconds (defaults to 3)"
    #     println "      -r <max retry> - CLI times out after certain number of attempts (defaults to 5)"
    #     println "      -verbose - Verbose mode"
    #     println
    #     println "      -h - Print this message"

    else
        println "Usage: "
        println "  deployCC.sh <Mode> [Flags]"
        println "    Modes:"
        println "      ${C_GREEN}packageChaincode${C_RESET} - Generate the genesis block of the channel"
        println "      ${C_GREEN}installChaincode${C_RESET} - Create the application channel"
        println "      ${C_GREEN}queryInstalled${C_RESET} - Join a peer to the channel"
        println "      ${C_GREEN}approveForMyOrg${C_RESET} - Set anchor peer"
        println "      ${C_GREEN}checkCommitReadiness${C_RESET} - Generate the genesis block of the channel"
        println "      ${C_GREEN}commitChaincodeDefinition${C_RESET} - Create the application channel"
        println "      ${C_GREEN}queryCommitted${C_RESET} - Join a peer to the channel"
        println "      ${C_GREEN}ChaincodeInvokeInit${C_RESET} - Set anchor peer"
        println "      ${C_GREEN}ChaincodeInvokeInit${C_RESET} - Set anchor peer"
        println
        println "    Flags:"
        println "      -c <channel name> - Name of the channel in which the chaincode definition will be checked (default to \"CER\")"
        println "      -d <delay> - CLI delays for a certain number of seconds (defaults to \"3\")"
        println "      -r <max retry> - CLI times out after certain number of attempts (defaults to \"5\")"
        println "      -ccn <name> - Chaincode name"
        println "      -ccp <path> - File path to the chaincode"
        println "      -ccv <version> - Chaincode version (default to \"1.0\")"
        println "      -ccs <sequence> - Chaincode definition sequence. Must be an integer (default to \"1\")"
        println "      -cci <fcn name> - Name of chaincode initialization function (default to \"NA\")"
        println "      -ccep <policy> - Chaincode endorsement policy using signature policy syntax (default to \"NA\")"
        println "      -cccg <collection-config> - File path to private data collections configuration file (default to \"NA\")"
        println "      -verbose - Verbose mode"
        println
        println "      -h - Print this message"
    fi
}