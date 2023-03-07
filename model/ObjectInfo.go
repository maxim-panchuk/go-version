package model

type ObjectInfo struct {
	BoId            int64
	BoSysName       string
	BoLocalName     string
	BoBriefName     string
	BoDescription   string
	BoAttributeInfo []*AttributeInfo
	UseVersion      bool
}
