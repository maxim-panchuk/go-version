package model

type ServiceInfo struct {
	PbcId                   int64
	PbcLocalName            string
	MicroserviceVersion     int64
	MicroserviceSysName     string
	MicroserviceLocalName   string
	MicroserviceDescription string
	BusinessObjectInfo      []*ObjectInfo
}
