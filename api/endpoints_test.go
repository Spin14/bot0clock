package api

import (
	"bytes"
	"fmt"
	"github.com/spin14/bot0clock/model"
	"net/http"
	"net/http/httptest"
	"testing"
)


const testUsername = "spin14"

func createTestUser(s *model.Storage) {
	_, _ = s.Create(testUsername)
}

func TestRetrieveUser(t *testing.T) {
	s := model.InitUserModelTable()
	defer model.CleanUserModelTable()

	path := fmt.Sprintf("/%s", testUsername)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	r := Router(s)

	rr404 := httptest.NewRecorder()

	r.ServeHTTP(rr404, req)

	if status := rr404.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
		return
	}

	createTestUser(s)
	rr200 := httptest.NewRecorder()

	r.ServeHTTP(rr200, req)

	if status := rr200.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
		return
	}

	expected := `{"username":"spin14"}
`
	if rr200.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got .%v. want .%v.",
			rr200.Body.String(), expected)
		return
	}
}

func TestCreateUser(t *testing.T) {
	s := model.InitUserModelTable()
	defer model.CleanUserModelTable()

	jsonData := `{"username":"spin14"}
`

	req, err := http.NewRequest("POST", "/", bytes.NewBufferString(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	r := Router(s)

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	expected := jsonData

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got .%v. want .%v.",
			rr.Body.String(), expected)
		return
	}

	if _, err := s.Retrieve(testUsername); err != nil {
		t.Error("user was not created")
	}
}


func TestUserList(t *testing.T) {
	s := model.InitUserModelTable()
	defer model.CleanUserModelTable()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	r := Router(s)

	rr := httptest.NewRecorder()

	_, _ = s.Create("a")
	_, _ = s.Create("b")
	_, _ = s.Create("c")

	r.ServeHTTP(rr, req)


	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
		return
	}

	expected := `[{"username":"a"},{"username":"b"},{"username":"c"}]
`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got .%v. want .%v.",
			rr.Body.String(), expected)
		return
	}
}

func TestUpdateUser404(t *testing.T) {
	s := model.InitUserModelTable()
	defer model.CleanUserModelTable()

	jsonData := `{"username":"acci"}
`

	path := fmt.Sprintf("/%s", testUsername)
	req, err := http.NewRequest("PUT", path, bytes.NewBufferString(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	r := Router(s)

	rr404 := httptest.NewRecorder()

	r.ServeHTTP(rr404, req)

	if status := rr404.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
		return
	}
}

func TestUpdateUser200(t *testing.T) {
	s := model.InitUserModelTable()
	defer model.CleanUserModelTable()

	jsonData := `{"username":"acci"}
`

	path := fmt.Sprintf("/%s", testUsername)
	req, err := http.NewRequest("PUT", path, bytes.NewBufferString(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	r := Router(s)

	createTestUser(s)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
		return
	}

	expected := `{"username":"acci"}
`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got .%v. want .%v.",
			rr.Body.String(), expected)
		return
	}
}