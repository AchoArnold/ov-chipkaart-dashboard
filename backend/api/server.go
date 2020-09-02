package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/proto/transactions"
	"google.golang.org/grpc"

	ov_chipkaart "github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/ov-chipkaart"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/middlewares"
	"github.com/gorilla/mux"

	"os"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/graph/validator"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/graph/validator/govalidator"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/services/password"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/errorhandler"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/logger"
	"github.com/getsentry/sentry-go"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/cache"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/cache/redis"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/database"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/database/mongodb"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/graph/generated"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/graph/resolver"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/services/jwt"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultPort = "8080"

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
	return handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: initializeResolver(),
			},
		),
	)
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
	)
}

func initializeOvChipkaartAPIClient() ov_chipkaart.APIClient {
	return ov_chipkaart.NewAPIService(ov_chipkaart.APIServiceConfig{
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
	errHandler, err := errorhandler.NewSentryErrorHandler(sentry.ClientOptions{
		// Either set your DSN here or set the SENTRY_DSN environment variable.
		Dsn: os.Getenv("SENTRY_DSN"),
		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug: true,
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	return errHandler
}

func initializeTransactionsServiceClient() transactions.TransactionsServiceClient {
	conn, err := grpc.Dial(os.Getenv("TRANSACTIONS_SERVICE_TARGET"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}

	return transactions.NewTransactionsServiceClient(conn)
}
