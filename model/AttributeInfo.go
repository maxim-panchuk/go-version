package model

type AttributeInfo struct {
	BoAttributeId             int64
	BoAttributeSysName        string
	BoAttributeLocalName      string
	BoAttributeDescription    string
	BoAttributeIsKey          int64
	DataTypeName              string
	LinkedBoId                int64
	LinkedMicroserviceSysName string
	LinkedBoSysName           string
}
