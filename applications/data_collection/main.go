package main

import (
	"context"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	fuota "github.com/chirpstack/chirpstack-fuota-server/v4/api/go"
)

func main() {
	deploymentId := "3c56a3b1-14b7-4bcd-981c-8f3c9684f302"

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

	err = os.WriteFile(deploymentId+"-data.json", []byte(protojson.Format(deploymentStatusResp)), 0644)

	if err != nil {
		log.Fatal(err)
	}

	deploymentLogJson := "{"

	for _, deviceStatus := range deploymentStatusResp.GetDeviceStatus() {
		deploymentLogJson += `"` + deviceStatus.GetDevEui() + `": `
		deploymentLog, err := client.GetDeploymentDeviceLogs(context.Background(), &fuota.GetDeploymentDeviceLogsRequest{
			DeploymentId: deploymentId,
			DevEui:       deviceStatus.GetDevEui(),
		})

		if err != nil {
			panic(err)
		}

		deploymentLogJson += protojson.Format(deploymentLog) + ", "
	}

	deploymentLogJson = strings.TrimSuffix(deploymentLogJson, ", ")
	deploymentLogJson += "}"

	err = os.WriteFile(deploymentId+"-logs.json", []byte(deploymentLogJson), 0644)

	if err != nil {
		log.Fatal(err)
	}
}
