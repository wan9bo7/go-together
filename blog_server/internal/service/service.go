package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	assetsPkg "together/blog_server/pkg/assets"
	pb "together/proto"

	"google.golang.org/grpc"
)

var assets assetsPkg.Assets

func New(addr string) {
	fmt.Println("addr:", addr)
	assets = assetsPkg.GetInstance()
	// 监听指定 TCP 端口，用于接受客户端请求
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	// 创建 gRPC Server 的实例对象
	s := grpc.NewServer()
	// gRPC Server 内部服务和路由的注册
	pb.RegisterBlogServerServer(s, new(blogServer))
	// Serve() 调用服务器以执行阻塞等待，直到进程被终止或被 Stop() 调用
	if err = s.Serve(lis); err != nil {
		panic(err)
	}
}

// 监听数据

type blogServer struct {
}

func (s *blogServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + r.GetName()}, nil
}

func (s *blogServer) GetList(ctx context.Context, r *pb.GetListRequest) (*pb.GetListReply, error) {
	data := getWebsite(r.GetUrl())
	if len(data) == 0 {
		return nil, errors.New("没有数据")
	}
	return &pb.GetListReply{Next: "", Data: data}, nil
}
