package maze

import (
	"encoding/json"
	"ftm-explorer/internal/config"
	"ftm-explorer/internal/types"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// Test that the phrase generator returns a valid phrase.
func TestMaze_GeneratePhrase(t *testing.T) {
	m := createMaze(t)
	phrase, err := m.generatePhrase()
	if err != nil {
		t.Fatalf("GeneratePhrase failed: %v", err)
	}
	if len(phrase) == 0 {
		t.Fatalf("Invalid phrase returned: %s", phrase)
	}
}

// Test that the challenge generator returns a valid challenge.
func TestMaze_GenerateChallenge(t *testing.T) {
	m := createMaze(t)
	challenge, err := m.GenerateChallenge()
	if err != nil {
		t.Fatalf("GenerateChallenge failed: %v", err)
	}
	if len(challenge) == 0 {
		t.Fatalf("Invalid challenge returned: %s", challenge)
	}
}

// Test that the maze exists.
func TestMaze_GetMaze(t *testing.T) {
	m := createMaze(t)
	maze := m.GetMaze(common.HexToAddress("0x6dDb82e5B91e5941f4633Cb3ee49560Fb4582EFA"))
	if maze == nil {
		t.Fatalf("GetMaze returned nil")
	}
}

// Test that the maze does not exist.
func TestMaze_GetMaze_NotFound(t *testing.T) {
	m := createMaze(t)
	maze := m.GetMaze(common.HexToAddress("0x6dDb82e5B91e5941f4633Cb3ee49560Fb4582EFB"))
	if maze != nil {
		t.Fatalf("GetMaze returned maze")
	}
}

// Test that the first tile is returned correctly.
func TestMaze_First_Tile(t *testing.T) {
	m := createMaze(t)
	tile, err := m.TileToPosition(common.HexToAddress("0x6dDb82e5B91e5941f4633Cb3ee49560Fb4582EFA"), 3145)
	if err != nil {
		t.Fatalf("TileToPosition failed: %v", err)
	}
	if tile.X != 0 || tile.Y != 0 {
		t.Fatalf("TileToPosition returned invalid tile: %v", tile)
	}
	if tile.Id != 3145 {
		t.Fatalf("TileToPosition returned invalid tile id: %v", tile)
	}
	if len(tile.Paths) != 1 {
		t.Fatalf("TileToPosition returned invalid paths: %v", tile)
	}
	if tile.Paths[0].Direction != types.MazeDirectionSouth {
		t.Fatalf("TileToPosition returned invalid path: %v", tile)
	}
	if tile.Paths[0].Length != 2 {
		t.Fatalf("TileToPosition returned invalid path length: %v", tile)
	}
	if tile.Paths[0].Id != 1809 {
		t.Fatalf("TileToPosition returned invalid path id: %v", tile)
	}
	if tile.Paths[0].X != 0 || tile.Paths[0].Y != 1 {
		t.Fatalf("TileToPosition returned invalid path coordinates: %v", tile)
	}
}

// Test that the exit tile is returned correctly.
func TestMaze_Exit_Tile(t *testing.T) {
	m := createMaze(t)
	tile, err := m.TileToPosition(common.HexToAddress("0x6dDb82e5B91e5941f4633Cb3ee49560Fb4582EFA"), 2145)
	if err != nil {
		t.Fatalf("TileToPosition failed: %v", err)
	}
	if tile.X != 0 || tile.Y != 3 {
		t.Fatalf("TileToPosition returned invalid tile: %v", tile)
	}
	if len(tile.Paths) != 1 {
		t.Fatalf("TileToPosition returned invalid paths: %v", tile)
	}
	if tile.Paths[0].Direction != types.MazeDirectionSouth {
		t.Fatalf("TileToPosition returned invalid path: %v", tile)
	}
	if tile.Paths[0].Length != 3 {
		t.Fatalf("TileToPosition returned invalid path length: %v", tile)
	}
	if tile.Paths[0].Id != 3010 {
		t.Fatalf("TileToPosition returned invalid path id: %v", tile)
	}
	if tile.Paths[0].X != 0 || tile.Paths[0].Y != 4 {
		t.Fatalf("TileToPosition returned invalid path coordinates: %v", tile)
	}
}

// Test that the tile with many paths is returned correctly.
func TestMaze_Tile_With_Many_Paths(t *testing.T) {
	m := createMaze(t)
	tile, err := m.TileToPosition(common.HexToAddress("0x6dDb82e5B91e5941f4633Cb3ee49560Fb4582EFA"), 882)
	if err != nil {
		t.Fatalf("TileToPosition failed: %v", err)
	}
	if tile.X != 4 || tile.Y != 7 {
		t.Fatalf("TileToPosition returned invalid tile: %v", tile)
	}
	if len(tile.Paths) != 3 {
		t.Fatalf("TileToPosition returned invalid paths: %v", tile)
	}
	if tile.Paths[0].Direction != types.MazeDirectionNorth {
		t.Fatalf("TileToPosition returned invalid path: %v", tile)
	}
	if tile.Paths[0].Length != 1 {
		t.Fatalf("TileToPosition returned invalid path length: %v", tile)
	}
	if tile.Paths[0].Id != 396 {
		t.Fatalf("TileToPosition returned invalid path id: %v", tile)
	}
	if tile.Paths[0].X != 4 || tile.Paths[0].Y != 6 {
		t.Fatalf("TileToPosition returned invalid path coordinates: %v", tile)
	}
	if tile.Paths[1].Direction != types.MazeDirectionEast {
		t.Fatalf("TileToPosition returned invalid path: %v", tile)
	}
	if tile.Paths[1].Length != 1 {
		t.Fatalf("TileToPosition returned invalid path length: %v", tile)
	}
	if tile.Paths[1].Id != 1172 {
		t.Fatalf("TileToPosition returned invalid path id: %v", tile)
	}
	if tile.Paths[1].X != 5 || tile.Paths[1].Y != 7 {
		t.Fatalf("TileToPosition returned invalid path coordinates: %v", tile)
	}
	if tile.Paths[2].Direction != types.MazeDirectionWest {
		t.Fatalf("TileToPosition returned invalid path: %v", tile)
	}
	if tile.Paths[2].Length != 3 {
		t.Fatalf("TileToPosition returned invalid path length: %v", tile)
	}
	if tile.Paths[2].Id != 2526 {
		t.Fatalf("TileToPosition returned invalid path id: %v", tile)
	}
	if tile.Paths[2].X != 3 || tile.Paths[2].Y != 7 {
		t.Fatalf("TileToPosition returned invalid path coordinates: %v", tile)
	}
}

// Create a maze and test that it exists.
func createMaze(t *testing.T) *Maze {
	t.Helper()
	cfg := config.Maze{
		VisibilityRange: 3,
		Configs: []config.MazeConfig{
			getMazeCfg(t),
		},
	}
	maze := NewMaze(&cfg)
	if maze == nil {
		t.Fatalf("Failed to create maze")
	}
	return maze
}

// Get the maze configuration.
func getMazeCfg(t *testing.T) config.MazeConfig {
	t.Helper()
	definition := `
{
  "name": "Red Squirrel",
  "address": "0x6dDb82e5B91e5941f4633Cb3ee49560Fb4582EFA",
  "width": 8,
  "height": 8,
  "entry": 0,
  "exit": 24,
  "tiles": [
    {
      "id": 3145,
      "position": {
        "x": 0,
        "y": 0
      },
      "paths": {
        "north": null,
        "south": 8,
        "east": null,
        "west": null
      }
    },
    {
      "id": 1672,
      "position": {
        "x": 1,
        "y": 0
      },
      "paths": {
        "north": null,
        "south": 9,
        "east": 2,
        "west": null
      }
    },
    {
      "id": 3578,
      "position": {
        "x": 2,
        "y": 0
      },
      "paths": {
        "north": null,
        "south": null,
        "east": 3,
        "west": 1
      }
    },
    {
      "id": 3343,
      "position": {
        "x": 3,
        "y": 0
      },
      "paths": {
        "north": null,
        "south": 11,
        "east": null,
        "west": 2
      }
    },
    {
      "id": 2062,
      "position": {
        "x": 4,
        "y": 0
      },
      "paths": {
        "north": null,
        "south": 12,
        "east": 5,
        "west": null
      }
    },
    {
      "id": 2649,
      "position": {
        "x": 5,
        "y": 0
      },
      "paths": {
        "north": null,
        "south": null,
        "east": 6,
        "west": 4
      }
    },
    {
      "id": 2997,
      "position": {
        "x": 6,
        "y": 0
      },
      "paths": {
        "north": null,
        "south": null,
        "east": 7,
        "west": 5
      }
    },
    {
      "id": 250,
      "position": {
        "x": 7,
        "y": 0
      },
      "paths": {
        "north": null,
        "south": 15,
        "east": null,
        "west": 6
      }
    },
    {
      "id": 1809,
      "position": {
        "x": 0,
        "y": 1
      },
      "paths": {
        "north": 0,
        "south": 16,
        "east": null,
        "west": null
      }
    },
    {
      "id": 280,
      "position": {
        "x": 1,
        "y": 1
      },
      "paths": {
        "north": 1,
        "south": null,
        "east": 10,
        "west": null
      }
    },
    {
      "id": 3845,
      "position": {
        "x": 2,
        "y": 1
      },
      "paths": {
        "north": null,
        "south": 18,
        "east": null,
        "west": 9
      }
    },
    {
      "id": 1115,
      "position": {
        "x": 3,
        "y": 1
      },
      "paths": {
        "north": 3,
        "south": null,
        "east": 12,
        "west": null
      }
    },
    {
      "id": 3309,
      "position": {
        "x": 4,
        "y": 1
      },
      "paths": {
        "north": 4,
        "south": null,
        "east": 13,
        "west": 11
      }
    },
    {
      "id": 2575,
      "position": {
        "x": 5,
        "y": 1
      },
      "paths": {
        "north": null,
        "south": null,
        "east": 14,
        "west": 12
      }
    },
    {
      "id": 1718,
      "position": {
        "x": 6,
        "y": 1
      },
      "paths": {
        "north": null,
        "south": 22,
        "east": null,
        "west": 13
      }
    },
    {
      "id": 1675,
      "position": {
        "x": 7,
        "y": 1
      },
      "paths": {
        "north": 7,
        "south": 23,
        "east": null,
        "west": null
      }
    },
    {
      "id": 503,
      "position": {
        "x": 0,
        "y": 2
      },
      "paths": {
        "north": 8,
        "south": null,
        "east": 17,
        "west": null
      }
    },
    {
      "id": 2321,
      "position": {
        "x": 1,
        "y": 2
      },
      "paths": {
        "north": null,
        "south": 25,
        "east": null,
        "west": 16
      }
    },
    {
      "id": 3988,
      "position": {
        "x": 2,
        "y": 2
      },
      "paths": {
        "north": 10,
        "south": null,
        "east": 19,
        "west": null
      }
    },
    {
      "id": 2247,
      "position": {
        "x": 3,
        "y": 2
      },
      "paths": {
        "north": null,
        "south": null,
        "east": null,
        "west": 18
      }
    },
    {
      "id": 3552,
      "position": {
        "x": 4,
        "y": 2
      },
      "paths": {
        "north": null,
        "south": 28,
        "east": 21,
        "west": null
      }
    },
    {
      "id": 1863,
      "position": {
        "x": 5,
        "y": 2
      },
      "paths": {
        "north": null,
        "south": 29,
        "east": null,
        "west": 20
      }
    },
    {
      "id": 121,
      "position": {
        "x": 6,
        "y": 2
      },
      "paths": {
        "north": 14,
        "south": null,
        "east": null,
        "west": null
      }
    },
    {
      "id": 2548,
      "position": {
        "x": 7,
        "y": 2
      },
      "paths": {
        "north": 15,
        "south": 31,
        "east": null,
        "west": null
      }
    },
    {
      "id": 2145,
      "position": {
        "x": 0,
        "y": 3
      },
      "paths": {
        "north": null,
        "south": 32,
        "east": null,
        "west": null
      }
    },
    {
      "id": 1846,
      "position": {
        "x": 1,
        "y": 3
      },
      "paths": {
        "north": 17,
        "south": 33,
        "east": null,
        "west": null
      }
    },
    {
      "id": 1247,
      "position": {
        "x": 2,
        "y": 3
      },
      "paths": {
        "north": null,
        "south": 34,
        "east": 27,
        "west": null
      }
    },
    {
      "id": 3569,
      "position": {
        "x": 3,
        "y": 3
      },
      "paths": {
        "north": null,
        "south": null,
        "east": 28,
        "west": 26
      }
    },
    {
      "id": 3083,
      "position": {
        "x": 4,
        "y": 3
      },
      "paths": {
        "north": 20,
        "south": null,
        "east": null,
        "west": 27
      }
    },
    {
      "id": 3509,
      "position": {
        "x": 5,
        "y": 3
      },
      "paths": {
        "north": 21,
        "south": 37,
        "east": null,
        "west": null
      }
    },
    {
      "id": 1781,
      "position": {
        "x": 6,
        "y": 3
      },
      "paths": {
        "north": null,
        "south": 38,
        "east": 31,
        "west": null
      }
    },
    {
      "id": 3952,
      "position": {
        "x": 7,
        "y": 3
      },
      "paths": {
        "north": 23,
        "south": 39,
        "east": null,
        "west": 30
      }
    },
    {
      "id": 3010,
      "position": {
        "x": 0,
        "y": 4
      },
      "paths": {
        "north": 24,
        "south": 40,
        "east": null,
        "west": null
      }
    },
    {
      "id": 2603,
      "position": {
        "x": 1,
        "y": 4
      },
      "paths": {
        "north": 25,
        "south": null,
        "east": 34,
        "west": null
      }
    },
    {
      "id": 1106,
      "position": {
        "x": 2,
        "y": 4
      },
      "paths": {
        "north": 26,
        "south": null,
        "east": null,
        "west": 33
      }
    },
    {
      "id": 3007,
      "position": {
        "x": 3,
        "y": 4
      },
      "paths": {
        "north": null,
        "south": 43,
        "east": 36,
        "west": null
      }
    },
    {
      "id": 1860,
      "position": {
        "x": 4,
        "y": 4
      },
      "paths": {
        "north": null,
        "south": 44,
        "east": null,
        "west": 35
      }
    },
    {
      "id": 1546,
      "position": {
        "x": 5,
        "y": 4
      },
      "paths": {
        "north": 29,
        "south": null,
        "east": 38,
        "west": null
      }
    },
    {
      "id": 113,
      "position": {
        "x": 6,
        "y": 4
      },
      "paths": {
        "north": 30,
        "south": null,
        "east": null,
        "west": 37
      }
    },
    {
      "id": 2420,
      "position": {
        "x": 7,
        "y": 4
      },
      "paths": {
        "north": 31,
        "south": 47,
        "east": null,
        "west": null
      }
    },
    {
      "id": 450,
      "position": {
        "x": 0,
        "y": 5
      },
      "paths": {
        "north": 32,
        "south": 48,
        "east": 41,
        "west": null
      }
    },
    {
      "id": 1138,
      "position": {
        "x": 1,
        "y": 5
      },
      "paths": {
        "north": null,
        "south": 49,
        "east": 42,
        "west": 40
      }
    },
    {
      "id": 3873,
      "position": {
        "x": 2,
        "y": 5
      },
      "paths": {
        "north": null,
        "south": 50,
        "east": null,
        "west": 41
      }
    },
    {
      "id": 3123,
      "position": {
        "x": 3,
        "y": 5
      },
      "paths": {
        "north": 35,
        "south": 51,
        "east": null,
        "west": null
      }
    },
    {
      "id": 3592,
      "position": {
        "x": 4,
        "y": 5
      },
      "paths": {
        "north": 36,
        "south": null,
        "east": 45,
        "west": null
      }
    },
    {
      "id": 4058,
      "position": {
        "x": 5,
        "y": 5
      },
      "paths": {
        "north": null,
        "south": null,
        "east": 46,
        "west": 44
      }
    },
    {
      "id": 87,
      "position": {
        "x": 6,
        "y": 5
      },
      "paths": {
        "north": null,
        "south": null,
        "east": 47,
        "west": 45
      }
    },
    {
      "id": 1731,
      "position": {
        "x": 7,
        "y": 5
      },
      "paths": {
        "north": 39,
        "south": null,
        "east": null,
        "west": 46
      }
    },
    {
      "id": 387,
      "position": {
        "x": 0,
        "y": 6
      },
      "paths": {
        "north": 40,
        "south": 56,
        "east": null,
        "west": null
      }
    },
    {
      "id": 3217,
      "position": {
        "x": 1,
        "y": 6
      },
      "paths": {
        "north": 41,
        "south": 57,
        "east": null,
        "west": null
      }
    },
    {
      "id": 1953,
      "position": {
        "x": 2,
        "y": 6
      },
      "paths": {
        "north": 42,
        "south": null,
        "east": null,
        "west": null
      }
    },
    {
      "id": 3266,
      "position": {
        "x": 3,
        "y": 6
      },
      "paths": {
        "north": 43,
        "south": null,
        "east": 52,
        "west": null
      }
    },
    {
      "id": 396,
      "position": {
        "x": 4,
        "y": 6
      },
      "paths": {
        "north": null,
        "south": 60,
        "east": null,
        "west": 51
      }
    },
    {
      "id": 385,
      "position": {
        "x": 5,
        "y": 6
      },
      "paths": {
        "north": null,
        "south": 61,
        "east": 54,
        "west": null
      }
    },
    {
      "id": 94,
      "position": {
        "x": 6,
        "y": 6
      },
      "paths": {
        "north": null,
        "south": null,
        "east": 55,
        "west": 53
      }
    },
    {
      "id": 3138,
      "position": {
        "x": 7,
        "y": 6
      },
      "paths": {
        "north": null,
        "south": 63,
        "east": null,
        "west": 54
      }
    },
    {
      "id": 6,
      "position": {
        "x": 0,
        "y": 7
      },
      "paths": {
        "north": 48,
        "south": null,
        "east": null,
        "west": null
      }
    },
    {
      "id": 1797,
      "position": {
        "x": 1,
        "y": 7
      },
      "paths": {
        "north": 49,
        "south": null,
        "east": 58,
        "west": null
      }
    },
    {
      "id": 1163,
      "position": {
        "x": 2,
        "y": 7
      },
      "paths": {
        "north": null,
        "south": null,
        "east": 59,
        "west": 57
      }
    },
    {
      "id": 2526,
      "position": {
        "x": 3,
        "y": 7
      },
      "paths": {
        "north": null,
        "south": null,
        "east": 60,
        "west": 58
      }
    },
    {
      "id": 882,
      "position": {
        "x": 4,
        "y": 7
      },
      "paths": {
        "north": 52,
        "south": null,
        "east": 61,
        "west": 59
      }
    },
    {
      "id": 1172,
      "position": {
        "x": 5,
        "y": 7
      },
      "paths": {
        "north": 53,
        "south": null,
        "east": null,
        "west": 60
      }
    },
    {
      "id": 574,
      "position": {
        "x": 6,
        "y": 7
      },
      "paths": {
        "north": null,
        "south": null,
        "east": 63,
        "west": null
      }
    },
    {
      "id": 1702,
      "position": {
        "x": 7,
        "y": 7
      },
      "paths": {
        "north": 55,
        "south": null,
        "east": null,
        "west": 62
      }
    }
  ]
}`
	// parse the maze definition
	var mazeCfg config.MazeConfig
	err := json.Unmarshal([]byte(definition), &mazeCfg)
	if err != nil {
		t.Fatalf("Failed to parse maze definition: %v", err)
	}
	return mazeCfg
}
