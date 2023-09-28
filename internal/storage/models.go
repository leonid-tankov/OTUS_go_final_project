package storage

type Click struct {
	BannerID         int64
	SlotID           int64
	SocialDemGroupID int64
}

type Counter struct {
	BannerID int64
	Count    int
}
