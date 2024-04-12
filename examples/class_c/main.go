package main

import (
	"context"
	"fmt"
	"log"
	"os"
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

	payloadSize := 320
	payload := make([]byte, payloadSize)
	for i := 0; i < payloadSize; i++ {
		payload[i] = byte(payloadSize + 1 - i)
	}

	fmt.Printf("checksum: %x\n", crc32.ChecksumIEEE(payload))

	devices := make([]*fuota.DeploymentDevice, 0)

	for _, DevEui := range DevEuis {
		devices = append(devices, &fuota.DeploymentDevice{
			DevEui:    DevEui,
			McRootKey: mcRootKey.String(),
		})
	}

	fragmentationDescriptor := []byte{15, 0, 0, 0}

	client := fuota.NewFuotaServerServiceClient(conn)
	resp, err := client.CreateDeployment(context.Background(), &fuota.CreateDeploymentRequest{
		Deployment: &fuota.Deployment{
			ApplicationId:                     ApplicationId,
			Devices:                           devices,
			MulticastGroupType:                fuota.MulticastGroupType_CLASS_C,
			MulticastDr:                       5,
			MulticastFrequency:                868300000,
			MulticastGroupId:                  0,
			MulticastTimeout:                  8,
			MulticastRegion:                   fuota.Region_EU868,
			UnicastTimeout:                    ptypes.DurationProto(60 * time.Second),
			UnicastAttemptCount:               1,
			TimeBetweenMissingAns:             ptypes.DurationProto(4 * time.Second),
			FragmentationFragmentSize:         64,
			Payload:                           payload,
			FragmentationRedundancy:           0,
			FragmentationSessionIndex:         0,
			FragmentationMatrix:               0,
			FragmentationBlockAckDelay:        1,
			FragmentationDescriptor:           fragmentationDescriptor,
			RequestFragmentationSessionStatus: fuota.RequestFragmentationSessionStatus_AFTER_SESSION_TIMEOUT,
		},
	})
	if err != nil {
		panic(err)
	}

	var id uuid.UUID
	copy(id[:], resp.GetId())

	fmt.Printf("deployment created: %s\n", resp.GetId())

	deploymentId := resp.GetId()

	err = os.WriteFile(deploymentId+"-descriptor.json", []byte("{\"deploymentId\": \""+deploymentId+"\", \"fragmentationDescriptor\": \""+fmt.Sprintf("%v", fragmentationDescriptor)+"\"}"), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
