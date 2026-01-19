package config

import (
	"testing"
)

func TestIsPlain_KamajiPlainEnvSet(t *testing.T) {
	ResetPlain()
	t.Setenv("KAMAJI_PLAIN", "1")
	t.Setenv("NO_COLOR", "")

	if !IsPlain() {
		t.Error("IsPlain() should return true when KAMAJI_PLAIN is set")
	}
}

func TestIsPlain_NoColorEnvSet(t *testing.T) {
	ResetPlain()
	t.Setenv("KAMAJI_PLAIN", "")
	t.Setenv("NO_COLOR", "1")

	if !IsPlain() {
		t.Error("IsPlain() should return true when NO_COLOR is set")
	}
}

func TestIsPlain_NeitherSet(t *testing.T) {
	ResetPlain()
	t.Setenv("KAMAJI_PLAIN", "")
	t.Setenv("NO_COLOR", "")

	if IsPlain() {
		t.Error("IsPlain() should return false when neither env var is set")
	}
}

func TestIsPlain_KamajiPlainFalsyValues(t *testing.T) {
	falsyValues := []string{"0", "false", "no", "off", "random"}
	for _, val := range falsyValues {
		t.Run(val, func(t *testing.T) {
			ResetPlain()
			t.Setenv("KAMAJI_PLAIN", val)
			t.Setenv("NO_COLOR", "")

			if IsPlain() {
				t.Errorf("IsPlain() should return false when KAMAJI_PLAIN=%q", val)
			}
		})
	}
}

func TestSetPlain_OverridesDetection(t *testing.T) {
	ResetPlain()
	t.Setenv("KAMAJI_PLAIN", "")
	t.Setenv("NO_COLOR", "")

	SetPlain(true)
	if !IsPlain() {
		t.Error("SetPlain(true) should force plain mode")
	}

	SetPlain(false)
	if IsPlain() {
		t.Error("SetPlain(false) should disable plain mode")
	}
}

func TestResetPlain_ClearsOverride(t *testing.T) {
	SetPlain(true)
	ResetPlain()
	t.Setenv("KAMAJI_PLAIN", "")
	t.Setenv("NO_COLOR", "")

	if IsPlain() {
		t.Error("ResetPlain() should clear override allowing env detection")
	}
}
