package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/palantir/stacktrace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ovchipkaart "github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/ov-chipkaart"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/proto/transactions"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type server struct {
	transactions.UnimplementedTransactionsServiceServer
	ovChipkaartAPIClient ovchipkaart.APIClient
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	log.Println("Server running ...")

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

	srv := grpc.NewServer()

	transactions.RegisterTransactionsServiceServer(srv, &server{ovChipkaartAPIClient: ovChipkaartAPIClient})

	log.Fatalln(srv.Serve(listener))
}

// RawRecordsFromBytes gets the raw records form a CSV file as bytes
func (s *server) RawRecordsFromBytes(_ context.Context, _ *transactions.BytesRawRecordsRequest) (*transactions.RawRecordsResponse, error) {
	return nil, nil
}

// RawRecordsWithCredentials gets raw records from the ov-chipkaart API
func (s *server) RawRecordsWithCredentials(_ context.Context, request *transactions.CredentialsRawRecordsRequest) (*transactions.RawRecordsResponse, error) {
	records, err := s.ovChipkaartAPIClient.FetchTransactions(ovchipkaart.TransactionFetchOptions{
		Username:   request.Username,
		Password:   request.Password,
		CardNumber: request.CardNumber,
		StartDate:  request.StartDate.AsTime(),
		EndDate:    request.EndDate.AsTime(),
	})

	if err != nil {
		return nil, status.Error(codes.Code(stacktrace.GetCode(err)), err.Error())
	}

	return s.makeResponse(records)
}

func (s server) makeResponse(records []ovchipkaart.RawRecord) (*transactions.RawRecordsResponse, error) {
	rawRecords := make([]*transactions.RawRecord, len(records))
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

		rawRecords[index] = &transactions.RawRecord{
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

	return &transactions.RawRecordsResponse{RawRecords: rawRecords}, nil
}
