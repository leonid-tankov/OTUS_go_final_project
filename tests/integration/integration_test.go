//go:build integration
// +build integration

package integration_test

import (
	"context"
	"testing"

	"github.com/leonid-tankov/OTUS_go_final_project/internal/server/grpc/pb"
	"github.com/leonid-tankov/OTUS_go_final_project/internal/storage"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AppSuite struct {
	suite.Suite
	ctx       context.Context
	appClient pb.BannerRotationClient
}

func (a *AppSuite) SetupSuite() {
	statHost := "127.0.0.1:8000"

	conn, err := grpc.Dial(statHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	a.Require().NoError(err)

	a.ctx = context.Background()
	a.appClient = pb.NewBannerRotationClient(conn)
}

func (a *AppSuite) SetupTest() {}

func (a *AppSuite) TestAddBanner() {
	a.addBanner(1, 1)

	_, err := a.appClient.AddBanner(a.ctx, &pb.BannerRequest{SlotId: 0, BannerId: 1}) // индекса 0 нет
	a.Require().Error(err)

	_, err = a.appClient.AddBanner(a.ctx, &pb.BannerRequest{SlotId: 1, BannerId: 0}) // индекса 0 нет
	a.Require().Error(err)

	_, err = a.appClient.AddBanner(a.ctx, &pb.BannerRequest{SlotId: 1, BannerId: 1})
	a.Require().ErrorAs(err, &storage.ErrHasRotation)

	a.deleteBanner(1, 1)
}

func (a *AppSuite) TestDeleteBanner() {
	_, err := a.appClient.DeleteBanner(a.ctx, &pb.BannerRequest{SlotId: 1, BannerId: 1})
	a.Require().ErrorAs(err, &storage.ErrNoRowsAffected)
}

func (a *AppSuite) TestClickBanner() {
	_, err := a.appClient.ClickBanner(a.ctx, &pb.BannerRequest{SlotId: 1, BannerId: 1, SocialDemGroupId: 1})
	a.Require().ErrorAs(err, &storage.ErrNoRowsAffected)
}

func (a *AppSuite) addBanner(bannerID, slotID int64) {
	req := pb.BannerRequest{
		BannerId: bannerID,
		SlotId:   slotID,
	}
	_, err := a.appClient.AddBanner(a.ctx, &req)
	a.Require().NoError(err)
}

func (a *AppSuite) deleteBanner(bannerID, slotID int64) {
	req := pb.BannerRequest{
		BannerId: bannerID,
		SlotId:   slotID,
	}
	_, err := a.appClient.DeleteBanner(a.ctx, &req)
	a.Require().NoError(err)
}

// func (a *AppSuite) clickBanner(bannerID, slotID, socialDemGroupID int64) {
//	req := pb.BannerRequest{
//		BannerId:         bannerID,
//		SlotId:           slotID,
//		SocialDemGroupId: socialDemGroupID,
//	}
//	_, err := a.appClient.ClickBanner(a.ctx, &req)
//	a.Require().NoError(err)
//}

func TestAppSuite(t *testing.T) {
	suite.Run(t, new(AppSuite))
}
