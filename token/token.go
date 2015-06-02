package token

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/matter123/chat/config"
	"time"
)

//Token is a transient Identifier for users
type tokenObj struct {
	user     string
	token    string
	lastUsed time.Time
}

var tokens []tokenObj

func notExpired(token tokenObj) bool {
	return time.Since(token.lastUsed) < time.Duration(config.Settings().TokenSettings.Expire)*time.Second
}

//Valid checks if the given token exists and is not expired
func Valid(token string) bool {
	for _, tok := range tokens {
		if tok.token == token {
			return notExpired(tok)
		}
	}
	return false
}

//User returns the user associated with the token or nil if the token is not valid
func User(token string) string {
	for _, tok := range tokens {
		if tok.token == token {
			if notExpired(tok) {
				return tok.user
			}
			break
		}
	}
	return ""
}

func genToken() string {
	tokenlen := config.Settings().TokenSettings.Size / 8
	token := make([]byte, tokenlen)
	for {
		_, err := rand.Read(token)
		if err == nil {
			return hex.EncodeToString(token)
		}
	}
}

//Token returns a token for the givin user, generating a new token if one doesn't exist or is expired
func Token(user string) string {
	for _, tok := range tokens {
		if tok.user == user {
			if notExpired(tok) {
				return tok.token
			}
			tok.token = genToken()
			tok.lastUsed = time.Now()
			return tok.token
		}
	}
	//if we got here the old token must be gone
	token := genToken()
	tokens = append(tokens, tokenObj{
		user:     user,
		token:    token,
		lastUsed: time.Now(),
	})
	return token
}

//Renew updates the expire time of the given token
//returns an error if token is invalid
func Renew(token string) error {
	for _, tok := range tokens {
		if tok.token == token {
			if notExpired(tok) {
				tok.lastUsed = time.Now()
				return nil
			}
			break
		}
	}
	return errors.New("Token is invalid")
}
