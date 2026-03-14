package security

import (
	"testing"
	"time"
)

func TestHMACTokenServiceIssueAndVerify(testingSuite *testing.T) {
	tokenService, creationError := NewHMACTokenService(WithTokenSecret("phase1-secret"))
	if creationError != nil {
		testingSuite.Fatalf("token service creation failure: %v", creationError)
	}
	claims := TokenClaims{SubjectUserID: "user", AssignedRole: "admin", TokenType: "access", IssuedAtUnix: time.Now().Unix(), ExpiresAtUnix: time.Now().Add(5 * time.Minute).Unix()}
	tokenValue, issueError := tokenService.IssueToken(claims)
	if issueError != nil {
		testingSuite.Fatalf("issue token failure: %v", issueError)
	}
	verifiedClaims, verificationError := tokenService.VerifyToken(tokenValue)
	if verificationError != nil {
		testingSuite.Fatalf("verify token failure: %v", verificationError)
	}
	if verifiedClaims.SubjectUserID != claims.SubjectUserID {
		testingSuite.Fatalf("unexpected claim subject")
	}
}

func TestSHA256PasswordHasher(testingSuite *testing.T) {
	passwordHasher := NewSHA256PasswordHasher()
	passwordHash := passwordHasher.HashPassword("password")
	if !passwordHasher.MatchesHash("password", passwordHash) {
		testingSuite.Fatalf("expected password hash match")
	}
	if passwordHasher.MatchesHash("wrong-password", passwordHash) {
		testingSuite.Fatalf("expected password hash mismatch")
	}
}
