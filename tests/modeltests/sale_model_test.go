package modeltests

import (
	"log"
	"testing"

	"github.com/diegoclair/ApiGolang/api/models"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllSales(t *testing.T) {

	err := refreshUserBuyAndSaleTable()
	if err != nil {
		log.Fatalf("Error refreshing user and sale table %v\n", err)
	}
	_, _, _, err = seedUsersBuysAndSales()
	if err != nil {
		log.Fatalf("Error seeding user and sale  table %v\n", err)
	}
	sales, err := saleInstance.FindAllSales(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the sales: %v\n", err)
		return
	}
	assert.Equal(t, len(*sales), 2)
}

func TestSaveSale(t *testing.T) {

	err := refreshUserBuyAndSaleTable()
	if err != nil {
		log.Fatalf("Error user and sale refreshing table %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	newSale := models.Sale{
		ID:            1,
		BitcoinAmount: "0.00254",
		AuthorID:      user.ID,
	}
	savedSale, err := newSale.SaveSale(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the sale: %v\n", err)
		return
	}
	assert.Equal(t, newSale.ID, savedSale.ID)
	assert.Equal(t, newSale.BitcoinAmount, savedSale.BitcoinAmount)
	assert.Equal(t, newSale.AuthorID, savedSale.AuthorID)

}

func TestGetSaleByID(t *testing.T) {

	err := refreshUserBuyAndSaleTable()
	if err != nil {
		log.Fatalf("Error refreshing user and sale table: %v\n", err)
	}
	_, sale, err := seedOneUserOneBuyAndOneSale()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	foundSale, err := saleInstance.FindSaleByID(server.DB, sale.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundSale.ID, sale.ID)
	assert.Equal(t, foundSale.BitcoinAmount, sale.BitcoinAmount)
}

func TestUpdateASale(t *testing.T) {

	err := refreshUserBuyAndSaleTable()
	if err != nil {
		log.Fatalf("Error refreshing user and sale table: %v\n", err)
	}
	_, sale, err := seedOneUserOneBuyAndOneSale()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	saleUpdate := models.Sale{
		ID:            1,
		BitcoinAmount: "0.00170004897",
		AuthorID:      sale.AuthorID,
	}
	updatedSale, err := saleUpdate.UpdateASale(server.DB)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, updatedSale.ID, saleUpdate.ID)
	assert.Equal(t, updatedSale.BitcoinAmount, saleUpdate.BitcoinAmount)
	assert.Equal(t, updatedSale.AuthorID, saleUpdate.AuthorID)
}

func TestDeleteASale(t *testing.T) {

	err := refreshUserBuyAndSaleTable()
	if err != nil {
		log.Fatalf("Error refreshing user and sale table: %v\n", err)
	}
	_, sale, err := seedOneUserOneBuyAndOneSale()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	isDeleted, err := saleInstance.DeleteASale(server.DB, sale.ID, sale.AuthorID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}
