package grpcclient

import (
	taskv1 "github.com/alisheikhattar/production-platform-lab/application/task-service/gen/task/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TaskClient struct {
	conn   *grpc.ClientConn
	Client taskv1.TaskServiceClient
}

func NewTaskClient(address string) (*TaskClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &TaskClient{
		conn:   conn,
		Client: taskv1.NewTaskServiceClient(conn),
	}, nil
}

func (c *TaskClient) Close() error {
	return c.conn.Close()
}
