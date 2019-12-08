package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/diegoclair/ApiGolang/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateSale(t *testing.T) {

	err := refreshUserBuyAndSaleTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}
	token, err := server.SignIn(user.Email, "password") //Note the password in the database is already hashed, we want unhashed
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		inputJSON    	string
		statusCode   	int
		bitcoinAmount	string
		author_id    	uint32
		tokenGiven   	string
		errorMessage 	string
	}{
		{
			inputJSON:    `{"bitcoinAmount":"0.0050", "author_id": 1}`,
			statusCode:   	201,
			tokenGiven:   	tokenString,
			bitcoinAmount:  "0.0050",
			author_id:    	user.ID,
			errorMessage: 	"",
		},
		{
			inputJSON:    `{"bitcoinAmount":"0.0050", "author_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Title Already Taken",
		},
		{
			// When no token is passed
			inputJSON:    `{"bitcoinAmount":"When no token is passed", "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"bitcoinAmount":"When incorrect token is passed", "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"bitcoinAmount": "", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Title",
		},
		{
			inputJSON:    `{"bitcoinAmount":"0.0050"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Author",
		},
		{
			// When user 2 uses user 1 token
			inputJSON:    `{"bitcoinAmount": "0.0050", "author_id": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/sales", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateSale)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["bitcoinAmount"], v.bitcoinAmount)
			assert.Equal(t, responseMap["author_id"], float64(v.author_id)) //just for both ids to have the same type
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetSales(t *testing.T) {

	err := refreshUserBuyAndSaleTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, _, err = seedUsersBuysAndSales()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/sales", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetSales)
	handler.ServeHTTP(rr, req)

	var sales []models.Sale
	err = json.Unmarshal([]byte(rr.Body.String()), &sales)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(sales), 2)
}
func TestGetSaleByID(t *testing.T) {

	err := refreshUserBuyAndSaleTable()
	if err != nil {
		log.Fatal(err)
	}
	_, sale, err := seedOneUserOneBuyAndOneSale()
	if err != nil {
		log.Fatal(err)
	}
	saleSample := []struct {
		id           	string
		statusCode   	int
		bitcoinAmount	string
		author_id   	uint32
		errorMessage 	string
	}{
		{
			id:         strconv.Itoa(int(sale.ID)),
			statusCode: 	200,
			bitcoinAmount:	sale.BitcoinAmount,
			author_id:  	sale.AuthorID,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range saleSample {

		req, err := http.NewRequest("GET", "/sales", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetSales)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, sale.BitcoinAmount, responseMap["bitcoinAmount"])
			assert.Equal(t, float64(sale.AuthorID), responseMap["author_id"]) //the response author id is float64
		}
	}
}

func TestUpdateSale(t *testing.T) {

	var SaleUserEmail, SaleUserPassword string
	var AuthSaleAuthorID uint32
	var AuthSaleID uint64

	err := refreshUserBuyAndSaleTable()
	if err != nil {
		log.Fatal(err)
	}
	users, _, sales, err := seedUsersBuysAndSales()
	if err != nil {
		log.Fatal(err)
	}
	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		SaleUserEmail = user.Email
		SaleUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(SaleUserEmail, SaleUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the first sale
	for _, sale := range sales {
		if sale.ID == 2 {
			continue
		}
		AuthSaleID = sale.ID
		AuthSaleAuthorID = sale.AuthorID
	}
	// fmt.Printf("this is the auth sale: %v\n", AuthSaleID)

	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		bitcoinAmount        string
		content      string
		author_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthSaleID)),
			updateJSON:   `{"bitcoinAmount":"The updated sale", "content": "This is the updated content", "author_id": 1}`,
			statusCode:   200,
			bitcoinAmount:        "The updated sale",
			content:      "This is the updated content",
			author_id:    AuthSaleAuthorID,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			// When no token is provided
			id:           strconv.Itoa(int(AuthSaleID)),
			updateJSON:   `{"bitcoinAmount":"This is still another bitcoinAmount", "content": "This is the updated content", "author_id": 1}`,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is provided
			id:           strconv.Itoa(int(AuthSaleID)),
			updateJSON:   `{"bitcoinAmount":"This is still another bitcoinAmount", "author_id": 1}`,
			tokenGiven:   "this is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			//Note: "Title 2" belongs to sale 2, and bitcoinAmount must be unique
			id:           strconv.Itoa(int(AuthSaleID)),
			updateJSON:   `{"bitcoinAmount":"Title 2", "author_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Title Already Taken",
		},
		{
			id:           strconv.Itoa(int(AuthSaleID)),
			updateJSON:   `{"bitcoinAmount":"", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Bitcoin Amount",
		},
		{
			id:           strconv.Itoa(int(AuthSaleID)),
			updateJSON:   `{"bitcoinAmount":"This is another bitcoinAmount"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(AuthSaleID)),
			updateJSON:   `{"bitcoinAmount":"This is still another bitcoinAmount", "author_id": 2}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/sales", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateSale)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["bitcoinAmount"], v.bitcoinAmount)
			assert.Equal(t, responseMap["author_id"], float64(v.author_id)) //just to match the type of the json we receive thats why we used float64
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteSale(t *testing.T) {

	var SaleUserEmail, SaleUserPassword string
	var SaleUserID uint32
	var AuthSaleID uint64

	err := refreshUserBuyAndSaleTable()
	if err != nil {
		log.Fatal(err)
	}
	users, _, sales, err := seedUsersBuysAndSales()
	if err != nil {
		log.Fatal(err)
	}
	//Let's get only the Second user
	for _, user := range users {
		if user.ID == 1 {
			continue
		}
		SaleUserEmail = user.Email
		SaleUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(SaleUserEmail, SaleUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the second buy
	for _, sale := range sales {
		if sale.ID == 1 {
			continue
		}
		AuthSaleID = sale.ID
		SaleUserID = sale.AuthorID
	}
	saleSample := []struct {
		id           string
		author_id    uint32
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthSaleID)),
			author_id:    SaleUserID,
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When empty token is passed
			id:           strconv.Itoa(int(AuthSaleID)),
			author_id:    SaleUserID,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			id:           strconv.Itoa(int(AuthSaleID)),
			author_id:    SaleUserID,
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(1)),
			author_id:    1,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range saleSample {

		req, _ := http.NewRequest("GET", "/sales", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteSale)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {

			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}