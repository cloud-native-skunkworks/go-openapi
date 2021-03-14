package restapi

import (
	"bytes"
	"encoding/json"
	"github.com/AlexsJones/go-openapi/models"
	"github.com/AlexsJones/go-openapi/restapi/operations"
	"github.com/go-openapi/loads"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getAPI() (*operations.GoOpenapiAPI, error) {

	swaggerSpec, err := loads.Embedded(SwaggerJSON, FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewGoOpenapiAPI(swaggerSpec)
	return api, nil
}

func GetAPIHandler() (http.Handler, error) {
	api, err := getAPI()
	if err != nil {
		return nil, err
	}
	h := configureAPI(api)
	err = api.Validate()
	if err != nil {
		return nil, err
	}
	return h, nil
}

func TestCreateUser(t *testing.T){
	handler, err := GetAPIHandler()
	if err != nil {
		t.Fatal("get api handler", err)
	}
	ts := httptest.NewServer(handler)
	defer ts.Close()

	user := models.User{
		Email:      "foo@bar.com",
		FirstName:  "alex",
		LastName:   "jones",
		Password:   "test",
		Phone:      "555-555",

		Username:   "alex",
	}
	jsonValue, _ := json.Marshal(user)
	resp, err := http.Post(ts.URL + "/v2/user", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatal()
	}
	if resp.StatusCode != 501 {
		t.Fatal()
	}
}