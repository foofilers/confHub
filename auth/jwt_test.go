package auth

import (
	"testing"
	"time"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func TestMain(m *testing.M) {
	logrus.SetLevel(logrus.DebugLevel)
	os.Exit(m.Run())
}

func TestValidJwt(t *testing.T) {
	origUserId := "Bob"
	token, err := GenerateJwt(origUserId, 1 * time.Second)
	if err != nil {
		t.Fatal(err)
	}
	userId, err := ValidateJwt(token)
	if err != nil {
		t.Fatal(err)
	}
	if userId != origUserId {
		t.Fatalf("UserId %s should be %s", userId, origUserId)
	}
}

func TestNoExpirationJwt(t *testing.T) {
	origUserId := "Alice"
	token, err := GenerateJwt(origUserId, 0)
	if err != nil {
		t.Fatal(err)
	}
	userId, err := ValidateJwt(token)
	if err != nil {
		t.Fatal(err)
	}
	if userId != origUserId {
		t.Fatalf("UserId %s should be %s", userId, origUserId)
	}
}

func TestExpirationJwt(t *testing.T) {
	origUserId := "Bob"
	token, err := GenerateJwt(origUserId, 1 * time.Second)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)
	_, err = ValidateJwt(token)
	if err == nil || !strings.Contains(err.Error(), "expired") {
		t.Fatalf("The token should be expired")
	}
}
