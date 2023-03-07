package db

import "github.com/maxim-panchuk/go-version/domain"

type VersionEntityRepository interface {
	ExistsByBoIdAndObjectId(boId int64, objectId string) (bool, error)
	SaveVersionEntity(versionEntity *domain.VersionEntity) error
	FindVersionEntityByBoIdAndObjectIdAndMaxNumber(boId int64, objectId string) ([]domain.VersionEntity, error)
}
