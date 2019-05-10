package api

type URL struct {
	Link        string `bson:"link" json:"link"`
	Hash        string `bson:"hash" json:"hash"`
	Transitions int    `bson:"transitions" json:"transitions"`
	Address     string `json:"address"`
	IsNew       bool   `json:"is_new"`
}
