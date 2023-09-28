//go:build integration
// +build integration

package integration_test

import (
	"context"
	"testing"
	"time"

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
	statHost := "localhost:8000"

	conn, err := grpc.Dial(statHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	a.Require().NoError(err)

	a.ctx = context.Background()
	a.appClient = pb.NewBannerRotationClient(conn)
}

func (a *AppSuite) SetupTest() {}

func (a *AppSuite) TestAddBanner() {
	a.addBanner(1, 1)

	_, err := a.appClient.AddBanner(a.ctx, &pb.BannerRequest{SlotId: 0, BannerId: 1}) // id начинаются с 1
	a.Require().Error(err)

	_, err = a.appClient.AddBanner(a.ctx, &pb.BannerRequest{SlotId: 1, BannerId: 0}) // id начинаются с 1
	a.Require().Error(err)

	_, err = a.appClient.AddBanner(a.ctx, &pb.BannerRequest{SlotId: 1, BannerId: 1})
	a.Require().ErrorAs(err, &storage.ErrHasRotation)

	time.Sleep(2 * time.Second)

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

func (a *AppSuite) TestGetBanner() {
	time.Sleep(1 * time.Second)
	_, err := a.appClient.GetBanner(a.ctx, &pb.BannerRequest{SlotId: 1, SocialDemGroupId: 1})
	a.Require().ErrorAs(err, &storage.ErrNoRowsAffected)
}

func (a *AppSuite) TestEnumerationOfAll() {
	for i := 1; i < 10; i++ {
		a.addBanner(int64(i), 1)
	}
	for i := 1; i < 10; i++ {
		a.clickBanner(int64(i), 1, 1)
	}
	ids := make([]int64, 9)
	j := 0
	for i := 0; i < 50; i++ {
		bannerId := a.getBanner(1, 1)
		if !a.contain(ids, bannerId) {
			ids[j] = bannerId
			j++
		}
	}
	for i := 1; i < 10; i++ {
		a.deleteBanner(int64(i), 1)
	}
	a.Require().Equal(9, j)
}

func (a *AppSuite) TestSelectionOfPopular() {
	for i := 1; i < 6; i++ {
		a.addBanner(int64(i), 1)
	}
	for i := 1; i < 5; i++ {
		a.clickBanner(int64(i), 1, 1)
	}
	for i := 1; i < 10; i++ {
		a.clickBanner(5, 1, 1)
	}
	counters := make(map[int64]int)
	for i := 0; i < 50; i++ {
		bannerId := a.getBanner(1, 1)
		counters[bannerId]++
	}
	a.Require().GreaterOrEqual(counters[5], 35)
	for i := 1; i < 6; i++ {
		a.deleteBanner(int64(i), 1)
	}
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

func (a *AppSuite) clickBanner(bannerID, slotID, socialDemGroupID int64) {
	req := pb.BannerRequest{
		BannerId:         bannerID,
		SlotId:           slotID,
		SocialDemGroupId: socialDemGroupID,
	}
	_, err := a.appClient.ClickBanner(a.ctx, &req)
	a.Require().NoError(err)
}

func (a *AppSuite) getBanner(slotID, socialDemGroupID int64) int64 {
	req := pb.BannerRequest{
		SlotId:           slotID,
		SocialDemGroupId: socialDemGroupID,
	}
	response, err := a.appClient.GetBanner(a.ctx, &req)
	a.Require().NoError(err)
	return response.BannerId
}

func (a *AppSuite) contain(ids []int64, id int64) bool {
	for _, i := range ids {
		if i == id {
			return true
		}
	}
	return false
}

func TestAppSuite(t *testing.T) {
	suite.Run(t, new(AppSuite))
}
