package types

import "github.com/ethereum/go-ethereum/common"

// MazePathDirection represents the direction of the maze path.
type MazePathDirection string

const (
	// MazeDirectionEast represents the direction east.
	MazeDirectionEast MazePathDirection = "EAST"

	// MazeDirectionWest represents the direction west.
	MazeDirectionWest MazePathDirection = "WEST"

	// MazeDirectionSouth represents the direction south.
	MazeDirectionSouth MazePathDirection = "SOUTH"

	// MazeDirectionNorth represents the direction north.
	MazeDirectionNorth MazePathDirection = "NORTH"
)

// Maze represents the maze.
type Maze struct {
	// Name is the name of the maze.
	Name string `json:"name"`
	// Address is the address of the maze.
	Address common.Address `json:"address"`
	// Width is the width of the maze.
	Width int32 `json:"width"`
	// Height is the height of the maze.
	Height int32 `json:"height"`
	// VisibilityRange is the visibility range of the maze.
	VisibilityRange int32 `json:"visibilityRange"`
	// StartX is the start x coordinate of the maze.
	StartX int32 `json:"startX"`
	// StartY is the start y coordinate of the maze.
	StartY int32 `json:"startY"`
	// EndX is the end x coordinate of the maze.
	EndX int32 `json:"endX"`
	// EndY is the end y coordinate of the maze.
	EndY int32 `json:"endY"`
}

// MazePath represents the path of the maze.
type MazePath struct {
	// Direction is the direction of the path.
	Direction MazePathDirection `json:"direction"`
	// Length is the visible length of the path.
	Length int32 `json:"length"`
	// Id is the id of the path.
	Id int32 `json:"id"`
	// X is the x coordinate of the path.
	X int32 `json:"x"`
	// Y is the y coordinate of the path.
	Y int32 `json:"y"`
}

// MazePosition represents the position of the maze.
type MazePosition struct {
	// Id is the id of the position.
	Id int32 `json:"id"`
	// X is the x coordinate.
	X int32 `json:"x"`
	// Y is the y coordinate.
	Y int32 `json:"y"`
	// Paths are the paths from the position.
	Paths []MazePath `json:"paths"`
}
