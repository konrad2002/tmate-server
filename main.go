package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/konrad2002/tmate-server/controller"
	"github.com/konrad2002/tmate-server/db"
	"github.com/konrad2002/tmate-server/repository"
	"github.com/konrad2002/tmate-server/service"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"time"
)

var (
	server *gin.Engine

	ms service.MemberService
	fs service.FieldService
	cs service.ConfigService
	qs service.QueryService

	mc controller.MemberController
	fc controller.FieldController
	cc controller.ConfigController
	qc controller.QueryController

	ctx         context.Context
	mongoClient *mongo.Client
)

func init() {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoCon, err := db.Connect()
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("mongo connection established")

	mr := repository.NewMemberRepository(mongoCon)
	fr := repository.NewFieldRepository(mongoCon)
	qr := repository.NewQueryRepository(mongoCon)

	fs = service.NewFieldService(fr)
	qs = service.NewQueryService(qr)
	cs = service.NewConfigService()
	ms = service.NewMemberService(mr, qs, fs, cs)

	mc = controller.NewMemberController(ms)
	fc = controller.NewFieldController(fs)
	cc = controller.NewConfigController(cs)
	qc = controller.NewQueryController(qs)

	server = gin.Default()
}

func main() {
	defer func(mongoClient *mongo.Client, ctx context.Context) {
		err := mongoClient.Disconnect(ctx)
		if err != nil {
			panic(err)
		}
	}(mongoClient, ctx)

	basePath := server.Group("/api/v1")

	mc.RegisterRoutes(basePath)
	fc.RegisterRoutes(basePath)
	cc.RegisterRoutes(basePath)
	qc.RegisterRoutes(basePath)

	port := os.Getenv("TMATE_PORT")

	if port == "" {
		fmt.Println("no application port given! Please set TMATE_PORT.")
		return
	}

	log.Fatal(server.Run(":" + port))
}
