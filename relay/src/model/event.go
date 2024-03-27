package model

type EventKind uint16

const (
	EventKindMetaData   = 0
	EventKindTextNote   = 1
	EventKindFollowList = 3
)

type Event struct {
	ID        string     `json:"id" bson:"_id"`
	PublicKey string     `json:"pubKey"`
	CreatedAt int        `json:"created_at"`
	Kind      EventKind  `json:"kind"`
	Tags      [][]string `json:"tags"`
	Content   string     `json:"content"`
	Signature string     `json:"sig"`
}
