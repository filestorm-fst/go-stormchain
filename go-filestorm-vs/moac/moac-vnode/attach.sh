#!/bin/bash
# attach the testnet1 with default moac ipc
echo "Make sure the moac is running!"
echo "Attaching the console to $HOME./moac/moac.ipc" 
./build/bin/moac attach $HOME/Library/MoacNode/moac.ipc
