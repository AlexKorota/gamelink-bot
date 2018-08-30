package command

import (
	"gamelinkBot/config"
	"gamelinkBot/prot"
	"google.golang.org/grpc"
	"log"
	"sync"
)

var (
	client  prot.AdminServiceClient
	rpcOnce sync.Once
)

func SharedClient() prot.AdminServiceClient {
	rpcOnce.Do(func() {
		conn, err := grpc.Dial(config.DialAddress, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		client = prot.NewAdminServiceClient(conn)
		if client == nil { //Но это не точно!
			log.Fatal("connection error")
		}
	})
	return client
}
