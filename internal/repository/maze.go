package repository

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

// MazePlayerPosition returns the position of the player in the maze.
func (r *Repository) MazePlayerPosition(mazeAddr common.Address, playerAddr common.Address) (uint16, error) {
	ctx, cancel := context.WithTimeout(context.Background(), kRpcTimeout)
	defer cancel()
	return r.rpc.MazePlayerPosition(ctx, mazeAddr, playerAddr)
}
