package db_types

// Block represents a block in database.
// We only need a few data, so we only define those fields.
type Block struct {
	Number    int64 `bson:"_id"`
	TxCount   int32 `bson:"txCount"`
	GasUsed   int64 `bson:"gasUsed"`
	Timestamp int64 `bson:"timestamp"`
}
