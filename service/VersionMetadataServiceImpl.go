package service

import (
	"strings"
	"sync"

	"github.com/maxim-panchuk/go-version/model"
)

type VersionMetadataServiceImpl struct {
	serviceInfo         *model.ServiceInfo
	microserviceSysName string
	boIdSysNameMap      map[int64]string
	boInfoMap           map[string]interface{}
}

var (
	singleton = VersionMetadataServiceImpl{}
	once      sync.Once
)

func GetVersionMetadateService() *VersionMetadataServiceImpl {
	once.Do(func() {
		singleton = *newVersionMetadataService()
	})

	return &singleton
}

func newVersionMetadataService() *VersionMetadataServiceImpl {
	return &VersionMetadataServiceImpl{
		microserviceSysName: "some_service",
		boIdSysNameMap:      make(map[int64]string, 0),
		boInfoMap:           make(map[string]interface{}, 0),
	}
}

func (s *VersionMetadataServiceImpl) ApplyServiceInfo(info *model.ServiceInfo) {
	s.clearData()

	if info == nil {
		return
	}

	if info.MicroserviceSysName != s.microserviceSysName {
		return
	}

	s.serviceInfo = info

	for _, objectInfo := range info.BusinessObjectInfo {
		s.processObjectInfo(objectInfo)
	}
}

func (s *VersionMetadataServiceImpl) UseVersionByBoSysName(boSysName string) bool {
	result := false
	boInfo, ok := s.boInfoMap[boSysName].(map[string]interface{})
	if ok {
		useVersion, ok := boInfo[USE_VERSION]
		if ok {
			result = useVersion.(bool)
		}
	}
	return result
}

func (s *VersionMetadataServiceImpl) KeyAttributeByBoSysName(boSysName string) string {
	var result string
	boInfo, ok := s.boInfoMap[boSysName].(map[string]interface{})
	if ok {
		keyAttr, ok := boInfo[KEY_ATTRIBUTE].(string)
		if ok {
			result = keyAttr
		}
	}
	return result
}

func (s *VersionMetadataServiceImpl) BoIdByBoSysName(boSysName string) int64 {
	var result int64 = 0
	boInfo, ok := s.boInfoMap[boSysName].(map[string]interface{})
	if ok {
		boId, ok := boInfo[BO_ID].(int64)
		if ok {
			result = boId
		}
	}
	return result
}

func (s *VersionMetadataServiceImpl) AttributesByBoId(boId int64) map[string]interface{} {
	result := make(map[string]interface{})
	boSysName := s.boIdSysNameMap[boId]
	if boSysName != "" {
		boInfo, ok := s.boInfoMap[boSysName].(map[string]interface{})
		if ok {
			result, _ = boInfo[ATTRIBUTES].(map[string]interface{})
		}
	}
	return result
}

func (s *VersionMetadataServiceImpl) UpdateUseVersion(boSysName string, useVersion bool) {
	boInfo, ok := s.boInfoMap[boSysName].(map[string]interface{})
	if ok {
		boInfo[USE_VERSION] = useVersion
	}

	if s.serviceInfo != nil && s.serviceInfo.BusinessObjectInfo != nil {
		for _, objectInfo := range s.serviceInfo.BusinessObjectInfo {
			if objectInfo != nil && objectInfo.BoSysName != "" && strings.EqualFold(objectInfo.BoSysName, boSysName) {
				objectInfo.UseVersion = useVersion
				break
			}
		}
	}
}

func (s *VersionMetadataServiceImpl) MicroserviceSysName() string {
	return s.microserviceSysName
}

func (s *VersionMetadataServiceImpl) BoSysNameByBoId(boId int64) string {
	return s.boIdSysNameMap[boId]
}

func (s *VersionMetadataServiceImpl) processObjectInfo(objectInfo *model.ObjectInfo) {
	if objectInfo == nil || objectInfo.BoSysName == "" || objectInfo.BoId <= 0 {
		return
	}

	boId := objectInfo.BoId
	boSysName := objectInfo.BoSysName

	s.boIdSysNameMap[boId] = boSysName

	attributes := make(map[string]interface{})
	keyAttribute := s.getKeyAttribute(objectInfo.BoAttributeInfo)

	for _, attributeInfo := range objectInfo.BoAttributeInfo {
		if attributeInfo == nil || attributeInfo.BoAttributeSysName == "" {
			continue
		}

		attributeData := s.getAttributeData(attributeInfo)
		attributes[attributeInfo.BoAttributeSysName] = attributeData
	}

	oneBo := map[string]interface{}{
		BO_ID:          boId,
		BO_SYS_NAME:    boSysName,
		BO_LOCAL_NAME:  objectInfo.BoLocalName,
		BO_DESCRIPTION: objectInfo.BoDescription,
		USE_VERSION:    objectInfo.UseVersion,
		KEY_ATTRIBUTE:  keyAttribute,
		ATTRIBUTES:     attributes,
	}

	s.boInfoMap[boSysName] = oneBo
}

func (s *VersionMetadataServiceImpl) getAttributeData(attributeInfo *model.AttributeInfo) map[string]interface{} {
	linkedBoId := attributeInfo.LinkedBoId
	isLinked := linkedBoId != 0 && linkedBoId > 0

	return map[string]interface{}{
		BO_ATTRIBUTE_ID:              attributeInfo.BoAttributeId,
		BO_ATTRIBUTE_SYS_NAME:        attributeInfo.BoAttributeSysName,
		BO_ATTRIBUTE_LOCAL_NAME:      attributeInfo.BoAttributeLocalName,
		IS_LINKED:                    isLinked,
		LINKED_MICROSERVICE_SYS_NAME: attributeInfo.LinkedMicroserviceSysName,
		LINKED_BO_SYS_NAME:           attributeInfo.LinkedBoSysName,
	}
}

func (s *VersionMetadataServiceImpl) getKeyAttribute(boAttributeInfo []*model.AttributeInfo) string {
	for _, attributeInfo := range boAttributeInfo {
		if attributeInfo != nil && attributeInfo.BoAttributeSysName != "" && attributeInfo.BoAttributeIsKey > 0 {
			return attributeInfo.BoAttributeSysName
		}
	}

	return ""
}

func (s *VersionMetadataServiceImpl) clearData() {
	s.serviceInfo = nil
	s.microserviceSysName = ""
	for k := range s.boIdSysNameMap {
		delete(s.boIdSysNameMap, k)
	}
	for k := range s.boInfoMap {
		delete(s.boInfoMap, k)
	}
}
