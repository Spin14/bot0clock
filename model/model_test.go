package model

import (
	"testing"
)

const testUsername = "spin14"


func TestStorage_Create(t *testing.T) {
	s := InitUserModelTable()
	defer CleanUserModelTable()

	if user, err := s.Create(testUsername); err != nil || user.Username != testUsername {
		t.Error("failed create error")
		return
	}

	if s.Count() != 1 {
		t.Error("count error")
		return
	}

	if user, err := s.Create(testUsername); err == nil || user != nil {
		t.Error("create error (not unique)")
		return
	}

	if s.Count() != 1 {
		t.Error("count error")
		return
	}
}


func TestStorage_ListAll(t *testing.T) {
	s := InitUserModelTable()
	defer CleanUserModelTable()

	_, _ = s.Create(testUsername)
	_, _ = s.Create("chi")
	_, _ = s.Create("vasmv")

	if users, err := s.ListAll(); err != nil || len(users) != 3 {
		t.Error("list all error")
		return
	}

}

func TestStorage_Retrieve(t *testing.T) {
	s := InitUserModelTable()
	defer CleanUserModelTable()

	_, _ = s.Create(testUsername)
	if user, err := s.Retrieve(testUsername); err != nil || user.Username != testUsername {
		t.Error("failed retrieve error")
		return
	}

	if _, err := s.Retrieve("notFound"); err == nil {
		t.Error("failed retrieve error")
		return
	}
}

func TestStorage_Update(t *testing.T) {
	s := InitUserModelTable()
	defer CleanUserModelTable()

	_, _ = s.Create(testUsername)
	if user, err := s.Update(testUsername, "newName"); err != nil || user.Username != "newName" {
		t.Error("failed update error")
		return
	}

	if _, err := s.Update("notFound", "-"); err == nil {
		t.Error("failed update error")
		return
	}
}
