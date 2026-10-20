package auth

import "testing"

func TestHashPassword(t *testing.T) {
	password := "thisisatestpasswordtotesthashing"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	if hash == "" {
		t.Fatal("HashPassword returned an empty hash")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "thisisatestpasswordtotesthashing"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned an error: %v", err)
	}

	if match == false {
		t.Fatal("CheckPasswordHash returned False")
	}

}
