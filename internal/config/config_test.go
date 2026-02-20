//nolint:testpackage // Testing GameConfig and SetDevInvincible from same package
package config

import (
	"testing"
)

func TestSetDevInvincible_OnlyWhenDebug(t *testing.T) {
	// When Debug is false, SetDevInvincible does not set Invincible
	cfg := NewConfig(WithDebug(false))
	cfg.SetDevInvincible(true)
	if cfg.Invincible {
		t.Error("Invincible should stay false when Debug is false")
	}

	// When Debug is true, SetDevInvincible sets Invincible
	cfg2 := NewConfig(WithDebug(true))
	cfg2.SetDevInvincible(true)
	if !cfg2.Invincible {
		t.Error("Invincible should be true after SetDevInvincible(true) when Debug is true")
	}
	cfg2.SetDevInvincible(false)
	if cfg2.Invincible {
		t.Error("Invincible should be false after SetDevInvincible(false)")
	}
}
