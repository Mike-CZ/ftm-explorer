package maze

import (
	"ftm-explorer/internal/types"

	"github.com/ethereum/go-ethereum/common"
)

//go:generate mockgen -source=interface.go -destination=maze_mock.go -package=maze -mock_names=IMaze=MockMaze

// IMaze represents a maze.
type IMaze interface {
	// Exists returns true if the maze exists on given address.
	Exists(common.Address) bool
	// GetMaze returns the maze.
	GetMaze(common.Address) *types.Maze
	// GetMazeList returns the list of mazes.
	GetMazeList() []*types.Maze
	// TileToPosition converts the tile id to position.
	TileToPosition(common.Address, int32) (*types.MazePosition, error)
	// GenerateChallenge generates a challenge for the maze.
	GenerateChallenge() (string, error)
}
