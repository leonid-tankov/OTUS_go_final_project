package main

import (
	"flag"

	"github.com/leonid-tankov/OTUS_go_final_project/internal/app"
	"github.com/leonid-tankov/OTUS_go_final_project/internal/config"
	"github.com/leonid-tankov/OTUS_go_final_project/internal/kafka"
	"github.com/leonid-tankov/OTUS_go_final_project/internal/logger"
	"github.com/leonid-tankov/OTUS_go_final_project/internal/server/grpc/server"
	"github.com/leonid-tankov/OTUS_go_final_project/internal/storage"
)

var (
	configFile string
	logLevel   int
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/banner-rotation/config.yaml", "Path to configuration file")
	flag.IntVar(&logLevel, "logLevel", 4, "Log level")
}

func main() {
	flag.Parse()

	logg := logger.New(logLevel)
	conf := config.NewConfig(logg, configFile)
	postgresStorage := storage.New(conf)
	kafkaProducer := kafka.NewProducer(conf)
	application := app.New(conf, logg, postgresStorage, kafkaProducer)
	serv := server.NewServer(application)
	serv.Run()
}
