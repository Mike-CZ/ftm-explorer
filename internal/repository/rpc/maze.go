package rpc

import (
	"context"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum"
	abi2 "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// MazePlayerPosition returns the player's position in the maze.
func (rpc *OperaRpc) MazePlayerPosition(ctx context.Context, mazeAddr common.Address, playerAddr common.Address) (uint16, error) {
	// contract abi definition to get the player's position in the maze
	definition := `[{
    "inputs": [
      {
        "internalType": "address",
        "name": "player",
        "type": "address"
      }
    ],
    "name": "onTile",
    "outputs": [
      {
        "internalType": "uint16",
        "name": "",
        "type": "uint16"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  }]`
	abi, err := abi2.JSON(strings.NewReader(definition))
	if err != nil {
		return 0, err
	}
	// pack constructor params
	data, err := abi.Pack("onTile", playerAddr)
	if err != nil {
		return 0, err
	}
	msg := ethereum.CallMsg{
		To:   &mazeAddr,
		Data: data,
	}
	output, err := ethclient.NewClient(rpc.ftm).CallContract(ctx, msg, nil)
	if err != nil {
		return 0, err
	}
	tile, err := abi.Unpack("onTile", output)
	if err != nil {
		log.Fatalf("Failed to unpack output: %v", err)
	}
	return tile[0].(uint16), nil
}
