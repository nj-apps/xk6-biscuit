package xk6Biscuit

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	biscuit "github.com/biscuit-auth/biscuit-go/v2"
	"github.com/biscuit-auth/biscuit-go/v2/parser"
	"go.k6.io/k6/js/modules"
	"golang.org/x/crypto/ed25519"
)

// init is called by the Go runtime at application startup.
func init() {
	modules.Register("k6/x/biscuit", new(Biscuit))
}

// ModuleInstance represents an instance of the JS module.
type (
	Biscuit struct{}
)

func biscuitFromToken(b64Token string) (*biscuit.Biscuit, error) {

	serializedToken, err := base64.URLEncoding.DecodeString(b64Token)
	if err != nil {
		return nil, fmt.Errorf("Error decoding token : %v", err)
	}

	biscuit, err := biscuit.Unmarshal(serializedToken)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize biscuit: %v", err)
	}

	return biscuit, nil
}

func (b *Biscuit) CreateKey() (ed25519.PublicKey, ed25519.PrivateKey) {
	rng := rand.Reader
	publicRoot, privateRoot, _ := ed25519.GenerateKey(rng)
	return publicRoot, privateRoot
}

// Inspect deserializes the token en returns the content in a string
func (b *Biscuit) Inspect(b64Token string) (string, error) {
	token, err := biscuitFromToken(b64Token)
	if err != nil {
		return "", err
	}
	return token.String(), nil
}

// Seal prevents a biscuit from being attenuated further
func (t *Biscuit) Seal(b64Token string) (string, error) {
	b, err := biscuitFromToken(b64Token)
	if err != nil {
		return "", err
	}
	sealed, err := b.Seal(rand.Reader)
	if err != nil {
		return "", err
	}

	return t.serialize(sealed)
}

func (t *Biscuit) serialize(b *biscuit.Biscuit) (string, error) {
	serialized, err := b.Serialize()
	if err != nil {
		return "", fmt.Errorf("failed to serialize biscuit: %v", err)
	}
	return base64.URLEncoding.Strict().EncodeToString(serialized), nil
}

// Attenuate creates a new token with the provided block appended
func (t *Biscuit) Attenuate(b64Token string, blocks []string) (string, error) {

	// Decode the base64 token
	serializedToken, err := base64.URLEncoding.DecodeString(b64Token)
	if err != nil {
		return "", fmt.Errorf("Error decoding token : %v", err)
	}

	// Deserialize the token
	token, err := biscuit.Unmarshal(serializedToken)
	if err != nil {
		return "", fmt.Errorf("failed to deserialize biscuit: %v", err)
	}
	// Create a new authorization block
	blockBuilder := token.CreateBlock()

	// Create constraints to be verified by server
	for _, s := range blocks {
		check, err := parser.FromStringCheck(s)
		if err != nil {
			return "", fmt.Errorf("failed to parse check: %v", err)
		}
		// Add validations to the autorization block
		err = blockBuilder.AddCheck(check)
		if err != nil {
			return "", fmt.Errorf("failed to add block check: %v", err)
		}
	}

	// Append the new autorization block to the token
	rng := rand.Reader
	token2, err := token.Append(rng, blockBuilder.Build())
	if err != nil {
		return "", fmt.Errorf("failed to append: %v", err)
	}
	// Return new token
	return t.serialize(token2)

}
