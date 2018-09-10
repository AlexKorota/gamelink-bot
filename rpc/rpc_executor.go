package rpc

import (
	"context"
	"gamelinkBot/config"
	"gamelinkBot/generalcmd"
	"gamelinkBot/iface"
	"gamelinkBot/prot"
	"google.golang.org/grpc"
	"log"
)

type (
	RpcWorker struct {
		client prot.AdminServiceClient
	}
)

func init() {
	w := NewRpcWorker()
	generalcmd.SetExecutor(w)
}

func NewRpcWorker() iface.GeneralExecutor {
	conn, err := grpc.Dial(config.DialAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	client := prot.NewAdminServiceClient(conn)
	if client == nil { //Но это не точно!
		log.Fatal("connection error")
	}
	return &RpcWorker{client: client}
}

func (r RpcWorker) Count(ctx context.Context, params []*prot.OneCriteriaStruct) (*prot.CountResponse, error) {
	data, err := r.client.Count(ctx, &prot.MultiCriteriaRequest{Params: params})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r RpcWorker) Delete(ctx context.Context, params []*prot.OneCriteriaStruct) (*prot.OneUserResponse, error) {
	data, err := r.client.Delete(ctx, &prot.MultiCriteriaRequest{Params: params})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r RpcWorker) Find(ctx context.Context, params []*prot.OneCriteriaStruct) (*prot.MultiUserResponse, error) {
	data, err := r.client.Find(ctx, &prot.MultiCriteriaRequest{Params: params})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r RpcWorker) Update(ctx context.Context, params []*prot.OneCriteriaStruct) (*prot.MultiUserResponse, error) {
	data, err := r.client.Update(ctx, &prot.MultiCriteriaRequest{Params: params})
	if err != nil {
		return nil, err
	}
	return data, nil
}
