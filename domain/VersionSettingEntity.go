package domain

type VersionSettingEntity struct {
	BoId             int64  `json:"boId"`
	UseVersion       bool   `json:"useVersion"`
	VersionSettingId int64  `json:"versionSettingId" gorm:"primary_key"`
	BoSysName        string `json:"boSysName"`
	BoLocalName      string `json:"boLocalName"`
	BoDescription    string `json:"boDescription"`
}
