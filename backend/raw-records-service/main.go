package main

import (
	"context"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/raw-records-service/database"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/raw-records-service/database/mongodb"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/raw-records-service/handlers"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/raw-records-service/transformers"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/logger"
	raw_records_service "github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/proto/raw-records-service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/palantir/stacktrace"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	log.Println("server running on port " + os.Getenv("SERVER_ADDRESS"))

	listener, err := net.Listen("tcp", os.Getenv("SERVER_ADDRESS"))
	if err != nil {
		log.Fatalln(err)
	}

	srv := grpc.NewServer()

	raw_records_service.RegisterRawRecordsServiceServer(srv, initializeServer())

	log.Fatalln(srv.Serve(listener))
}

func initializeServer() *handlers.Server {
	return &handlers.Server{
		DB: initializeDB(),
		Logger: initializeLogger(),
		Transformers: transformers.Transformers{},
	}
}

func initializeDB() database.DB {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(stacktrace.Propagate(err, "cannot connect to mongoDB"))
	}
	db := client.Database(os.Getenv("MONGODB_DB_NAME"))

	return mongodb.NewMongoDB(db)
}

func initializeLogger() logger.Logger {
	return logger.NewGoKitLogger(os.Stdout)
}

