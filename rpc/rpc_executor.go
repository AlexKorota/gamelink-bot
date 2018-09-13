package rpc

import (
	"context"
	"gamelinkBot/config"
	"gamelinkBot/generalcmd"
	"gamelinkBot/iface"
	"google.golang.org/grpc"
	"log"
)

type (
	RpcWorker struct {
		client AdminServiceClient
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
	client := NewAdminServiceClient(conn)
	if client == nil { //Но это не точно!
		log.Fatal("connection error")
	}
	return &RpcWorker{client: client}
}

func (r RpcWorker) Count(ctx context.Context, params []*iface.OneCriteriaStruct) (*iface.CountResponse, error) {
	data, err := r.client.Count(ctx, &iface.MultiCriteriaRequest{Params: params})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r RpcWorker) Delete(ctx context.Context, params []*iface.OneCriteriaStruct) (*iface.OneUserResponse, error) {
	data, err := r.client.Delete(ctx, &iface.MultiCriteriaRequest{Params: params})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r RpcWorker) Find(ctx context.Context, params []*iface.OneCriteriaStruct) (*iface.MultiUserResponse, error) {
	data, err := r.client.Find(ctx, &iface.MultiCriteriaRequest{Params: params})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r RpcWorker) Update(ctx context.Context, params []*iface.OneCriteriaStruct) (*iface.MultiUserResponse, error) {
	data, err := r.client.Update(ctx, &iface.MultiCriteriaRequest{Params: params})
	if err != nil {
		return nil, err
	}
	return data, nil
}
