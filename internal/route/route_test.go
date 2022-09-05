package route

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"main/internal/model"
)

var tokens []string //0 user, 1 admin
var events []string //0 past, 1 future

func TestInit(t *testing.T) {
	router := SetupRouter()
	message := "Welcome to API event-organizer"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)
	messageIn, _ := searchElemInBody(w.Body.String(), "message")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, message, messageIn)
}

func TestLogin(t *testing.T) {
	type caseStruct struct {
		Username string
		Password string
		Status   int
	}
	var cases = []caseStruct{
		{Username: "usuarioNormal", Password: "123456", Status: http.StatusOK},
		{Username: "usuarioAdministrador", Password: "123456", Status: http.StatusOK},
		{Username: "usuarioInexistente", Password: "123456", Status: http.StatusUnauthorized},
	}

	router := SetupRouter()

	for _, testData := range cases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/login", nil)
		req.SetBasicAuth(testData.Username, testData.Password)
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			if token, err := searchElemInBody(w.Body.String(), "token"); err == nil {
				tokens = append(tokens, token)
			}
		}
		t.Log("datos: ", tokens[len(tokens)-1])

		assert.Equal(t, testData.Status, w.Code)
	}

}

func TestGetEvents(t *testing.T) {
	type caseStruct struct {
		UserToken string
		Status    int
	}
	var cases = []caseStruct{
		{UserToken: tokens[0], Status: http.StatusOK},
		{UserToken: tokens[1], Status: http.StatusOK},
	}

	router := SetupRouter()

	for _, testData := range cases {

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/events", nil)
		req.Header.Set("Authorization", "Bearer "+testData.UserToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, testData.Status, w.Code)
	}
}

func TestPostEvent(t *testing.T) {
	pastDate, _ := createJsonEvent("2020-12-30")
	futureDate, _ := createJsonEvent("2040-12-30")
	type caseStruct struct {
		UserToken string
		Status    int
		Json      bytes.Buffer
	}
	var cases = []caseStruct{
		{UserToken: tokens[0], Status: http.StatusUnauthorized, Json: pastDate},
		{UserToken: tokens[1], Status: http.StatusCreated, Json: pastDate},
		{UserToken: tokens[1], Status: http.StatusCreated, Json: futureDate},
		{UserToken: tokens[1], Status: http.StatusCreated, Json: futureDate},
	}

	router := SetupRouter()

	for _, testData := range cases {

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/events", &testData.Json)
		req.Header.Set("Authorization", "Bearer "+testData.UserToken)
		router.ServeHTTP(w, req)

		if w.Code == http.StatusCreated {
			idEvent, _ := searchElemInBody(w.Body.String(), "id")
			if idEvent != "" {
				events = append(events, idEvent)
			}
		}

		assert.Equal(t, testData.Status, w.Code)
	}
}

func TestGetEvent(t *testing.T) {
	type caseStruct struct {
		UserToken string
		Status    int
	}
	var cases = []caseStruct{
		{UserToken: tokens[0], Status: http.StatusOK},
		{UserToken: tokens[1], Status: http.StatusOK},
	}

	router := SetupRouter()

	for _, testData := range cases {

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/events/"+events[len(events)-1], nil)
		req.Header.Set("Authorization", "Bearer "+testData.UserToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, testData.Status, w.Code)
	}
}

func TestPutEvent(t *testing.T) {
	type caseStruct struct {
		UserToken string
		Status    int
	}
	var cases = []caseStruct{
		{UserToken: tokens[0], Status: http.StatusUnauthorized},
		{UserToken: tokens[1], Status: http.StatusOK},
	}

	router := SetupRouter()

	for _, testData := range cases {

		w := httptest.NewRecorder()
		js, _ := createJsonEvent("2022-12-30")
		req, _ := http.NewRequest(http.MethodPut, "/events/"+events[len(events)-1], &js)
		req.Header.Set("Authorization", "Bearer "+testData.UserToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, testData.Status, w.Code)
	}
}

func TestDeleteEvent(t *testing.T) {
	type caseStruct struct {
		UserToken string
		Status    int
	}
	var cases = []caseStruct{
		{UserToken: tokens[0], Status: http.StatusUnauthorized},
		{UserToken: tokens[1], Status: http.StatusOK},
	}

	router := SetupRouter()

	for _, testData := range cases {

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/events/"+events[len(events)-1], nil)
		req.Header.Set("Authorization", "Bearer "+testData.UserToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, testData.Status, w.Code)
	}
}

func TestPostInscription(t *testing.T) {
	past, _ := createJsonInscriptionPost(events[0])
	future, _ := createJsonInscriptionPost(events[1])
	type caseStruct struct {
		UserToken string
		Status    int
		Json      bytes.Buffer
	}
	var cases = []caseStruct{
		{UserToken: tokens[1], Status: http.StatusNotAcceptable, Json: past},
		{UserToken: tokens[0], Status: http.StatusNotAcceptable, Json: past},
		{UserToken: tokens[1], Status: http.StatusCreated, Json: future},
	}

	router := SetupRouter()

	for _, testData := range cases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/events/inscription", &testData.Json)
		req.Header.Set("Authorization", "Bearer "+testData.UserToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, testData.Status, w.Code)
	}
}

func TestGetInscription(t *testing.T) {
	type caseStruct struct {
		UserToken string
		Status    int
	}
	var cases = []caseStruct{
		{UserToken: tokens[1], Status: http.StatusOK},
		{UserToken: tokens[0], Status: http.StatusOK},
	}

	router := SetupRouter()

	for _, testData := range cases {
		w := httptest.NewRecorder()
		js, _ := createJsonInscriptionGet()
		req, _ := http.NewRequest(http.MethodGet, "/events/inscription", &js)
		req.Header.Set("Authorization", "Bearer "+testData.UserToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, testData.Status, w.Code)
	}
}

func TestNotFound(t *testing.T) {
	router := SetupRouter()
	message := "Page not found"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/admin", nil)
	router.ServeHTTP(w, req)
	messageIn, _ := searchElemInBody(w.Body.String(), "message")

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, message, messageIn)
}

func createJsonEvent(dateString string) (bytes.Buffer, error) {
	organizer, _ := primitive.ObjectIDFromHex("630bf88273e7d5e13cde99b1")
	nDate, _ := time.Parse("2006-01-02", dateString)
	date := primitive.NewDateTimeFromTime(nDate)
	je := model.Event{
		Title:             "titulo test interno",
		Description_small: "descript corta",
		Description_large: "descript larga",
		DateOfEvent:       date,
		Organizer:         organizer,
		Place:             "lugar",
		Status:            true,
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(je)
	return buf, err
}

func createJsonInscriptionPost(idEvent string) (bytes.Buffer, error) {
	oidEvent, _ := primitive.ObjectIDFromHex(idEvent)
	oidUser, _ := primitive.ObjectIDFromHex("630bf88273e7d5e13cde99b0")

	je := model.Inscription{
		Event: oidEvent,
		User:  oidUser,
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(je)
	return buf, err
}

func createJsonInscriptionGet() (bytes.Buffer, error) {
	oidUser, _ := primitive.ObjectIDFromHex("630bf88273e7d5e13cde99b0")

	je := model.Inscription{
		User: oidUser,
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(je)
	return buf, err
}

func searchElemInBody(body string, element string) (string, error) {
	var response map[string]string
	err := json.Unmarshal([]byte(body), &response)
	if value, exists := response[element]; exists {
		return value, nil
	}
	return "", err
}
