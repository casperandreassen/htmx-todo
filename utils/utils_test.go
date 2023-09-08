package utils

import "testing"

func TestIssueJwtToken(t *testing.T) {
	token, err := IssueJwtToken(1, "ola")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	VerifyJwtToken(token)
}
