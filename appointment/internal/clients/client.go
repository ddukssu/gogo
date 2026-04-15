package clients

import (
	"context"
	pb "github.com/ddukssu/gogo/doc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type DoctorClient struct {
	client pb.DoctorServiceClient
}

func NewDoctorClient(target string) (*DoctorClient, error) {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &DoctorClient{client: pb.NewDoctorServiceClient(conn)}, nil
}

func (c *DoctorClient) CheckDoctorExists(doctorID string) (bool, error) {
	_, err := c.client.GetDoctor(context.Background(), &pb.GetDoctorRequest{Id: doctorID})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return false, nil
		}
		return false, status.Errorf(codes.Unavailable, "doctor service unavailable")
	}
	return true, nil
}
