package main

import (
	"fmt"

	"github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"

	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"

	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"

	service2 "github.com/baoshuai123/go-micro-order/domain/service"

	"github.com/baoshuai123/go-micro-order/domain/repository"

	common "github.com/baoshuai123/go-micro-common"
	"github.com/baoshuai123/go-micro-order/handler"
	order "github.com/baoshuai123/go-micro-order/proto/order"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	consul2 "github.com/micro/go-plugins/registry/consul/v2"
	"github.com/opentracing/opentracing-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	QPS = 100
)

func main() {
	//1.配置中心
	consulConfig, err := common.GetConsulConfig("localhost", 8500, "/micro/config")
	if err != nil {
		log.Error(err)
	}
	//2.注册中心
	consul := consul2.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"localhost:8500",
		}
	})
	//3.链路追踪
	tracer, io, err := common.NewTracer("go.micro.service.order", "localhost:6831")
	if err != nil {
		log.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(tracer)
	//4.数据库配置
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlInfo.User,
		mysqlInfo.Pwd,
		mysqlInfo.Host,
		mysqlInfo.Port,
		mysqlInfo.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		log.Error(err)
	}
	//初始化表
	tableInit := repository.NewOrderRepository(db)
	if err := tableInit.InitTable(); err != nil {
		log.Error(err)
	}
	//5.创建实例
	orderDataService := service2.NewOrderDataService(tableInit)
	//6.暴露监控地址
	common.PrometheusBoot(9092)
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.order"),
		micro.Version("latest"),
		//暴露服务地址
		micro.Address(":9085"),
		//添加consul
		micro.Registry(consul),
		//添加链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		//添加限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
		//添加监控
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
	)

	// Initialise service
	service.Init()

	// Register Handler
	order.RegisterOrderHandler(service.Server(), &handler.Order{OrderDataService: orderDataService})

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
