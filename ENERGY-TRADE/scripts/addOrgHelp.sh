#!/bin/bash

C_RESET='\033[0m'
C_RED='\033[0;31m'
C_GREEN='\033[0;32m'
C_BLUE='\033[0;34m'
C_YELLOW='\033[1;33m'

# Print the usage message
function addOrgHelp() {
  USAGE="$1"
  if [ "$USAGE" == "files" ]; then
    println "Usage: "
    println "  addOrg.sh ${C_GREEN}files <Org> <Port 1> <Port 2>${C_RESET} [Flags]"
    println
    println "    Flags:"
    println
    println "      -h - Print this message"
  elif [ "$USAGE" == "add" ]; then
    println "Usage: "
    println "  addOrg.sh ${C_GREEN}add <Org> <Port 1> <Port 2>${C_RESET} [Flags]"
    println
    println "    Flags:"
    println
    println "      -h - Print this message"
  else
    println "Usage: "
    println "  addOrg.sh <Mode> <Org> <Port 1> <Port 2> [Flags]"
    println "    Modes:"
    println "      ${C_GREEN}files${C_RESET} - Create all the configuration files for the org"
    println "      ${C_GREEN}add${C_RESET} - Add the org to the network"
    println
    println "    Flags:"
    println
    println "      -h - Print this message"
  fi
}