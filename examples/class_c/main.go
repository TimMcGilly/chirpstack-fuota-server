package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"

	"github.com/brocaar/lorawan/applayer/multicastsetup"
	fuota "github.com/chirpstack/chirpstack-fuota-server/v4/api/go"
)

func main() {
	mcRootKey, err := multicastsetup.GetMcRootKeyForGenAppKey(GenAppKey)
	if err != nil {
		log.Fatal(err)
	}

	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial("localhost:8070", dialOpts...)
	if err != nil {
		panic(err)
	}

	client := fuota.NewFuotaServerServiceClient(conn)
	resp, err := client.CreateDeployment(context.Background(), &fuota.CreateDeploymentRequest{
		Deployment: &fuota.Deployment{
			ApplicationId: ApplicationId,
			Devices: []*fuota.DeploymentDevice{
				{
					DevEui:    DevEui1,
					McRootKey: mcRootKey.String(),
				},
			},
			MulticastGroupType:                fuota.MulticastGroupType_CLASS_C,
			MulticastDr:                       5,
			MulticastFrequency:                868100000,
			MulticastGroupId:                  0,
			MulticastTimeout:                  6,
			MulticastRegion:                   fuota.Region_EU868,
			UnicastTimeout:                    ptypes.DurationProto(60 * time.Second),
			UnicastAttemptCount:               1,
			FragmentationFragmentSize:         50,
			Payload:                           make([]byte, 100),
			FragmentationRedundancy:           1,
			FragmentationSessionIndex:         0,
			FragmentationMatrix:               0,
			FragmentationBlockAckDelay:        1,
			FragmentationDescriptor:           []byte{0, 0, 0, 0},
			RequestFragmentationSessionStatus: fuota.RequestFragmentationSessionStatus_AFTER_SESSION_TIMEOUT,
		},
	})
	if err != nil {
		panic(err)
	}

	var id uuid.UUID
	copy(id[:], resp.GetId())

	fmt.Printf("deployment created: %s\n", id)
}
