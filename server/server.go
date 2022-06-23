package server

import (
	"context"

	"github.com/sayshu-7s/grpc-gateway-example/gen/go/example"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ExampleAPIServer struct {
	nextID int64
	msgs   map[int64]*example.ExampleMessage
	example.UnimplementedExampleApiServer
}

func NewExampleAPIServer() (*ExampleAPIServer, error) {
	msgs := make(map[int64]*example.ExampleMessage)
	// use id=1 so that you can try out get method.
	msgs[1] = &example.ExampleMessage{Id: 1, ExampleField: "example"}
	return &ExampleAPIServer{msgs: msgs, nextID: 2}, nil
}

func (s ExampleAPIServer) GetMessage(ctx context.Context, r *example.GetMessageRequest) (*example.ExampleMessage, error) {
	msg, ok := s.msgs[r.GetId()]
	if !ok {
		return nil, status.Error(codes.NotFound, codes.NotFound.String())
	}
	return msg, nil
}

func (s ExampleAPIServer) BatchGetMessages(r *example.BatchGetMessagesRequest, stream example.ExampleApi_BatchGetMessagesServer) error {
	ids := r.GetIds()
	for _, id := range ids {
		msg, ok := s.msgs[id]
		var res *example.BatchGetMessagesResponse

		if ok {
			res = &example.BatchGetMessagesResponse{Result: &example.BatchGetMessagesResponse_Found{Found: msg}}
		} else {
			res = &example.BatchGetMessagesResponse{Result: &example.BatchGetMessagesResponse_Missing{Missing: id}}
		}
		if err := stream.Send(res); err != nil {
			return status.Error(codes.Internal, codes.Internal.String())
		}
	}
	return nil
}

func (s ExampleAPIServer) CreateMessage(ctx context.Context, r *example.CreateMessageRequest) (*example.ExampleMessage, error) {
	msg := r.GetMessage()
	if msg == nil {
		status.Error(codes.InvalidArgument, codes.InvalidArgument.String())
	}

	if msg.GetId() == 0 {
		msg.Id = s.nextID
	} else if _, ok := s.msgs[msg.GetId()]; ok {
		return nil, status.Error(codes.AlreadyExists, codes.AlreadyExists.String())
	}
	s.nextID = msg.GetId() + 1

	s.msgs[msg.GetId()] = msg
	return msg, nil
}

func (s ExampleAPIServer) DeleteMessage(ctx context.Context, r *example.DeleteMessageRequest) (*emptypb.Empty, error) {
	delete(s.msgs, r.GetId())
	return new(emptypb.Empty), nil
}
