package service

import "github.com/maxim-panchuk/go-version/model"

type VersionSettingService interface {
	ApplyServiceInfo(serviceInfo *model.ServiceInfo) error

	/*
				void qvrnUpdateVersionSetting(
		        QvrnUpdateVersionSettingParam qvrnUpdateVersionSettingParam);

				Page<VersionSetting> qvrnFindLstVersionSettingByParam(
		        	Predicate predicate,
		        	Pageable pageable);

		    	VersionSetting qvrnFindVersionSettingById(
		        	sLong versionSettingId);

	*/
}
