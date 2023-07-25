package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"google.golang.org/grpc"

	"log-service/data"
	"log-service/logs"

	"github.com/caarlos0/env/v6"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Config struct {
	MongoURL      string `env:"MONGO_URL"`
	MongoUser     string `env:"MONGO_INITDB_ROOT_USERNAME"`
	MongoPassword string `env:"MONGO_INITDB_ROOT_PASSWORD"`
	WebPort       string `env:"WEB_PORT"`
	RpcPort       string `env:"RPC_PORT"`
	GpcPort       string `env:"GRPC_PORT"`
	Models        data.Models
}

func main() {
	app := Config{}
	if err := env.Parse(&app); err != nil {
		log.Fatal(err.Error())
	}

	mongoClient, err := app.connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app.Models = data.New(client)

	err = rpc.Register(new(RPCServer))
	go app.listenRPC()

	go app.listenGRPC()

	log.Printf("Starting log service on port %s\n", app.WebPort)
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", app.WebPort),
		Handler: app.routes(),
	}

	if err = srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func (app *Config) connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(app.MongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: app.MongoUser,
		Password: app.MongoPassword,
	})

	conn, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Printf("%e", fmt.Errorf("%w", err))
		return nil, err
	}

	return conn, nil
}

func (app *Config) listenRPC() {
	log.Println("Starting RPC server")

	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", app.RpcPort))
	if err != nil {
		log.Panic(err)
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func (app *Config) listenGRPC() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", app.GpcPort))
	if err != nil {
		log.Fatalf("Failed to listen gRPC: %v", err)
	}

	server := grpc.NewServer()

	logs.RegisterLogServiceServer(server, &LogServer{Models: app.Models})
	log.Printf("gRPC server started on port %s", app.GpcPort)

	if err = server.Serve(listen); err != nil {
		log.Fatalf("Failed to listen gRPC: %v", err)
	}
}
