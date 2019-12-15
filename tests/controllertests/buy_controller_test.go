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

	"github.com/diegoclair/ApiGolang/api/models"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateBuy(t *testing.T) {

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
		inputJSON      string
		statusCode     int
		bitcoin_amount float64
		author_id      uint32
		tokenGiven     string
		errorMessage   string
	}{
		{
			inputJSON:      `{"bitcoin_amount":0.0050, "author_id": 1}`,
			statusCode:     201,
			tokenGiven:     tokenString,
			bitcoin_amount: 0.0050,
			author_id:      user.ID,
			errorMessage:   "",
		},
		{
			// When no token is passed
			inputJSON:    `{"bitcoin_amount":0.00801, "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"bitcoin_amount":0.0054747, "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"bitcoin_amount": 0, "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Bitcoin Amount",
		},
		{
			inputJSON:    `{"bitcoin_amount":0.0050}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Author",
		},
		{
			// When user 2 uses user 1 token
			inputJSON:    `{"bitcoin_amount": 0.0050, "author_id": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/buys", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateBuy)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["bitcoin_amount"], v.bitcoin_amount)
			assert.Equal(t, responseMap["author_id"], float64(v.author_id)) //just for both ids to have the same type
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetBuys(t *testing.T) {

	err := refreshUserBuyAndSaleTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, _, err = seedUsersBuysAndSales()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/buys", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetBuys)
	handler.ServeHTTP(rr, req)

	var buys []models.Buy
	err = json.Unmarshal([]byte(rr.Body.String()), &buys)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(buys), 2)
}
func TestGetBuyByID(t *testing.T) {

	err := refreshUserBuyAndSaleTable()
	if err != nil {
		log.Fatal(err)
	}
	buy, _, err := seedOneUserOneBuyAndOneSale()
	if err != nil {
		log.Fatal(err)
	}
	buySample := []struct {
		id             string
		statusCode     int
		bitcoin_amount float64
		author_id      uint32
		errorMessage   string
	}{
		{
			id:             strconv.Itoa(int(buy.ID)),
			statusCode:     200,
			bitcoin_amount: buy.BitcoinAmount,
			author_id:      buy.AuthorID,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}

	for _, v := range buySample {
		req, err := http.NewRequest("GET", "/buys", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetBuys)
		handler.ServeHTTP(rr, req)
		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)

		if err != nil {
			log.Fatalf("Cannot convert to json4: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, buy.BitcoinAmount, responseMap["bitcoin_amount"])
			assert.Equal(t, float64(buy.AuthorID), responseMap["author_id"]) //the response author id is float64
		}
	}
}

func TestUpdateBuy(t *testing.T) {

	var BuyUserEmail, BuyUserPassword string
	var AuthBuyAuthorID uint32
	var AuthBuyID uint64

	err := refreshUserBuyAndSaleTable()
	if err != nil {
		log.Fatal(err)
	}
	users, buys, _, err := seedUsersBuysAndSales()
	if err != nil {
		log.Fatal(err)
	}
	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		BuyUserEmail = user.Email
		BuyUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(BuyUserEmail, BuyUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the first buy
	for _, buy := range buys {
		if buy.ID == 2 {
			continue
		}
		AuthBuyID = buy.ID
		AuthBuyAuthorID = buy.AuthorID
	}
	// fmt.Printf("this is the auth buy: %v\n", AuthBuyID)

	samples := []struct {
		id             string
		updateJSON     string
		statusCode     int
		bitcoin_amount string
		content        string
		author_id      uint32
		tokenGiven     string
		errorMessage   string
	}{
		{
			// Convert int64 to int first before converting to string
			id:             strconv.Itoa(int(AuthBuyID)),
			updateJSON:     `{"bitcoin_amount":"The updated buy", "content": "This is the updated content", "author_id": 1}`,
			statusCode:     200,
			bitcoin_amount: "The updated buy",
			content:        "This is the updated content",
			author_id:      AuthBuyAuthorID,
			tokenGiven:     tokenString,
			errorMessage:   "",
		},
		{
			// When no token is provided
			id:           strconv.Itoa(int(AuthBuyID)),
			updateJSON:   `{"bitcoin_amount":"This is still another bitcoin_amount", "content": "This is the updated content", "author_id": 1}`,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is provided
			id:           strconv.Itoa(int(AuthBuyID)),
			updateJSON:   `{"bitcoin_amount":"This is still another bitcoin_amount", "author_id": 1}`,
			tokenGiven:   "this is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			//Note: "Title 2" belongs to buy 2, and bitcoin_amount must be unique
			id:           strconv.Itoa(int(AuthBuyID)),
			updateJSON:   `{"bitcoin_amount":"Title 2", "author_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Title Already Taken",
		},
		{
			id:           strconv.Itoa(int(AuthBuyID)),
			updateJSON:   `{"bitcoin_amount":"", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Bitcoin Amount",
		},
		{
			id:           strconv.Itoa(int(AuthBuyID)),
			updateJSON:   `{"bitcoin_amount":"This is another bitcoin_amount"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(AuthBuyID)),
			updateJSON:   `{"bitcoin_amount":"This is still another bitcoin_amount", "author_id": 2}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/buys", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateBuy)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["bitcoin_amount"], v.bitcoin_amount)
			assert.Equal(t, responseMap["author_id"], float64(v.author_id)) //just to match the type of the json we receive thats why we used float64
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteBuy(t *testing.T) {

	var BuyUserEmail, BuyUserPassword string
	var BuyUserID uint32
	var AuthBuyID uint64

	err := refreshUserBuyAndSaleTable()
	if err != nil {
		log.Fatal(err)
	}
	users, buys, _, err := seedUsersBuysAndSales()
	if err != nil {
		log.Fatal(err)
	}
	//Let's get only the Second user
	for _, user := range users {
		if user.ID == 1 {
			continue
		}
		BuyUserEmail = user.Email
		BuyUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(BuyUserEmail, BuyUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the second buy
	for _, buy := range buys {
		if buy.ID == 1 {
			continue
		}
		AuthBuyID = buy.ID
		BuyUserID = buy.AuthorID
	}
	buySample := []struct {
		id           string
		author_id    uint32
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthBuyID)),
			author_id:    BuyUserID,
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When empty token is passed
			id:           strconv.Itoa(int(AuthBuyID)),
			author_id:    BuyUserID,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			id:           strconv.Itoa(int(AuthBuyID)),
			author_id:    BuyUserID,
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
	for _, v := range buySample {

		req, _ := http.NewRequest("GET", "/buys", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteBuy)

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
