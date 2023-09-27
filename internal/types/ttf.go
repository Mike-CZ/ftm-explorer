package types

// Ttf is time to finality type.
type Ttf struct {
	Timestamp int64   `bson:"timestamp" json:"timestamp"`
	Value     float64 `bson:"value" json:"value"`
}
