package maze

import (
	"fmt"
	"ftm-explorer/internal/config"
	"ftm-explorer/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tyler-smith/go-bip39"
)

// kMazeChallengePrefix is prefix of message to be signed using Metamask.
const kMazeChallengePrefix = "Please sign following text to obtain your maze position:\n\n"

// mazeTile represents the path of the maze.
type mazeTile struct {
	X int32
	Y int32
	// Directions to other tiles, if nil, there is no path
	// Holds the id of the tile
	North *int32
	South *int32
	East  *int32
	West  *int32
}

// Maze represents a maze.
type Maze struct {
	cfg *config.Maze
	// address => maze
	mazes map[string]*types.Maze
	// address => tile id => maze tile
	tiles map[string]map[int32]mazeTile
}

// NewMaze creates a new maze.
func NewMaze(cfg *config.Maze) *Maze {
	m := &Maze{
		cfg: cfg,
	}
	// initialize mazes
	m.mazes = make(map[string]*types.Maze)
	m.tiles = make(map[string]map[int32]mazeTile)
	for _, mcfg := range cfg.Configs {
		m.mazes[mcfg.Address] = &types.Maze{
			Name:            mcfg.Name,
			Address:         common.HexToAddress(mcfg.Address),
			Width:           mcfg.Width,
			Height:          mcfg.Height,
			VisibilityRange: int32(cfg.VisibilityRange),
			StartX:          mcfg.Tiles[mcfg.Entry].Position.X,
			StartY:          mcfg.Tiles[mcfg.Entry].Position.Y,
			EndX:            mcfg.Tiles[mcfg.Exit].Position.X,
			EndY:            mcfg.Tiles[mcfg.Exit].Position.Y,
		}

		// initialize tiles
		m.tiles[mcfg.Address] = make(map[int32]mazeTile)
		for _, tile := range mcfg.Tiles {
			path := mazeTile{
				X: tile.Position.X,
				Y: tile.Position.Y,
			}
			// initialize paths, config holds indexes to array, we need to convert them to ids
			if tile.Paths.North != nil {
				path.North = &mcfg.Tiles[*tile.Paths.North].Id
			}
			if tile.Paths.South != nil {
				path.South = &mcfg.Tiles[*tile.Paths.South].Id
			}
			if tile.Paths.East != nil {
				path.East = &mcfg.Tiles[*tile.Paths.East].Id
			}
			if tile.Paths.West != nil {
				path.West = &mcfg.Tiles[*tile.Paths.West].Id
			}
			m.tiles[mcfg.Address][tile.Id] = path
		}
	}
	return m
}

// Exists returns true if the maze exists on given address.
func (m *Maze) Exists(address common.Address) bool {
	_, ok := m.mazes[address.Hex()]
	return ok
}

// GetMaze returns the maze.
func (m *Maze) GetMaze(address common.Address) *types.Maze {
	maze, ok := m.mazes[address.Hex()]
	if !ok {
		return nil
	}
	return maze
}

// GetMazeList returns the list of mazes.
func (m *Maze) GetMazeList() []*types.Maze {
	mazes := make([]*types.Maze, 0)
	for _, maze := range m.mazes {
		mazes = append(mazes, maze)
	}
	return mazes
}

// TileToPosition converts the tile id to position.
func (m *Maze) TileToPosition(mazeAddress common.Address, tileId int32) (*types.MazePosition, error) {
	tiles, ok := m.tiles[mazeAddress.Hex()]
	if !ok {
		return nil, fmt.Errorf("maze not found: %s", mazeAddress)
	}
	tile, ok := tiles[tileId]
	if !ok {
		return nil, fmt.Errorf("tile not found: %d", tileId)
	}

	// initialize all paths from the tile
	paths := make([]types.MazePath, 0)
	if tile.North != nil {
		next := tiles[*tile.North]
		northPath := types.MazePath{
			Direction: types.MazeDirectionNorth,
			Id:        *tile.North,
			X:         next.X,
			Y:         next.Y,
			Length:    1,
		}
		// calculate length
		for {
			if next.North == nil || uint(northPath.Length) >= m.cfg.VisibilityRange {
				break
			}
			next = tiles[*next.North]
			northPath.Length++
		}
		paths = append(paths, northPath)
	}
	if tile.South != nil {
		next := tiles[*tile.South]
		southPath := types.MazePath{
			Direction: types.MazeDirectionSouth,
			Id:        *tile.South,
			X:         next.X,
			Y:         next.Y,
			Length:    1,
		}
		// calculate length
		for {
			if next.South == nil || uint(southPath.Length) >= m.cfg.VisibilityRange {
				break
			}
			next = tiles[*next.South]
			southPath.Length++
		}
		paths = append(paths, southPath)
	}
	if tile.East != nil {
		next := tiles[*tile.East]
		eastPath := types.MazePath{
			Direction: types.MazeDirectionEast,
			Id:        *tile.East,
			X:         next.X,
			Y:         next.Y,
			Length:    1,
		}
		// calculate length
		for {
			if next.East == nil || uint(eastPath.Length) >= m.cfg.VisibilityRange {
				break
			}
			next = tiles[*next.East]
			eastPath.Length++
		}
		paths = append(paths, eastPath)
	}
	if tile.West != nil {
		next := tiles[*tile.West]
		westPath := types.MazePath{
			Direction: types.MazeDirectionWest,
			Id:        *tile.West,
			X:         next.X,
			Y:         next.Y,
			Length:    1,
		}
		// calculate length
		for {
			if next.West == nil || uint(westPath.Length) >= m.cfg.VisibilityRange {
				break
			}
			next = tiles[*next.West]
			westPath.Length++
		}
		paths = append(paths, westPath)
	}

	return &types.MazePosition{
		Id:    tileId,
		X:     tile.X,
		Y:     tile.Y,
		Paths: paths,
	}, nil
}

// GenerateChallenge generates a challenge for the maze.
func (m *Maze) GenerateChallenge() (string, error) {
	phrase, err := m.generatePhrase()
	if err != nil {
		return "", fmt.Errorf("error generating phrase: %v", err)
	}
	return kMazeChallengePrefix + phrase, nil
}

// generatePhrase generates a phrase for the maze.
func (m *Maze) generatePhrase() (string, error) {
	// generate phrase based on bip-39 standard
	entropy, err := bip39.NewEntropy(256) // 256 bits to get a 24-word mnemonic
	if err != nil {
		return "", fmt.Errorf("error generating entropy: %v", err)
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("error generating mnemonic: %v", err)
	}
	// calculate hash from mnemonic
	hash, err := bip39.MnemonicToByteArray(mnemonic)
	if err != nil {
		return "", fmt.Errorf("error calculating hash from mnemonic: %v", err)
	}
	// return hex encoded hash
	return hexutil.Encode(hash), nil
}
