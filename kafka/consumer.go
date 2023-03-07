package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/maxim-panchuk/go-version/config"
	"github.com/maxim-panchuk/go-version/model"
	"github.com/maxim-panchuk/go-version/service"
	"github.com/segmentio/kafka-go"
)

func ListenForServiceInfo(metadataService service.VersionMetadataService, settingService service.VersionSettingService) {

	serviceName := config.GetServiceName()
	kafkaHost := config.GetKafkaHost()
	kafkaPort := config.GetKafkaPort()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{kafkaHost + ":" + kafkaPort},
		Topic:       "dqCompositionAfterStartLinkBusinessObject",
		Partition:   0,
		MinBytes:    10e3,
		MaxBytes:    10e6,
		MaxWait:     time.Second,
		StartOffset: kafka.LastOffset,
	})

	defer r.Close()

	for {
		m, err := r.FetchMessage(context.Background())
		if err != nil {
			log.Printf("error while fetching message: %v", err)
			break
		}

		var serviceInfo model.ServiceInfo
		err = json.Unmarshal(m.Value, &serviceInfo)
		if err != nil {
			log.Printf("error while unmarshaling message: %v", err)
			continue
		}

		if serviceName != "" && serviceName == serviceInfo.MicroserviceSysName {
			settingService.ApplyServiceInfo(&serviceInfo)
			metadataService.ApplyServiceInfo(&serviceInfo)
		}

		err = r.CommitMessages(context.Background(), m)
		if err != nil {
			log.Printf("error while committing message: %v", err)
		}
	}
}
