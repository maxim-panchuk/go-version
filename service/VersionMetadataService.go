package service

import "github.com/maxim-panchuk/go-version/model"

const (
	BO_ID                        = "boId"
	BO_SYS_NAME                  = "boSysName"
	BO_LOCAL_NAME                = "boLocalName"
	BO_DESCRIPTION               = "boDescription"
	USE_VERSION                  = "useVersion"
	BO_ATTRIBUTE_ID              = "boAttributeId"
	BO_ATTRIBUTE_SYS_NAME        = "boAttributeSysName"
	BO_ATTRIBUTE_LOCAL_NAME      = "boAttributeLocalName"
	IS_LINKED                    = "isLinked"
	LINKED_MICROSERVICE_SYS_NAME = "linkedMicroserviceSysName"
	LINKED_BO_SYS_NAME           = "linkedBoSysName"
	KEY_ATTRIBUTE                = "keyAttribute"
	ATTRIBUTES                   = "attributes"
)

type VersionMetadataService interface {
	ApplyServiceInfo(serviceInfo *model.ServiceInfo)
	UseVersionByBoSysName(boSysName string) bool
	KeyAttributeByBoSysName(boSysName string) string
	BoIdByBoSysName(boSysName string) int64
	AttributesByBoId(boId int64) map[string]interface{}
	UpdateUseVersion(boSysName string, useVersion bool)
	MicroserviceSysName() string
	BoSysNameByBoId(boId int64) string
}
