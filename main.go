package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/konrad2002/tmate-server/attest"
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

	ms  service.MemberService
	fs  service.FieldService
	cs  service.ConfigService
	qs  service.QueryService
	as  service.AttestService
	es  service.ExportService
	us  service.UserService
	hs  service.HistoryService
	ems service.EmailService

	mc  controller.MemberController
	fc  controller.FieldController
	cc  controller.ConfigController
	qc  controller.QueryController
	ec  controller.ExportController
	uc  controller.UserController
	emc controller.EmailController
	ac  controller.AttestController

	ctx         context.Context
	mongoClient *mongo.Client
)

func init() {
	println("initializing...")
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
	ur := repository.NewUserRepository(mongoCon)
	hr := repository.NewHistoryRepository(mongoCon)

	hs = service.NewHistoryService(hr)
	fs = service.NewFieldService(fr)
	qs = service.NewQueryService(qr, hs)
	cs = service.NewConfigService()
	ms = service.NewMemberService(mr, qs, fs, cs, hs)
	es = service.NewExportService(ms)
	us = service.NewUserService(ur)
	ems = service.NewEmailService(cs, ms, hs)
	as = service.NewAttestService(ms, fs, cs, ems)

	mc = controller.NewMemberController(ms, us)
	fc = controller.NewFieldController(fs, us)
	cc = controller.NewConfigController(cs, us)
	qc = controller.NewQueryController(qs, us)
	ec = controller.NewExportController(es, us)
	uc = controller.NewUserController(us)
	emc = controller.NewEmailController(ems, us)
	ac = controller.NewAttestController(as, us)

	server = gin.Default()
}

func main() {
	println("starting...")
	defer func(mongoClient *mongo.Client, ctx context.Context) {
		err := mongoClient.Disconnect(ctx)
		if err != nil {
			panic(err)
		}
	}(mongoClient, ctx)

	server.Use(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Status(204)
			return
		}
		c.Next()
	})

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://tmate.weiss-konrad.de", "http://localhost:4200", "https://tmate.st-erz.de"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	basePath := server.Group("/api/v1")

	mc.RegisterRoutes(basePath)
	fc.RegisterRoutes(basePath)
	cc.RegisterRoutes(basePath)
	qc.RegisterRoutes(basePath)
	ec.RegisterRoutes(basePath)
	uc.RegisterRoutes(basePath)
	emc.RegisterRoutes(basePath)
	ac.RegisterRoutes(basePath)

	port := os.Getenv("TMATE_PORT")

	if port == "" {
		fmt.Println("no application port given! Please set TMATE_PORT.")
		return
	}

	attest.StartAttestRoutine(as)

	log.Fatal(server.Run(":" + port))
}
