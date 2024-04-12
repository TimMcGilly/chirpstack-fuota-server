package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	fuota "github.com/chirpstack/chirpstack-fuota-server/v4/api/go"
)

func main() {
	deploymentId := "02fa88a3-3391-44a6-8633-c9eb71849d27"

	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial("localhost:8070", dialOpts...)
	if err != nil {
		panic(err)
	}

	client := fuota.NewFuotaServerServiceClient(conn)
	deploymentStatusResp, err := client.GetDeploymentStatus(context.Background(), &fuota.GetDeploymentStatusRequest{
		Id: deploymentId,
	})

	if err != nil {
		panic(err)
	}

	err = os.WriteFile(deploymentId+".json", []byte(protojson.Format(deploymentStatusResp)), 0644)

	if err != nil {
		log.Fatal(err)
	}
}
