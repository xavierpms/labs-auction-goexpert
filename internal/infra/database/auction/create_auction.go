package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection                  *mongo.Collection
	auctionInterval             time.Duration
	insertAuctionFn             func(ctx context.Context, auction *AuctionEntityMongo) error
	closeAuctionIfExpiredByIDFn func(ctx context.Context, auctionId string, auctionEndTime time.Time) error
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	repository := &AuctionRepository{
		Collection:      database.Collection("auctions"),
		auctionInterval: getAuctionInterval(),
	}

	repository.insertAuctionFn = repository.insertAuction
	repository.closeAuctionIfExpiredByIDFn = repository.closeAuctionIfExpiredByID

	return repository
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}

	err := ar.insertAuctionFn(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	go ar.autoCloseAuction(auctionEntity.Id, auctionEntity.Timestamp)

	return nil
}

func (ar *AuctionRepository) insertAuction(ctx context.Context, auction *AuctionEntityMongo) error {
	_, err := ar.Collection.InsertOne(ctx, auction)
	return err
}

func (ar *AuctionRepository) autoCloseAuction(auctionId string, auctionTimestamp time.Time) {
	auctionEndTime := calculateAuctionEndTime(auctionTimestamp, ar.auctionInterval)
	waitTime := time.Until(auctionEndTime)
	if waitTime > 0 {
		time.Sleep(waitTime)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := ar.closeAuctionIfExpiredByIDFn(ctx, auctionId, auctionEndTime); err != nil {
		logger.Error("Error trying to close auction automatically", err)
	}
}

func (ar *AuctionRepository) closeAuctionIfExpiredByID(ctx context.Context, auctionId string, auctionEndTime time.Time) error {
	filter := bson.M{
		"_id":       auctionId,
		"status":    auction_entity.Active,
		"timestamp": bson.M{"$lte": auctionEndTime.Unix()},
	}

	update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}

	_, err := ar.Collection.UpdateOne(ctx, filter, update)
	return err
}

func getAuctionInterval() time.Duration {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return time.Minute * 5
	}

	return duration
}

func calculateAuctionEndTime(auctionTimestamp time.Time, auctionInterval time.Duration) time.Time {
	return auctionTimestamp.Add(auctionInterval)
}
