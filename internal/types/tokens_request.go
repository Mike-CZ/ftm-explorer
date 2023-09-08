package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TokensRequest represents a request for tokens.
type TokensRequest struct {
	Id primitive.ObjectID `bson:"_id"`

	// Address is the address of the applicant.
	IpAddress string `bson:"ip"`

	// Phrase is the phrase for the applicant that will be signed.
	Phrase string `bson:"phrase"`

	// ClaimedAt is the time when the tokens were claimed.
	ClaimedAt *int64 `bson:"claimed_at"`
}
