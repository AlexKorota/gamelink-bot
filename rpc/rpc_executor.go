package rpc

import (
	"context"
	"gamelink-bot/config"
	"gamelink-bot/generalcmd"
	"gamelink-bot/iface"
	msg "gamelink-go/proto_msg"
	service "gamelink-go/proto_service"
	"google.golang.org/grpc"
	"log"
)

type (
	RpcWorker struct {
		client service.AdminServiceClient
	}
)

func init() {
	w := NewRpcWorker()
	generalcmd.SetExecutor(w)
}

func NewRpcWorker() iface.GeneralExecutor {
	conn, err := grpc.Dial(config.GRPCDialAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	client := service.NewAdminServiceClient(conn)
	if client == nil { //Но это не точно!
		log.Fatal("connection error")
	}
	return &RpcWorker{client: client}
}

func (r RpcWorker) Count(ctx context.Context, params []*msg.OneCriteriaStruct) (*msg.CountResponse, error) {
	data, err := r.client.Count(ctx, &msg.MultiCriteriaRequest{Params: params})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r RpcWorker) Delete(ctx context.Context, params []*msg.OneCriteriaStruct) (*msg.OneUserResponse, error) {
	data, err := r.client.Delete(ctx, &msg.MultiCriteriaRequest{Params: params})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r RpcWorker) Find(ctx context.Context, params []*msg.OneCriteriaStruct) (*msg.MultiUserResponse, error) {
	data, err := r.client.Find(ctx, &msg.MultiCriteriaRequest{Params: params})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r RpcWorker) Update(ctx context.Context, findParams []*msg.OneCriteriaStruct, updParams []*msg.UpdateCriteriaStruct) (*msg.StringResponse, error) {
	data, err := r.client.Update(ctx, &msg.UpdateCriteriaRequest{FindParams: findParams, UpdParams: updParams})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r RpcWorker) SendPush(ctx context.Context, params []*msg.OneCriteriaStruct, message string) (*msg.StringResponse, error) {
	data, err := r.client.SendPush(ctx, &msg.PushCriteriaRequest{Params: params, Message: message})
	if err != nil {
		return nil, err
	}
	return data, nil
}
