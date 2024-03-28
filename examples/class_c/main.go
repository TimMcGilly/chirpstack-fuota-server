package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"

	"hash/crc32"

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

	payloadSize := 128
	payload := make([]byte, payloadSize)
	for i := 0; i < payloadSize; i++ {
		payload[i] = byte(payloadSize + 1 - i)
	}

	fmt.Printf("checksum: %x", crc32.ChecksumIEEE(payload))

	client := fuota.NewFuotaServerServiceClient(conn)
	resp, err := client.CreateDeployment(context.Background(), &fuota.CreateDeploymentRequest{
		Deployment: &fuota.Deployment{
			ApplicationId: ApplicationId,
			Devices: []*fuota.DeploymentDevice{
				{
					DevEui:    DevEui1,
					McRootKey: mcRootKey.String(),
				},
				{
					DevEui:    DevEui2,
					McRootKey: mcRootKey.String(),
				},
			},
			MulticastGroupType:                fuota.MulticastGroupType_CLASS_C,
			MulticastDr:                       3,
			MulticastFrequency:                868300000,
			MulticastGroupId:                  0,
			MulticastTimeout:                  9,
			MulticastRegion:                   fuota.Region_EU868,
			UnicastTimeout:                    ptypes.DurationProto(60 * time.Second),
			UnicastAttemptCount:               1,
			FragmentationFragmentSize:         64,
			Payload:                           payload,
			FragmentationRedundancy:           0,
			FragmentationSessionIndex:         0,
			FragmentationMatrix:               0,
			FragmentationBlockAckDelay:        1,
			FragmentationDescriptor:           []byte{0, 0, 0, 2},
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
