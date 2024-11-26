package ext_grpc

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetGRPCProvider(host, port string) (grpc.ClientConnInterface, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return conn, nil
}
