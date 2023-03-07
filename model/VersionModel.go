package model

import "time"

type VersionModel struct {
	VersionID  int64     `json:"versionId"`
	BoID       int64     `json:"boId"`
	VersionKey string    `json:"versionKey"`
	Data       string    `json:"data"`
	ObjectID   string    `json:"objectId"`
	Comment    string    `json:"comment"`
	Number     int       `json:"number"`
	Date       time.Time `json:"date" format:"2006-01-02:15:04:05.000Z"`
	Login      string    `json:"login"`
	IsCurrent  bool      `json:"isCurrent"`
	IsFirst    bool      `json:"isFirst"`
}
