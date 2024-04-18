package tests

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
	"github.com/victorsteven/fullstack/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateDrug(t *testing.T) {

	err := refreshTable()
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
		inputJSON    string
		statusCode   int
		name         string
		approved     bool
		min_dose     uint32
		max_dose     uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"name":"Aspirina", "approved": true, "min_dose": 1, "max_dose": 1}`,
			statusCode:   201,
			tokenGiven:   tokenString,
			name:         "The title",
			approved:     true,
			min_dose:     1,
			max_dose:     1,
			errorMessage: "",
		},
		{
			inputJSON:    `{"name":"genurin", "approved": false, "min_dose": 1, "max_dose": 10}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Title Already Taken",
		},
		{
			// When no token is passed
			inputJSON:    `{"name":"genurin", "approved": false, "min_dose": 1, "max_dose": 10}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"name":"genurin", "approved": false, "min_dose": 1, "max_dose": 10}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"name":"genurin", "approved": false, "min_dose": 1, "max_dose": 10}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Name",
		},
		{
			inputJSON:    `{"name":"genurin", "approved": false, "min_dose": 1, "max_dose": 10}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Approving",
		},
		{
			inputJSON:    `{"name":"genurin", "approved": false, "min_dose": 1, "max_dose": 10}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Name",
		},
		{
			// When user 2 uses user 1 token
			inputJSON:    `{"name":"genurin", "approved": false, "min_dose": 1, "max_dose": 10}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/drugs", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateDrug)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["name"], v.name)
			assert.Equal(t, responseMap["approved"], v.approved)
			assert.Equal(t, responseMap["min_dose"], float64(v.min_dose))
			assert.Equal(t, responseMap["max_dose"], float64(v.max_dose))
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetDrugs(t *testing.T) {

	err := refreshTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedDataTable()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/drugs", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetPosts)
	handler.ServeHTTP(rr, req)

	var drugs []models.Drug
	err = json.Unmarshal([]byte(rr.Body.String()), &drugs)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(drugs), 2)
}
func TestFindDrugByID(t *testing.T) {

	err := refreshTable()
	if err != nil {
		log.Fatal(err)
	}
	drug, err := seedOneDataTable()
	if err != nil {
		log.Fatal(err)
	}
	drugSample := []struct {
		id           string
		statusCode   int
		name         string
		approved     bool
		min_dose     uint32
		max_dose     uint32
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(drug.ID)),
			statusCode: 200,
			name:       drug.Name,
			approved:   drug.Approved,
			min_dose:   drug.Min_dose,
			max_dose:   drug.Max_dose,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range drugSample {

		req, err := http.NewRequest("GET", "/drugs", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetDrug)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, responseMap["name"], v.name)
			assert.Equal(t, responseMap["approved"], v.approved)
			assert.Equal(t, responseMap["min_dose"], float64(v.min_dose))
			assert.Equal(t, responseMap["max_dose"], float64(v.max_dose))
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestUpdateDrug(t *testing.T) {

	var DrugUserEmail, DrugUserPassword string
	var AuthPostID uint64

	err := refreshTable()
	if err != nil {
		log.Fatal(err)
	}
	users, drugs, err := seedDataTable()
	if err != nil {
		log.Fatal(err)
	}
	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		DrugUserEmail = user.Email
		DrugUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(DrugUserEmail, DrugUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)
	// Get only the first post
	for _, drug := range drugs {
		if drug.ID == 2 {
			continue
		}
		AuthPostID = drug.ID
	}
	// fmt.Printf("this is the auth post: %v\n", AuthPostID)

	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		name         string
		approved     bool
		min_dose     uint32
		max_dose     uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"name":"genurin", "approved": false, "min_dose": 1, "max_dose": 10}`,
			statusCode:   200,
			name:         "The title",
			approved:     true,
			min_dose:     1,
			max_dose:     1,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			// When no token is provided
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"name":"genurin", "approved": false, "min_dose": 1, "max_dose": 10}`,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is provided
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"name":"genurin", "approved": false, "min_dose": 1, "max_dose": 10}`,
			tokenGiven:   "this is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			//Note: "Title 2" belongs to post 2, and title must be unique
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"name":"genurin", "approved": false, "min_dose": 1, "max_dose": 10}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Title Already Taken",
		},
		{
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"name":"genurin", "approved": false, "min_dose": 1, "max_dose": 10}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Title",
		},
		{
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"name":"genurin", "approved": false, "min_dose": 1, "max_dose": 10}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Content",
		},
		{
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"name":"genurin", "approved": false, "min_dose": 1, "max_dose": 10}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Author",
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"name":"genurin", "approved": false, "min_dose": 1, "max_dose": 10}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/drugs", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateDrug)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["name"], v.name)
			assert.Equal(t, responseMap["approved"], v.approved)
			assert.Equal(t, responseMap["min_dose"], float64(v.min_dose))
			assert.Equal(t, responseMap["max_dose"], float64(v.max_dose))
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteDrug(t *testing.T) {

	var PostUserEmail, PostUserPassword string
	var PostUserID uint32
	var AuthPostID uint64

	err := refreshTable()
	if err != nil {
		log.Fatal(err)
	}
	users, posts, err := seedDataTable()
	if err != nil {
		log.Fatal(err)
	}
	//Let's get only the Second user
	for _, user := range users {
		if user.ID == 1 {
			continue
		}
		PostUserEmail = user.Email
		PostUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(PostUserEmail, PostUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the second post
	for _, post := range posts {
		if post.ID == 1 {
			continue
		}
		AuthPostID = post.ID
		PostUserID = post.Min_dose
	}
	postSample := []struct {
		id           string
		author_id    uint32
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthPostID)),
			author_id:    PostUserID,
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When empty token is passed
			id:           strconv.Itoa(int(AuthPostID)),
			author_id:    PostUserID,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			id:           strconv.Itoa(int(AuthPostID)),
			author_id:    PostUserID,
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

	for _, v := range postSample {

		req, _ := http.NewRequest("GET", "/posts", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteDrug)

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
