package main

import "testing"

func TestHashPassword_ProducesVerifiableHash(t *testing.T) {
	hash, err := hashPassword("correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("hashPassword returned an error: %v", err)
	}

	if hash == "" {
		t.Fatal("hashPassword returned an empty hash")
	}

	if hash == "correct-horse-battery-staple" {
		t.Fatal("hashPassword returned the plaintext password unchanged")
	}
}

func TestCheckPassword_AcceptsCorrectPassword(t *testing.T) {
	hash, err := hashPassword("correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("hashPassword returned an error: %v", err)
	}

	if !checkPassword("correct-horse-battery-staple", hash) {
		t.Fatal("checkPassword rejected the correct password")
	}
}

func TestCheckPassword_RejectsWrongPassword(t *testing.T) {
	hash, err := hashPassword("correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("hashPassword returned an error: %v", err)
	}

	if checkPassword("wrong-password", hash) {
		t.Fatal("checkPassword accepted an incorrect password")
	}
}
