package model

import (
	"testing"
)

func TestStorage_ListAll(t *testing.T) {

}

func TestStorage_Create(t *testing.T) {
	initUserModelTable()
	defer cleanUserModelTable()

	s := testStorage()
	if s.Count() != 0 {
		t.Error()
		return
	}

	testUsername := "spin14"

	if user, err := s.Create(testUsername); err != nil || user.Username != testUsername {
		t.Error("failed create error")
		return
	}

	if s.Count() != 1 {
		t.Error()
		return
	}

	if user, err := s.Create(testUsername); err == nil || user != nil {
		t.Error("create error (not unique)")
		return
	}

	if s.Count() != 1 {
		t.Error()
		return
	}
}
