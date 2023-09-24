package app

import (
	"context"
	"time"

	"github.com/leonid-tankov/OTUS_go_final_project/internal/algorithms"
	"github.com/leonid-tankov/OTUS_go_final_project/internal/config"
	"github.com/leonid-tankov/OTUS_go_final_project/internal/kafka"
	"github.com/leonid-tankov/OTUS_go_final_project/internal/storage"
	"github.com/sirupsen/logrus"
)

type App struct {
	Config   *config.Config
	Logger   *logrus.Logger
	Storage  *storage.Storage
	Producer *kafka.Producer
}

func New(cfg *config.Config, logger *logrus.Logger, storage *storage.Storage, producer *kafka.Producer) *App {
	return &App{
		Config:   cfg,
		Logger:   logger,
		Storage:  storage,
		Producer: producer,
	}
}

func (a *App) AddBanner(ctx context.Context, bannerID, slotID int64) error {
	if err := a.Storage.Connect(ctx); err != nil {
		return err
	}
	defer a.Storage.Close(ctx)
	if err := a.Storage.CheckBanner(ctx, bannerID); err != nil {
		return err
	}
	if err := a.Storage.CheckSlot(ctx, slotID); err != nil {
		return err
	}
	if err := a.Storage.CheckClicks(ctx, storage.Click{
		BannerID: bannerID,
		SlotID:   slotID,
	}); err != nil {
		return err
	}
	IDs, err := a.Storage.GetSocialDemGroups(ctx)
	if err != nil {
		return err
	}
	for _, ID := range IDs {
		if err := a.Storage.InsertClicks(ctx, storage.Click{
			BannerID:         bannerID,
			SlotID:           slotID,
			SocialDemGroupID: int64(ID),
		}); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) DeleteBanner(ctx context.Context, bannerID, slotID int64) error {
	if err := a.Storage.Connect(ctx); err != nil {
		return err
	}
	defer a.Storage.Close(ctx)
	if err := a.Storage.DeleteClicks(ctx, storage.Click{
		BannerID: bannerID,
		SlotID:   slotID,
	}); err != nil {
		return err
	}
	return nil
}

func (a *App) ClickBanner(ctx context.Context, bannerID, slotID, socialDemGroupID int64) error {
	if err := a.Storage.Connect(ctx); err != nil {
		return err
	}
	defer a.Storage.Close(ctx)
	if err := a.Storage.UpdateClicks(ctx, storage.Click{
		BannerID:         bannerID,
		SlotID:           slotID,
		SocialDemGroupID: socialDemGroupID,
	}); err != nil {
		return err
	}
	if err := a.Producer.Produce(ctx, kafka.Message{
		EventType:        kafka.Click,
		BannerID:         bannerID,
		SlotID:           slotID,
		SocialDemGroupID: socialDemGroupID,
		Timestamp:        time.Now(),
	}); err != nil {
		return err
	}
	return nil
}

func (a *App) GetBanner(ctx context.Context, slotID, socialDemGroupID int64) (int64, error) {
	if err := a.Storage.Connect(ctx); err != nil {
		return 0, err
	}
	defer a.Storage.Close(ctx)
	counters, err := a.Storage.GetClicks(ctx, storage.Click{
		SlotID:           slotID,
		SocialDemGroupID: socialDemGroupID,
	})
	if err != nil {
		return 0, err
	}
	ID := algorithms.MultiArmedBandit(counters, a.Config.Probability)
	if err = a.Producer.Produce(ctx, kafka.Message{
		EventType:        kafka.Show,
		BannerID:         ID,
		SlotID:           slotID,
		SocialDemGroupID: socialDemGroupID,
		Timestamp:        time.Now(),
	}); err != nil {
		return ID, err
	}
	return ID, nil
}
