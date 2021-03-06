package main

import (
	"context"
	rawRecordsService "github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/proto/raw-records-service"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql/handler/apollotracing"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/proto/transactions-service"
	"google.golang.org/grpc"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/ovchipkaart"

	"os"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/middlewares"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/graph/validator"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/graph/validator/govalidator"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/services/password"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/errorhandler"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/logger"
	"github.com/getsentry/sentry-go"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/cache"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/cache/redis"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/database"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/database/mongodb"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/graph/generated"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/graph/resolver"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/services/jwt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultPort = "8080"

type Singletons struct {
	errorHandler errorhandler.ErrorHandler
}

var (
	singletons = Singletons{}
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := mux.NewRouter()

	middlewareClient := initializeMiddlewares()

	router.Use(middlewareClient.LogRequest(initializeLogger()))
	router.Use(middlewareClient.EnrichUserID(initializeJWTService()))
	router.Use(middlewareClient.AddLanguageTag())

	router.HandleFunc("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", initializeGraphQLServer())

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, cors.AllowAll().Handler(router)))
}

func initializeGraphQLServer() *handler.Server {
	server := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: initializeResolver(),
			},
		),
	)
	server.Use(apollotracing.Tracer{})
	return server
}
func initializeResolver() *resolver.Resolver {
	return resolver.NewResolver(
		initializeDB(),
		initializeValidator(),
		initializePasswordService(),
		initializeErrorHandler(),
		initializeLogger(),
		initializeJWTService(),
		initializeTransactionsServiceClient(),
		initializeRawRecordsServiceClient(),
	)
}

func initializeOvChipkaartAPIClient() ovchipkaart.APIClient {
	return ovchipkaart.NewAPIService(ovchipkaart.APIServiceConfig{
		ClientID:     os.Getenv("OV_CHIPKAART_API_CLIENT_ID"),
		ClientSecret: os.Getenv("OV_CHIPKAART_API_CLIENT_SECRET"),
		Locale:       "en",
		Client:       &http.Client{},
	})
}

func initializeValidatorHelpers() validator.Helpers {
	return validator.NewHelpers(initializeOvChipkaartAPIClient())
}

func initializeValidator() validator.Validator {
	return govalidator.New(initializeDB(), initializeValidatorHelpers(), initializeErrorHandler())
}

func initializePasswordService() password.Service {
	return password.NewBcryptService()
}

func initializeJWTService() jwt.Service {
	sessionDays, err := strconv.Atoi(os.Getenv("AUTH_SESSION_DAYS"))
	if err != nil {
		log.Fatal(err.Error())
	}

	if sessionDays < 1 {
		log.Fatal("AUTH_SESSION_DAYS cannot be < 1")
	}

	return jwt.NewService(os.Getenv("JWT_SECRET"), initializeCache(), sessionDays)
}

func initializeDB() database.DB {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(errors.Wrapf(err, "cannot connect to mongoDB"))
	}

	db := client.Database(os.Getenv("MONGODB_DB_NAME"))

	return mongodb.NewMongoDB(db)
}

func initializeCache() cache.Cache {
	return redis.NewClient(redis.Options{
		Address:  os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}

func initializeMiddlewares() middlewares.Client {
	return middlewares.New()
}

func initializeLogger() logger.Logger {
	return logger.NewGoKitLogger(os.Stdout)
}

func initializeErrorHandler() errorhandler.ErrorHandler {
	if singletons.errorHandler != nil {
		return singletons.errorHandler
	}

	errorHandlerSingleton, err := errorhandler.NewSentryErrorHandler(sentry.ClientOptions{
		// Either set your DSN here or set the SENTRY_DSN environment variable.
		Dsn: os.Getenv("SENTRY_DSN"),
		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug: true,
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	singletons.errorHandler = errorHandlerSingleton
	return singletons.errorHandler
}

func initializeTransactionsServiceClient() transactions_service.TransactionsServiceClient {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	conn, err := grpc.DialContext(ctx, os.Getenv("TRANSACTIONS_SERVICE_TARGET"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}

	return transactions_service.NewTransactionsServiceClient(conn)
}

func initializeRawRecordsServiceClient() rawRecordsService.RawRecordsServiceClient {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	conn, err := grpc.DialContext(ctx, os.Getenv("RAW_RECORDS_SERVICE_TARGET"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}

	return rawRecordsService.NewRawRecordsServiceClient(conn)
}