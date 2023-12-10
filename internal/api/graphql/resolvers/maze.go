package resolvers

import (
	"fmt"
	"ftm-explorer/internal/auth"
	"ftm-explorer/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Maze returns the maze.
func (rs *RootResolver) Maze(args struct {
	Address common.Address
}) (*types.Maze, error) {
	if rs.maze == nil {
		return nil, fmt.Errorf("maze is not initialized")
	}
	return rs.maze.GetMaze(args.Address), nil
}

// MazeList returns the list of mazes.
func (rs *RootResolver) MazeList() ([]*types.Maze, error) {
	if rs.maze == nil {
		return nil, fmt.Errorf("maze is not initialized")
	}
	return rs.maze.GetMazeList(), nil
}

// MazeGameSession returns the maze game session.
func (rs *RootResolver) MazeGameSession() (string, error) {
	if rs.maze == nil {
		return "", fmt.Errorf("maze is not initialized")
	}
	return rs.maze.GenerateChallenge()
}

// MazeMyPosition returns the position of the player.
func (rs *RootResolver) MazeMyPosition(args struct {
	Address     common.Address
	Challenge   string
	Signature   string
	MazeAddress common.Address
}) (*types.MazePosition, error) {
	if rs.maze == nil {
		return nil, fmt.Errorf("maze is not initialized")
	}
	// decode the signature
	decodedSignature, err := hexutil.Decode(args.Signature)
	if err != nil {
		return nil, fmt.Errorf("signature hex decoding failed; %s", err)
	}
	// verify signature
	_, err = auth.VerifySignature(args.Challenge, args.Address, decodedSignature)
	if err != nil {
		return nil, fmt.Errorf("signature verification failed; %s", err)
	}
	// return nil if the maze does not exist
	if !rs.maze.Exists(args.MazeAddress) {
		return nil, fmt.Errorf("maze does not exist")
	}
	// get the player position
	tileId, err := rs.repository.MazePlayerPosition(args.MazeAddress, args.Address)
	if err != nil {
		rs.log.Errorf("error getting player position: %v", err)
		return nil, fmt.Errorf("error getting player position")
	}
	// return nil if the player is not on the maze
	if tileId == 0 {
		return nil, nil
	}
	return rs.maze.TileToPosition(args.MazeAddress, int32(tileId))
}
