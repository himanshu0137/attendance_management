#!/bin/bash

# Set up environment variables from .env file
set -a
source .env
set +a

# Run the server
./server