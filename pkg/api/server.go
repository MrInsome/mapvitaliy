package api

import (
	"apitraning/pkg"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type serverStruct struct {
	UnimplementedAccountServiceServer
	*pkg.Repository
	accountServiceClient
}

func (s *serverStruct) Unsubscribe(ctx context.Context, req *UnsubscribeRequest) (*UnsubscribeResponse, error) {
	err := s.UnsubscribeAccount(int(req.AccountId))
	if err != nil {
		return &UnsubscribeResponse{Success: false}, nil
	}
	return &UnsubscribeResponse{Success: true}, nil
}

func OpenGRPC(r *pkg.Repository) {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("%v", err)
	}
	s := grpc.NewServer()
	srv := &serverStruct{UnimplementedAccountServiceServer{}, r, accountServiceClient{}}
	RegisterAccountServiceServer(s, srv)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("%v", err)
	}
}
