package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"
)

// Define Queries for tokens table.
const (
	INSERT_TOKEN_QUERY = `INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES ($1, $2, $3, $4)
	`

	DELETE_ALL_FOR_USER_QUERY = `DELETE FROM tokens WHERE scope = $1 AND user_id = $2
	`
)

// Define constants for the token scope. For now we just define the scope "activation"
// Add ScopeAuthetication (v2)
const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

type TokenModel struct {
	DB *sql.DB
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	// Initialize a zero-value slice of bytes with a length of 16 bytes.
	randomBytes := make([]byte, 16)

	// Use the Read() function from the crypto/rand package to fill the byte slice with
	// random bytes from your operating system's CSPRNG.
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hashArray := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hashArray[:]

	return token, nil
}

func (tokenModel TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = tokenModel.Insert(token)
	return token, err
}

func (tokenModel TokenModel) Insert(token *Token) error {
	args := []interface{}{
		token.Hash,
		token.UserID,
		token.Expiry,
		token.Scope,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := tokenModel.DB.ExecContext(ctx, INSERT_TOKEN_QUERY, args...)
	return err
}

func (tokenModel TokenModel) DeleteAllForUser(scope string, userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := tokenModel.DB.ExecContext(ctx, DELETE_ALL_FOR_USER_QUERY, scope, userID)
	return err
}
