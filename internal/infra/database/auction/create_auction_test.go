package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"testing"
	"time"
)

func TestCreateAuction_ShouldCloseAuctionAutomatically(t *testing.T) {
	repo := &AuctionRepository{
		auctionInterval: 20 * time.Millisecond,
	}

	repo.insertAuctionFn = func(ctx context.Context, auction *AuctionEntityMongo) error {
		return nil
	}

	closeCalled := make(chan struct{}, 1)
	repo.closeAuctionIfExpiredByIDFn = func(ctx context.Context, auctionId string, auctionEndTime time.Time) error {
		closeCalled <- struct{}{}
		return nil
	}

	err := repo.CreateAuction(context.Background(), &auction_entity.Auction{
		Id:          "auction-id-1",
		ProductName: "Produto Teste",
		Category:    "Eletronicos",
		Description: "Descricao de teste valida para leilao",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now(),
	})
	if err != nil {
		t.Fatalf("expected no error creating auction, got: %v", err)
	}

	select {
	case <-closeCalled:
	case <-time.After(250 * time.Millisecond):
		t.Fatalf("expected automatic close routine to be called")
	}
}
