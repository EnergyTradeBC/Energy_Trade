#!/bin/bash

#!/bin/bash

C_RESET='\033[0m'
C_RED='\033[0;31m'
C_GREEN='\033[0;32m'
C_BLUE='\033[0;34m'
C_YELLOW='\033[1;33m'

# Print the usage message
function printHelp() {
  USAGE="$1"
  if [ "$USAGE" == "up" ]; then
    println "Usage: "
    println "  network.sh ${C_GREEN}up${C_RESET} [Flags]"
    println
    println "    Flags:"
    println "    -ca <use CAs> -  Use Certificate Authorities to generate network crypto material"
    println "    -s <dbtype> - Peer state database to deploy: goleveldb (default) or couchdb"
    println "    -verbose - Verbose mode"
    println
    println "    -h - Print this message"
  elif [ "$USAGE" == "down" ]; then
    println "Usage: "
    println "  network.sh ${C_GREEN}down${C_RESET} [Flags]"
    println
    println "    Flags:"
    println "    -verbose - Verbose mode"
    println
    println "    -h - Print this message"
  elif [ "$USAGE" == "restart" ]; then
    println "Usage: "
    println "  network.sh ${C_GREEN}restart${C_RESET} [Flags]"
    println
    println "    Flags:"
    println "    -ca <use CAs> -  Use Certificate Authorities to generate network crypto material"
    println "    -s <dbtype> - Peer state database to deploy: goleveldb (default) or couchdb"
    println "    -verbose - Verbose mode"
    println
    println "    -h - Print this message"
  else
    println "Usage: "
    println "  network.sh <Mode> [Flags]"
    println "    Modes:"
    println "      ${C_GREEN}up${C_RESET} - Bring up Fabric orderer and peer nodes. No channel is created"
    println "      ${C_GREEN}down${C_RESET} - Bring down the network"
    println "      ${C_GREEN}restart${C_RESET} - Restart the network"
    println
    println "    Flags:"
    println "    Used with ${C_GREEN}network.sh up${C_RESET}, ${C_GREEN}network.sh restart${C_RESET}:"
    println "    -ca <use CAs> -  Use Certificate Authorities to generate network crypto material"
    println "    -s <dbtype> - Peer state database to deploy: goleveldb (default) or couchdb"
    println
    println "    -verbose - Verbose mode"
    println "    -h - Print this message"
  fi
}