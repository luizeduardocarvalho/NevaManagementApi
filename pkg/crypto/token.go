package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateInvitationToken generates a cryptographically secure random token
// with the prefix "inv_" for identification
func GenerateInvitationToken() (string, error) {
	// Generate 32 random bytes (256 bits of entropy)
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode as hex and add prefix
	return fmt.Sprintf("inv_%s", hex.EncodeToString(randomBytes)), nil
}
