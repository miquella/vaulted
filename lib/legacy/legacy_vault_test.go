package legacy_test

import (
	"encoding/json"
	"testing"

	"github.com/miquella/vaulted/lib/legacy"
)

const (
	VAULT_LEGACY = `{"keyDetails":{"digest":"sha-512","salt":"r9Yc86Nw9Vnjf64CX9PtAzSflWYjF1893jgwqb5UprE=","iterations":65536},"macDigest":"sha-256","cipher":"aes","cipherMode":"ctr","mac":"5aUIwZpewtrt/HPHqd4ei+lOqqJuROJKFWlzIL/Xrlc=","iv":"heeKgQjdoDDWNARV7ugs/g==","environments":"p37jXVGr4m1UM9vRvwu/bwSYkwrhKSg4q9f6eQSn/imMbixaXwrabCAGPuAviSKRrg=="}`
)

func TestVault(t *testing.T) {
	v := legacy.Vault{}
	err := json.Unmarshal([]byte(VAULT_LEGACY), &v)
	if err != nil {
		t.Fatalf("failed to unmarshal vault: %v", err)
	}

	_, err = v.DecryptEnvironments("invalid password")
	if err == nil {
		t.Fatal("Invalid password accepted!")
	}

	envs, err := v.DecryptEnvironments("test")
	if err != nil {
		t.Fatalf("Valid password not accepted! %v", err)
	}

	if envs["test"].Vars["TEST"] != "legacy" {
		t.Fatalf("Expected: legacy, got %s", envs["test"].Vars["TEST"])
	}
}
