package models

import (
	"fmt"
	"testing"

	"github.com/keydotcat/backend/util"
)

var a32b = make([]byte, 32)

func getDummyUser() *User {
	ctx := getCtx()
	uid := util.GenerateRandomToken(10)
	vkp := VaultKeyPair{a32b, map[string][]byte{uid: a32b}}
	u, err := NewUser(ctx, uid, "uid fullname", uid+"@nowhere.net", uid, a32b, a32b, vkp)
	if err != nil {
		panic(err)
	}
	return u
}

func TestCreateUser(t *testing.T) {
	ctx := getCtx()
	vkp := VaultKeyPair{make([]byte, 32), map[string][]byte{"test": []byte("crap")}}
	u, err := NewUser(ctx, "test", "easdsa", "asdas@asdas.com", "somepass", make([]byte, 32), make([]byte, 32), vkp)
	if err != nil {
		fmt.Println(util.GetStack(err))
		t.Fatal(err)
	}
	if err = u.CheckPassword("somepass"); err != nil {
		fmt.Println("Password didn't check")
	}
	if u.Id != "test" {
		fmt.Println(util.GetStack(err))
		t.Errorf("Invalid username: %s vs test", u.Id)
	}
	u, err = NewUser(ctx, "test", "easdsa", "asdas@asdas.com", "somepass", make([]byte, 32), make([]byte, 32), vkp)
	if err != nil && "Username already taken" != err.Error() {
		t.Fatal(err)
	}
	teams := u.GetTeams(ctx)
	if len(teams) != 1 {
		t.Fatalf("Expected to have 1 team and got %d", len(teams))
	}
	team := teams[0]
	if team.Name != u.FullName {
		t.Errorf("Team name mismatch expected %s and got %s", u.FullName, team.Name)
	}
	if team.Owner != u.Id {
		t.Errorf("Team owner mismatch expected %s and got %s", u.Id, team.Owner)
	}
	vaults := team.GetVaultsForUser(ctx, u)
	if len(vaults) != 1 {
		t.Fatalf("Expected to have 1 vaults and got %d", len(vaults))
	}
	vault := vaults[0]
	if vault.Id != DEFAULT_VAULT_NAME {
		t.Errorf("Team name mismatch expected %s and got %s", DEFAULT_VAULT_NAME, vault.Id)
	}
}
