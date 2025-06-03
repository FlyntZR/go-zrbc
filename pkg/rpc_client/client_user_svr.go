package rpcclient

import (
	"fmt"

	pb "go-zrbc/pkg/protobuf"

	"google.golang.org/grpc"
)

type UserSvrClient struct {
	conn   *grpc.ClientConn
	Handle pb.UserServiceClient
}

// "user-rpc-service.default.svc.cluster.local:9090"
func NewUserSvrClient(addr string) (*UserSvrClient, error) {
	if addr == "" {
		return nil, fmt.Errorf("empty addr!!")
	}
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		//xlog.Error("rpc did not connect: %v", err)
		return nil, err
	}
	//defer conn.Close()

	client := UserSvrClient{conn: conn}
	client.Handle = pb.NewUserServiceClient(conn)
	//client := pb.NewUserServiceClient(conn)
	return &client, nil
}

func (c UserSvrClient) Close() {
	c.conn.Close()
}
