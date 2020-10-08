package main

import (
	"bytes"
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/palantir/stacktrace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/ovchipkaart"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/proto/transactions-service"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type server struct {
	transactions_service.UnimplementedTransactionsServiceServer
	ovChipkaartAPIClient   ovchipkaart.APIClient
	csvTransactionsService *TransactionFetcherCSVService
}

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

	ovChipkaartAPIClient := ovchipkaart.NewAPIService(ovchipkaart.APIServiceConfig{
		ClientID:     os.Getenv("OV_CHIPKAART_API_CLIENT_ID"),
		ClientSecret: os.Getenv("OV_CHIPKAART_API_CLIENT_SECRET"),
		Locale:       "en",
		Client:       &http.Client{},
	})

	csvTransactionsService := NewTransactionFetcherCSVService(NewCSVIOReader())

	srv := grpc.NewServer()

	transactions_service.RegisterTransactionsServiceServer(srv, &server{ovChipkaartAPIClient: ovChipkaartAPIClient, csvTransactionsService: csvTransactionsService})

	log.Fatalln(srv.Serve(listener))
}

// RawRecordsFromBytes gets the raw records form a CSV file as bytes
func (s *server) RawRecordsFromBytes(_ context.Context, request *transactions_service.BytesRawRecordsRequest) (*transactions_service.RawRecordsResponse, error) {
	csvFile := bytes.NewBuffer(request.GetData())

	records, err := s.csvTransactionsService.FetchTransactionRecords(CSVTransactionFetchOptions{
		data:       csvFile,
		cardNumber: request.GetCardNumber(),
		startDate:  request.GetStartDate().AsTime(),
		endDate:    request.GetEndDate().AsTime(),
	})

	if err != nil {
		return nil, status.Error(codes.Code(stacktrace.GetCode(err)), err.Error())
	}

	return s.makeResponse(records)
}

// RawRecordsWithCredentials gets raw records from the ov-chipkaart API
func (s *server) RawRecordsWithCredentials(_ context.Context, request *transactions_service.CredentialsRawRecordsRequest) (*transactions_service.RawRecordsResponse, error) {
	records, err := s.ovChipkaartAPIClient.FetchTransactions(ovchipkaart.TransactionFetchOptions{
		Username:   request.GetUsername(),
		Password:   request.GetPassword(),
		CardNumber: request.GetCardNumber(),
		StartDate:  request.GetStartDate().AsTime(),
		EndDate:    request.GetEndDate().AsTime(),
	})

	if err != nil {
		return nil, status.Error(codes.Code(stacktrace.GetCode(err)), err.Error())
	}

	return s.makeResponse(records)
}

func (s server) makeResponse(records []ovchipkaart.RawRecord) (*transactions_service.RawRecordsResponse, error) {
	rawRecords := make([]*transactions_service.RawRecord, len(records))
	for index, record := range records {
		tDateTime, err := ptypes.TimestampProto(record.TransactionDateTime.ToTime())
		if err != nil {
			return nil, err
		}

		var fare *wrappers.DoubleValue
		if record.Fare != nil {
			fare = wrapperspb.Double(*record.Fare)
		}

		var ePurseMut *wrappers.DoubleValue
		if record.EPurseMut != nil {
			ePurseMut = wrapperspb.Double(*record.EPurseMut)
		}

		rawRecords[index] = &transactions_service.RawRecord{
			CheckInInfo:            record.CheckInInfo,
			CheckInText:            record.CheckInText,
			Fare:                   fare,
			FareCalculation:        record.FareCalculation,
			FareText:               record.FareText,
			ModalType:              record.ModalType,
			ProductInfo:            record.ProductInfo,
			ProductText:            record.ProductText,
			Pto:                    record.Pto,
			TransactionDateTime:    tDateTime,
			TransactionInfo:        record.TransactionInfo,
			TransactionName:        string(record.TransactionName),
			EPurseMut:              ePurseMut,
			EPurseMutInfo:          record.EPurseMutInfo,
			TransactionExplanation: record.TransactionExplanation,
			TransactionPriority:    record.TransactionPriority,
		}
	}

	return &transactions_service.RawRecordsResponse{RawRecords: rawRecords}, nil
}
