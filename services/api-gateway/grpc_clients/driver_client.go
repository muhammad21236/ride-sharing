package grpc_clients

import (
	"os"
	pb "ride-sharing/shared/proto/driver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type driverServiceClient struct {
	Client pb.DriverServiceClient
	Conn   *grpc.ClientConn
}

func NewDriverServiceClient() (*driverServiceClient, error) {
	driverServiceURL := os.Getenv("DRIVER_SERVICE_URL")
	if driverServiceURL == "" {
		driverServiceURL = "driver-service:9095"
	}

	conn, err := grpc.NewClient(driverServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewDriverServiceClient(conn)

	return &driverServiceClient{
		Client: client,
		Conn:   conn,
	}, nil
}

func (c *driverServiceClient) Close() {
	if c.Conn != nil {
		if err := c.Conn.Close(); err != nil {
			return
		}
	}
}
