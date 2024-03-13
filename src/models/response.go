package models

type UpdateResponse struct {
	MatchedCount  int64       `json:"MatchedCount"`
	ModifiedCount int64       `json:"ModifiedCount"`
	UpsertedCount int64       `json:"UpsertedCount"`
	UpsertedID    interface{} `json:"UpsertedID,omitempty"`
}
