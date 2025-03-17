package main

import (
	"context"
	g "go_kafka"
	handler "go_kafka/pkg/handlers"
	"go_kafka/pkg/kafk"
	"go_kafka/pkg/repistory"
	"go_kafka/pkg/repistory/cache"
	"go_kafka/pkg/service"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/spf13/viper"
)

// @title Todo App API
// @version 1.0
// @description API Server for TodoList Application

// @host localhost:2006
// @BasePath /
// @in header
// @name AccountTransaction
func main() {
	producerConfig := kafka.ConfigMap{
		"bootstrap.servers": viper.GetString("kafka.host"),
		"acks":              "all",
	}
	consumerConfig := kafka.ConfigMap{
		"bootstrap.servers": viper.GetString("kafka.host"),
		"group.id":          "test-group-id",
		"auto.offset.reset": "latest",
	}
	kafkaProduser := kafk.NewProducer(producerConfig)
	kafkaConsumer, err := kafk.NewConsumer(consumerConfig)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	redisConfig := cache.RedisConfig{
		Address:  "localhost:6379",
		Password: "",
		DB:       0,
	}
	redisClient := cache.RedisConnect(redisConfig)
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	dbPool := repistory.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.user"),
		DBName:   viper.GetString("db.dbname"),
		Password: viper.GetString("db.password"),
		SSLMode:  viper.GetString("db.sslmode"),
	}
	db := repistory.GetConnectionPool(ctx, dbPool)
	repos := repistory.NewRepository(db)
	services := service.NewService(*repos, redisClient, kafkaProduser, kafkaConsumer)
	handlers := handler.NewHandler(services)

	srv := new(g.Server)
	if err := srv.Run("2006", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("/home/merjen/url/config")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
