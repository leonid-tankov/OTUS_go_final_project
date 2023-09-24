package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/leonid-tankov/OTUS_go_final_project/internal/app"
	"github.com/leonid-tankov/OTUS_go_final_project/internal/server/grpc/pb"
	"github.com/leonid-tankov/OTUS_go_final_project/internal/storage/migrations"
	"github.com/pressly/goose/v3"
	"google.golang.org/grpc"
)

type BannerRotationServer struct {
	pb.UnimplementedBannerRotationServer
	app    *app.App
	server *grpc.Server
}

func NewServer(app *app.App) *BannerRotationServer {
	return &BannerRotationServer{
		app:    app,
		server: grpc.NewServer(),
	}
}

func (brs *BannerRotationServer) AddBanner(ctx context.Context, req *pb.BannerRequest) (*empty.Empty, error) {
	if err := brs.app.AddBanner(ctx, req.GetBannerId(), req.GetSlotId()); err != nil {
		brs.app.Logger.Error(err)
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (brs *BannerRotationServer) DeleteBanner(ctx context.Context, req *pb.BannerRequest) (*empty.Empty, error) {
	if err := brs.app.DeleteBanner(ctx, req.GetBannerId(), req.GetSlotId()); err != nil {
		brs.app.Logger.Error(err)
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (brs *BannerRotationServer) ClickBanner(ctx context.Context, req *pb.BannerRequest) (*empty.Empty, error) {
	if err := brs.app.ClickBanner(ctx, req.GetBannerId(), req.GetSlotId(), req.GetSocialDemGroupId()); err != nil {
		brs.app.Logger.Error(err)
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (brs *BannerRotationServer) GetBanner(ctx context.Context, req *pb.BannerRequest) (*pb.BannerResponse, error) {
	bannerID, err := brs.app.GetBanner(ctx, req.GetSlotId(), req.GetSocialDemGroupId())
	if err != nil {
		brs.app.Logger.Error(err)
		return nil, err
	}
	return &pb.BannerResponse{
		BannerId: bannerID,
	}, nil
}

func (brs *BannerRotationServer) Start(address string) error {
	lsn, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	pb.RegisterBannerRotationServer(brs.server, brs)

	brs.app.Logger.Infof("starting server on %s", lsn.Addr().String())
	return brs.server.Serve(lsn)
}

func (brs *BannerRotationServer) Stop() {
	brs.app.Logger.Info("Stopping server...")
	brs.server.GracefulStop()
}

func (brs *BannerRotationServer) Run() {
	err := brs.migrate()
	if err != nil {
		brs.app.Logger.Fatal(err)
	}
	err = brs.app.Producer.CreateTopic("banner-rotation-analytics")
	if err != nil {
		brs.app.Logger.Fatal(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()
		brs.Stop()
	}()

	if err := brs.Start(fmt.Sprintf(":%s", brs.app.Config.Grpc.Port)); err != nil {
		cancel()
		brs.app.Logger.Fatal("failed to start grpc server: " + err.Error())
	}
}

func (brs *BannerRotationServer) migrate() error {
	stdlib.GetDefaultDriver()

	db, err := goose.OpenDBWithDriver("pgx", brs.app.Storage.Dsn)
	if err != nil {
		return err
	}
	goose.SetLogger(brs.app.Logger)
	goose.SetBaseFS(&migrations.EmbedMigrations)

	err = goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	brs.app.Logger.Info("starting migrations...")
	err = goose.Up(db, ".")
	if err != nil {
		return err
	}

	err = db.Close()
	if err != nil {
		return err
	}
	brs.app.Logger.Info("end migrations...")

	return nil
}
