package service

import (
	"log"
	"strings"

	"github.com/maxim-panchuk/go-version/db"
	"github.com/maxim-panchuk/go-version/domain"
	"github.com/maxim-panchuk/go-version/model"
)

type VersionSettingServiceImpl struct {
	repo                   db.VersionSettingEntityRepository
	versionMetadataService VersionMetadataService
}

func NewVersionSettingServiceImpl(repo db.VersionSettingEntityRepository, versionMetadataService VersionMetadataService) *VersionSettingServiceImpl {
	return &VersionSettingServiceImpl{
		repo:                   repo,
		versionMetadataService: versionMetadataService,
	}
}

func (s *VersionSettingServiceImpl) ApplyServiceInfo(serviceInfo *model.ServiceInfo) error {
	if serviceInfo == nil || len(serviceInfo.BusinessObjectInfo) == 0 {
		return nil
	}

	validBoIds := make([]int64, 0)

	for _, objectInfo := range serviceInfo.BusinessObjectInfo {
		if objectInfo == nil {
			continue
		}

		boId := objectInfo.BoId
		boSysName := objectInfo.BoSysName

		if boId == 0 || boSysName == "" {
			continue
		}

		versionSetting, err := s.repo.FindVersionSettingEntityByBoId(boId)
		if err != nil {
			log.Printf("Error while finding version setting entity: %v", err)
			continue
		}

		if versionSetting != nil {
			entity := versionSetting
			entity.BoSysName = boSysName
			entity.BoLocalName = objectInfo.BoLocalName
			entity.BoDescription = objectInfo.BoDescription
			err = s.repo.SaveVersionSettingEntity(entity)

			if err != nil {
				log.Println(err)
				continue
			}

			objectInfo.UseVersion = entity.UseVersion
		} else {
			entity := &domain.VersionSettingEntity{
				BoId:          boId,
				BoSysName:     boSysName,
				BoLocalName:   objectInfo.BoLocalName,
				BoDescription: objectInfo.BoDescription,
				UseVersion:    false,
			}

			if strings.EqualFold(boSysName, "Individual") {
				entity.UseVersion = true
			}

			err = s.repo.SaveVersionSettingEntity(entity)

			if err != nil {
				log.Println(err)
				continue
			}

			objectInfo.UseVersion = entity.UseVersion
		}

		validBoIds = append(validBoIds, boId)
	}

	allEntities, err := s.repo.FindAllVersionSettingEntities()
	if err != nil {
		log.Println(err)
		return err
	}

	for _, entity := range allEntities {
		boId := entity.BoId

		if boId != 0 && !contains(validBoIds, boId) {
			err = s.repo.DeleteVersionSettingEntity(*entity)

			if err != nil {
				log.Printf("Error while deleting version setting entity: %v", err)
			}
		}
	}

	return nil
}

func contains(slice []int64, element int64) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}
