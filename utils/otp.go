package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"
)

type OTPEntry struct {
	OTP       string
	ExpiresAt time.Time
}

var otpStore = make(map[string]OTPEntry)
var mu sync.Mutex

// GenerateOTP generates a 6-digit numeric OTP and stores it with email or phone
func GenerateOTP(key string) (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	otp := fmt.Sprintf("%06d", n.Int64())

	mu.Lock()
	otpStore[key] = OTPEntry{
		OTP:       otp,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	mu.Unlock()

	log.Printf("✅ Generated OTP for %s: %s\n", key, otp) // 🔥 Add this line for debugging

	return otp, nil
}


// VerifyOTP checks if the submitted OTP is correct and not expired
func VerifyOTP(key, submitted string) bool {
	mu.Lock()
	entry, exists := otpStore[key]
	mu.Unlock()

	if !exists {
		return false
	}

	if time.Now().After(entry.ExpiresAt) {
		return false
	}

	return entry.OTP == submitted
}

// Optional: DeleteOTP after use (for better security)
func DeleteOTP(key string) {
	mu.Lock()
	delete(otpStore, key)
	mu.Unlock()
}
