package domain

import "time"

type VersionEntity struct {
	VersionId  int64     `json:"versionId"`
	BoId       int64     `json:"boId"`
	VersionKey string    `json:"versionKey"`
	Data       string    `json:"data"`
	ObjectId   string    `json:"objectId"`
	Comment    string    `json:"comment"`
	Number     int       `json:"number"`
	Date       time.Time `json:"date"`
	Login      string    `json:"login"`
	IsCurrent  bool      `json:"isCurrent"`
	IsFirst    bool      `json:"isFirst"`
}
