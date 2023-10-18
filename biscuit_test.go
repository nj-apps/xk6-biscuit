package xk6Biscuit

import (
	"encoding/base64"
	"testing"

	"crypto/rand"

	"fmt"

	biscuit "github.com/biscuit-auth/biscuit-go/v2"
	"github.com/biscuit-auth/biscuit-go/v2/parser"
	"golang.org/x/crypto/ed25519"
)

/*
	type mockVU struct {
		name string
	}

// instantiate the interface

	func (mockVU) Context() context.Context {
		return context.TODO()
	}

	func (mockVU) InitEnv() *common.InitEnvironment {
		return nil
	}

	func (mockVU) State() *lib.State {
		return nil
	}

	func (mockVU) Runtime() *goja.Runtime {
		return &goja.Runtime{}
	}

	func (mockVU) RegisterCallback() (enqueueCallback func(func() error)) {
		return enqueueCallback
	}

	func TestFIFO(t *testing.T) {
		Arguments := []goja.Value{}

		client1 := newClient(Arguments, mockVU{})
		client1.Push("first value")
		client1.Push("2nd value")
		client2 := newClient(Arguments, mockVU{})
		out1, _ := client1.Pop()
		out2, _ := client2.Pop()
		if out1 != "first value" || out2 != "2nd value" {
			t.Errorf("Single fifo : out1=%s out2=%s", out1, out2)
		}
	}

	func TestNamedFIFO(t *testing.T) {
		vm := goja.New()

		Arguments := []goja.Value{
			vm.ToValue("liste A"),
		}
		client1 := newClient(Arguments, mockVU{})
		client1.Push("first value A")
		client1.Push("2nd value A")

		Arguments = []goja.Value{
			vm.ToValue("liste B"),
		}
		client2 := newClient(Arguments, mockVU{})
		client2.Push("first value B")
		client2.Push("2nd value B")

		client3 := newClient(Arguments, mockVU{})

		out1, _ := client1.Pop()
		out2, _ := client2.Pop()
		out3, _ := client3.Pop()

		if out1 != "first value A" || out2 != "first value B" || out3 != "2nd value B" {
			t.Errorf("Named fifos : out1=%s out2=%s out3=%s", out1, out2, out3)
		}

}
*/
func TestBiscuit(t *testing.T) {

	rng := rand.Reader
	publicRoot, privateRoot, _ := ed25519.GenerateKey(rng)
	fmt.Print("public:", publicRoot)

	token := createToken(privateRoot)
	fmt.Printf("token=%v", token)

	b := Biscuit{}
	blocks := []string{
		"check if time($now), $now < 2023-12-30T15:44:00Z",
		`check if txn::service("Canal")`}

	t.Error(b.Inspect(token))

	sealed, err := b.Seal(token)
	if err != nil {
		t.Fatalf("Error sealing token : %v", err)
	}

	att, err := b.Attenuate(sealed, blocks)
	if err != nil {
		t.Fatalf("Error Attenuating token : %v", err)
	}

	t.Errorf("attenuated token=%v", att)

	t.Error(b.Inspect(att))

	/*

		b0 := Biscuit{}

		b := Biscuit{
			checks: []string{
				"check if time($now), $now < 2023-12-30T15:44:00Z",
				`check if txn::service("Canal")`},
		}

		bToken, err := base64.URLEncoding.DecodeString(token)
		if err != nil {
			t.Fatalf("Error decoding token : %v", err)
		}

		att0, err := b0.Attenuate(bToken, &publicRoot)
		if err != nil {
			t.Fatalf("Error Attenuating token : %v", err)
		}

		t.Errorf("original:\n%v\nattenuated:\n%v\n", token, base64.URLEncoding.EncodeToString(att0))

		att, err := b.Attenuate(bToken, &publicRoot)
		if err != nil {
			t.Fatalf("Error Attenuating token : %v", err)
		}

		t.Errorf("original:\n%v\nattenuated:\n%v\n", token, base64.URLEncoding.EncodeToString(att))
	*/
}

func createToken(privateRoot ed25519.PrivateKey) string {

	//fmt.Printf("publicRoot=%v", publicRoot)

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

	// token is now a []byte, ready to be shared
	// The biscuit spec mandates the use of URL-safe base64 encoding for textual representation:

	fmt.Println(base64.URLEncoding.EncodeToString(token))
	return base64.URLEncoding.EncodeToString(token)

}
