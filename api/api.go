package api

import (
	"errors"
	"sync"
	"time"

	"gitflex.diasoft.ru/mvp-go/golang-libraries/go-version/db"
	"gitflex.diasoft.ru/mvp-go/golang-libraries/go-version/domain"
	"gitflex.diasoft.ru/mvp-go/golang-libraries/go-version/model"
	"gitflex.diasoft.ru/mvp-go/golang-libraries/go-version/service"
	"github.com/google/uuid"
)

type SaveVersion interface {
	LatchVersion(objectSysName, objectId, jsonedObject string) error
}

type saveVersion struct {
	versionMetadataService  service.VersionMetadataService
	versionEntityRepository db.VersionEntityRepository
}

var (
	singleton = saveVersion{}
	once      sync.Once
)

func newSaveVersion() *saveVersion {
	return &saveVersion{
		versionMetadataService:  service.GetVersionMetadateService(),
		versionEntityRepository: db.GetRepo(nil),
	}
}

func GetSaveVersion() SaveVersion {
	once.Do(func() {
		singleton = *newSaveVersion()
	})
	return &singleton
}

func (s *saveVersion) LatchVersion(objectSysName, objectId, jsonedObject string) error {

	if objectSysName == "" || objectId == "" || jsonedObject == "" {
		return errors.New("SaveVersion: objectSysName or objectId or jsonedObject is nil")
	}

	if s.versionMetadataService.UseVersionByBoSysName(objectSysName) {
		boId := s.versionMetadataService.BoIdByBoSysName(objectSysName)
		if boId != 0 {
			versionKey := uuid.New().String()
			err := s.saveVersionData(boId, objectId, versionKey, jsonedObject)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *saveVersion) saveVersionData(boId int64, objectId, versionKey, data string) error {

	currentTime := time.Now()

	versionModel := model.VersionModel{
		BoID:       boId,
		VersionKey: versionKey,
		Data:       data,
		ObjectID:   objectId,
		Date:       currentTime,
	}

	err := s.saveData(versionModel)
	if err != nil {
		return err
	}

	return nil

}

func (s *saveVersion) saveData(versionModel model.VersionModel) error {

	if versionModel.BoID == 0 {
		return errors.New("BoID is not presented")
	}

	if versionModel.ObjectID == "" {
		return errors.New("ObjectId is not presented")
	}

	if exists, err := s.versionEntityRepository.ExistsByBoIdAndObjectId(versionModel.BoID, versionModel.ObjectID); !exists || err != nil {
		if err != nil {
			return err
		}

		versionEntity := domain.VersionEntity{
			BoId:       versionModel.BoID,
			VersionKey: versionModel.VersionKey,
			Data:       versionModel.Data,
			ObjectId:   versionModel.ObjectID,
			Comment:    versionModel.Comment,
			Number:     1,
			Date:       versionModel.Date,
			Login:      "",
			IsCurrent:  true,
			IsFirst:    true,
		}

		return s.versionEntityRepository.SaveVersionEntity(&versionEntity)
	}

	list, err := s.versionEntityRepository.FindVersionEntityByBoIdAndObjectIdAndMaxNumber(versionModel.BoID, versionModel.ObjectID)
	if err != nil {
		return err
	}

	var number int = 1

	if len(list) != 0 {
		oldVersion := list[0]
		oldVersion.Data = versionModel.Data
		oldVersion.IsCurrent = false
		number = oldVersion.Number + 1

		if err := s.versionEntityRepository.SaveVersionEntity(&oldVersion); err != nil {
			return err
		}
	}

	entity := domain.VersionEntity{
		BoId:       versionModel.BoID,
		VersionKey: versionModel.VersionKey,
		Data:       "",
		ObjectId:   versionModel.ObjectID,
		Comment:    versionModel.Comment,
		Number:     number,
		Date:       versionModel.Date,
		Login:      "",
		IsCurrent:  true,
		IsFirst:    false,
	}

	if err = s.versionEntityRepository.SaveVersionEntity(&entity); err != nil {
		return err
	}

	return nil
}
