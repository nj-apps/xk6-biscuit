package xk6Biscuit

import (
	"encoding/base64"
	"strings"
	"testing"

	"crypto/rand"

	"fmt"

	biscuit "github.com/biscuit-auth/biscuit-go/v2"
	"github.com/biscuit-auth/biscuit-go/v2/parser"
	"golang.org/x/crypto/ed25519"
)

func TestAttenuate(t *testing.T) {
	b := Biscuit{}

	// generate a new token
	rng := rand.Reader
	_, privateRoot, _ := ed25519.GenerateKey(rng)
	token := createToken(privateRoot)

	// inspect the token
	s, e := b.Inspect(token)
	if e != nil {
		t.Fatalf("Failed to inspect token: %s", e.Error())
	}
	t.Logf("Token details: %s\n", s)

	// Attenuate
	blocks := []string{
		"check if time($now), $now < 2023-12-30T15:44:00Z",
		`check if txn::service("Canal")`}

	att, e := b.Attenuate(token, blocks)
	if e != nil {
		t.Fatalf("Failed to attenuate token: %s", e.Error())
	}

	s, _ = b.Inspect(att)
	if !(strings.Contains(s, "Canal") && strings.Contains(s, "2023-12-30T15:44:00Z")) {
		t.Fatalf("Failed to attenuate token blocks not found in attenuated token : %s \n", s)
	}

}

func TestSealing(t *testing.T) {
	b := Biscuit{}

	// generate a new token
	rng := rand.Reader
	_, privateRoot, _ := ed25519.GenerateKey(rng)
	token := createToken(privateRoot)

	// seal the token
	sealed, e := b.Seal(token)
	if e != nil {
		t.Fatalf("Error sealing token : %s", e.Error())
	}

	// try to attenuate
	_, e = b.Attenuate(sealed, []string{})
	if e == nil {
		t.Fatal("Error : attenuate succeed but a sealed token may NOT be Attenuated.")
	}

}

func createToken(privateRoot ed25519.PrivateKey) string {

	authority, err := parser.FromStringBlockWithParams(`
		right("/a/file1.txt", {read});
		right("/a/file1.txt", {write});
		right("/a/file2.txt", {read});
		right("/a/file3.txt", {write});
	`, map[string]biscuit.Term{"read": biscuit.String("read"), "write": biscuit.String("write")})

	if err != nil {
		panic(fmt.Errorf("failed to parse authority block: %v", err))
	}

	builder := biscuit.NewBuilder(privateRoot)
	builder.AddBlock(authority)

	b, err := builder.Build()
	if err != nil {
		panic(fmt.Errorf("failed to build biscuit: %v", err))
	}

	token, err := b.Serialize()
	if err != nil {
		panic(fmt.Errorf("failed to serialize biscuit: %v", err))
	}

	return base64.URLEncoding.EncodeToString(token)
}
