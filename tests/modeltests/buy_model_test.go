package modeltests

import (
	"log"
	"testing"

	"github.com/diegoclair/ApiGolang/api/models"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllBuys(t *testing.T) {

	err := refreshUserBuyAndSaleTable()
	if err != nil {
		log.Fatalf("Error refreshing user and buy table %v\n", err)
	}
	_, _, _, err = seedUsersBuysAndSales()
	if err != nil {
		log.Fatalf("Error seeding user and buy  table %v\n", err)
	}
	buys, err := buyInstance.FindAllBuys(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the buys: %v\n", err)
		return
	}
	assert.Equal(t, len(*buys), 2)
}

func TestSaveBuy(t *testing.T) {

	err := refreshUserBuyAndSaleTable()
	if err != nil {
		log.Fatalf("Error user and buy refreshing table %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	newBuy := models.Buy{
		ID:            1,
		BitcoinAmount: "0.00254",
		AuthorID:      user.ID,
	}
	savedBuy, err := newBuy.SaveBuy(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the buy: %v\n", err)
		return
	}
	assert.Equal(t, newBuy.ID, savedBuy.ID)
	assert.Equal(t, newBuy.BitcoinAmount, savedBuy.BitcoinAmount)
	assert.Equal(t, newBuy.AuthorID, savedBuy.AuthorID)

}

func TestGetBuyByID(t *testing.T) {

	err := refreshUserBuyAndSaleTable()
	if err != nil {
		log.Fatalf("Error refreshing user and buy table: %v\n", err)
	}
	buy, _, err := seedOneUserOneBuyAndOneSale()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	foundBuy, err := buyInstance.FindBuyByID(server.DB, buy.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundBuy.ID, buy.ID)
	assert.Equal(t, foundBuy.BitcoinAmount, buy.BitcoinAmount)
}

func TestUpdateABuy(t *testing.T) {

	err := refreshUserBuyAndSaleTable()
	if err != nil {
		log.Fatalf("Error refreshing user and buy table: %v\n", err)
	}
	buy, _, err := seedOneUserOneBuyAndOneSale()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	buyUpdate := models.Buy{
		ID:            1,
		BitcoinAmount: "0.00170004897",
		AuthorID:      buy.AuthorID,
	}
	updatedBuy, err := buyUpdate.UpdateABuy(server.DB)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, updatedBuy.ID, buyUpdate.ID)
	assert.Equal(t, updatedBuy.BitcoinAmount, buyUpdate.BitcoinAmount)
	assert.Equal(t, updatedBuy.AuthorID, buyUpdate.AuthorID)
}

func TestDeleteABuy(t *testing.T) {

	err := refreshUserBuyAndSaleTable()
	if err != nil {
		log.Fatalf("Error refreshing user and buy table: %v\n", err)
	}
	buy, _, err := seedOneUserOneBuyAndOneSale()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	isDeleted, err := buyInstance.DeleteABuy(server.DB, buy.ID, buy.AuthorID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}
