package db

import "github.com/maxim-panchuk/go-version/domain"

type VersionSettingEntityRepository interface {
	FindVersionSettingEntityByBoId(boId int64) (*domain.VersionSettingEntity, error)
	FindVersionSettingEntityByBoSysName(boSysName string) (*domain.VersionSettingEntity, error)
	SaveVersionSettingEntity(versionSettingEntity *domain.VersionSettingEntity) error
	FindAllVersionSettingEntities() ([]*domain.VersionSettingEntity, error)
	DeleteVersionSettingEntity(versionSettingEntity domain.VersionSettingEntity) error
}
