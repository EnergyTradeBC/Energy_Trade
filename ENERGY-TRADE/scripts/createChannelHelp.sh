#!/bin/bash

C_RESET='\033[0m'
C_RED='\033[0;31m'
C_GREEN='\033[0;32m'
C_BLUE='\033[0;34m'
C_YELLOW='\033[1;33m'

# Print the usage message
function createChannelHelp() {
  USAGE="$1"
  if [ "$USAGE" == "genesis" ]; then
    println "Usage: "
    println "  createChannel.sh ${C_GREEN}genesis${C_RESET} [Flags]"
    println
    println "    Flags:"
    println "      -c <channel name> - Name of channel to create (defaults to \"cer\")"
    println "      -verbose - Verbose mode"
    println
    println "      -h - Print this message"
  elif [ "$USAGE" == "create" ]; then
    println "Usage: "
    println "  createChannel.sh ${C_GREEN}create${C_RESET} [Flags]"
    println
    println "    Flags:"
    println "      -c <channel name> - Name of channel to create (defaults to \"cer\")"
    println "      -d <delay> - CLI delays for a certain number of seconds (defaults to 3)"
    println "      -r <max retry> - CLI times out after certain number of attempts (defaults to 5)"
    println "      -verbose - Verbose mode"
    println
    println "      -h - Print this message"
  elif [ "$USAGE" == "updateConfig" ]; then
    println "Usage: "
    println "  createChannel.sh ${C_GREEN}updateConfig <Org>${C_RESET} [Flags]"
    println
    println "    Flags:"
    println "      -c <channel name> - Name of channel to create (defaults to \"cer\")"
    println "      -d <delay> - CLI delays for a certain number of seconds (defaults to 3)"
    println "      -r <max retry> - CLI times out after certain number of attempts (defaults to 5)"
    println "      -verbose - Verbose mode"
    println
    println "      -h - Print this message"
  elif [ "$USAGE" == "signConfig" ]; then
    println "Usage: "
    println "  createChannel.sh ${C_GREEN}signConfig <Org> <as_Org>${C_RESET} [Flags]"
    println
    println "    Flags:"
    println "      -c <channel name> - Name of channel to create (defaults to \"cer\")"
    println "      -d <delay> - CLI delays for a certain number of seconds (defaults to 3)"
    println "      -r <max retry> - CLI times out after certain number of attempts (defaults to 5)"
    println "      -verbose - Verbose mode"
    println
    println "      -h - Print this message"
  elif [ "$USAGE" == "submitConfig" ]; then
    println "Usage: "
    println "  createChannel.sh ${C_GREEN}submitConfig <Org> <as_Org>${C_RESET} [Flags]"
    println
    println "    Flags:"
    println "      -c <channel name> - Name of channel to create (defaults to \"cer\")"
    println "      -d <delay> - CLI delays for a certain number of seconds (defaults to 3)"
    println "      -r <max retry> - CLI times out after certain number of attempts (defaults to 5)"
    println "      -verbose - Verbose mode"
    println
    println "      -h - Print this message"
  elif [ "$USAGE" == "join" ]; then
    println "Usage: "
    println "  createChannel.sh ${C_GREEN}join <Org>${C_RESET} [Flags]"
    println
    println "    Flags:"
    println "      -c <channel name> - Name of channel to create (defaults to \"cer\")"
    println "      -d <delay> - CLI delays for a certain number of seconds (defaults to 3)"
    println "      -r <max retry> - CLI times out after certain number of attempts (defaults to 5)"
    println "      -verbose - Verbose mode"
    println
    println "      -h - Print this message"
  elif [ "$USAGE" == "anchor" ]; then
    println "Usage: "
    println "  createChannel.sh ${C_GREEN}anchor${C_RESET} [Flags]"
    println
    println "    Flags:"
    println "      -c <channel name> - Name of channel to create (defaults to \"cer\")"
    println "      -verbose - Verbose mode"
    println
    println "      -h - Print this message"
  else
    println "Usage: "
    println "  createChannel.sh <Mode> [Flags]"
    println "    Modes:"
    println "      ${C_GREEN}genesis${C_RESET} - Generate the genesis block of the channel"
    println "      ${C_GREEN}create${C_RESET} - Create the application channel"
    println "      ${C_GREEN}updateConfig${C_RESET} - Update the channel configuration to allow a peer to join the channel"
    println "      ${C_GREEN}signConfig${C_RESET} - Sign the channel configuration to allow a peer to join the channel"
    println "      ${C_GREEN}submitConfig${C_RESET} - Submit the channel configuration to allow a peer to join the channel"
    println "      ${C_GREEN}join${C_RESET} - Join a peer to the channel"
    println "      ${C_GREEN}anchor${C_RESET} - Set anchor peer"
    println
    println "    Flags:"
    println "      Used with ${C_GREEN}createChannel.sh genesis${C_RESET}, ${C_GREEN}createChannel.sh updateConfig${C_RESET}, ${C_GREEN}createChannel.sh signConfig${C_RESET}, ${C_GREEN}createChannel.sh submitConfig${C_RESET}, ${C_GREEN}createChannel.sh anchor${C_RESET}:"
    println "        -c <channel name> - Name of channel to create (defaults to \"cer\")"
    println "        -verbose - Verbose mode"
    println
    println "      Used with ${C_GREEN}createChannel.sh create${C_RESET}, ${C_GREEN}createChannel.sh join${C_RESET}:"
    println "        -c <channel name> - Name of channel to create (defaults to \"cer\")"
    println "        -d <delay> - CLI delays for a certain number of seconds (defaults to 3)"
    println "        -r <max retry> - CLI times out after certain number of attempts (defaults to 5)"
    println "        -verbose - Verbose mode"
    println
    println "      -h - Print this message"
  fi
}