package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"net"
	"time"
	"together/blog_server/internal/task"
	assetsPkg "together/blog_server/pkg/assets"
	"together/global"
	pb "together/proto"

	"google.golang.org/grpc"
)

var assets assetsPkg.Assets

func New(addr string) {
	fmt.Println("addr:", addr)
	assets = assetsPkg.GetInstance()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterBlogServerServer(s, new(blogServer))

	// 开启定时任务
	job := task.New([]string{global.BlogServer.IxugoDomain, global.BlogServer.WangboDomain}, context.Background())
	go cronJob(job)
	if err = s.Serve(lis); err != nil {
		panic(err)
	}
}

func cronJob(job task.BlogJob) {
	location, _ := time.LoadLocation("Asia/Shanghai")
	c := cron.New(cron.WithLocation(location))
	c.AddJob("@daily", cron.NewChain(cron.Recover(cron.DefaultLogger)).Then(job))
	c.Start()
	select {}
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
