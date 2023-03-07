package config

import "os"

type Config struct {
	ServiceName    string
	DbName         string
	SchemeName     string
	DbPort         string
	DbHost         string
	DbLogin        string
	DbPassword     string
	DbDriver       string
	KafkaHost      string
	KafkaPort      string
	KafkaJsonTopic string
}

var Conf Config

func InitConfig() {
	Conf = Config{
		ServiceName:    os.Getenv("SERVICE"),
		DbName:         os.Getenv("DB_DATABASE"),
		SchemeName:     os.Getenv("DB_SCHEMA"),
		DbPort:         os.Getenv("DB_PORT"),
		DbHost:         os.Getenv("DB_HOST"),
		DbLogin:        os.Getenv("DB_LOGIN"),
		DbPassword:     os.Getenv("DB_PASSWORD"),
		DbDriver:       os.Getenv("DB_DRIVER"),
		KafkaHost:      os.Getenv("KAFKA_HOST"),
		KafkaPort:      os.Getenv("KAFKA_PORT"),
		KafkaJsonTopic: os.Getenv("KAFKA_TOPIC"),
	}
}

func GetServiceName() string {
	return Conf.ServiceName
}

func GetDbName() string {
	return Conf.DbName
}

func GetSchemeName() string {
	return Conf.SchemeName
}

func GetDbPort() string {
	return Conf.DbPort
}

func GetDbHost() string {
	return Conf.DbHost
}

func GetDbLogin() string {
	return Conf.DbLogin
}

func GetDbPassword() string {
	return Conf.DbPassword
}

func GetDbDriver() string {
	return Conf.DbDriver
}

func GetKafkaHost() string {
	return Conf.KafkaHost
}

func GetKafkaPort() string {
	return Conf.KafkaPort
}

func GetKafkaJsonTopic() string {
	return Conf.KafkaJsonTopic
}
