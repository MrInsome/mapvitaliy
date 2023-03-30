package api

import (
	"apitraning/pkg"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

type serverStruct struct {
	repository *pkg.Repository
	UnsubscribeServiceServer
}

func (s *serverStruct) Unsubscribe(ctx context.Context, account *AccountRequest) (*emptypb.Empty, error) {
	r, err := s.repository.GetAccount(int(account.AccountId))
	err = s.repository.UnsubscribeAccount(s.repository.DBReturn(), r.AccountID)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func OpenGRPC() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("%v", err)
	}
	s := grpc.NewServer()
	srv := &serverStruct{}
	RegisterUnsubscribeServiceServer(s, srv)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("%v", err)
	}
}
