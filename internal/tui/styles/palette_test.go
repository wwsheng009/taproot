package styles

import "testing"

func TestPaletteColors(t *testing.T) {
	// Test that all palette color constants are valid hex colors
	colors := []struct {
		name  string
		value string
	}{
		{"Smoke", ColorSmoke},
		{"Malibu", ColorMalibu},
		{"Zest", ColorZest},
		{"Charple", ColorCharple},
		{"Guac", ColorGuac},
		{"Charcoal", ColorCharcoal},
		{"Coral", ColorCoral},
		{"Butter", ColorButter},
		{"Sriracha", ColorSriracha},
		{"Oyster", ColorOyster},
		{"Bengal", ColorBengal},
		{"Pony", ColorPony},
		{"Guppy", ColorGuppy},
		{"Salmon", ColorSalmon},
		{"Cheeky", ColorCheeky},
		{"Mauve", ColorMauve},
		{"Hazy", ColorHazy},
		{"Salt", ColorSalt},
		{"Citron", ColorCitron},
		{"Julep", ColorJulep},
		{"Cumin", ColorCumin},
		{"Bok", ColorBok},
		{"Zinc", ColorZinc},
		{"Squid", ColorSquid},
	}

	for _, c := range colors {
		t.Run(c.name, func(t *testing.T) {
			if len(c.value) != 7 {
				t.Errorf("Color %s = %q, want length 7 (hex format)", c.name, c.value)
			}
			if c.value[0] != '#' {
				t.Errorf("Color %s = %q, want '#' prefix", c.name, c.value)
			}
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	t.Run("stringPtr", func(t *testing.T) {
		s := "test"
		ptr := stringPtr(s)
		if ptr == nil {
			t.Fatal("stringPtr() returned nil")
		}
		if *ptr != s {
			t.Errorf("stringPtr() = %p, want %q", ptr, s)
		}
	})

	t.Run("boolPtr", func(t *testing.T) {
		b := true
		ptr := boolPtr(b)
		if ptr == nil {
			t.Fatal("boolPtr() returned nil")
		}
		if *ptr != b {
			t.Errorf("boolPtr() = %v, want %v", *ptr, b)
		}
	})

	t.Run("uintPtr", func(t *testing.T) {
		u := uint(42)
		ptr := uintPtr(u)
		if ptr == nil {
			t.Fatal("uintPtr() returned nil")
		}
		if *ptr != u {
			t.Errorf("uintPtr() = %d, want %d", *ptr, u)
		}
	})
}
