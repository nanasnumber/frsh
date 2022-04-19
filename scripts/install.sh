#!/bin/bash

# Script to install program
PROGRAM_NAME="frsh"

BIN_PATH="/usr/local/bin"
USER_HOME=$HOME;
PROGRAM_DIR=".${PROGRAM_NAME}";

echo "Installing ${PROGRAM_NAME} to $BIN_PATH...";

sudo cp ./dist/${PROGRAM_NAME} $BIN_PATH;

echo "${PROGRAM_NAME} installed!"


