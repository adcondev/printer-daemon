package models

import (
	"encoding/json"
	"testing"
)

func TestBoolFlex_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected BoolFlex
		wantErr  bool
	}{
		{"true boolean", `true`, true, false},
		{"false boolean", `false`, false, false},
		{"string 1", `"1"`, true, false},
		{"string 0", `"0"`, false, false},
		{"string true", `"true"`, true, false},
		{"string false", `"false"`, false, false},
		{"empty string", `""`, false, false},
		{"null", `null`, false, false},
		{"invalid", `"invalid"`, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b BoolFlex
			err := json.Unmarshal([]byte(tt.input), &b)
			if (err != nil) != tt.wantErr {
				t.Errorf("BoolFlex.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && b != tt.expected {
				t.Errorf("BoolFlex.UnmarshalJSON() = %v, want %v", b, tt.expected)
			}
		})
	}
}

func TestIntFlex_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected IntFlex
		wantErr  bool
	}{
		{"integer", `42`, 42, false},
		{"string integer", `"123"`, 123, false},
		{"zero", `0`, 0, false},
		{"string zero", `"0"`, 0, false},
		{"empty string", `""`, 0, false},
		{"null", `null`, 0, false},
		{"invalid string", `"abc"`, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var i IntFlex
			err := json.Unmarshal([]byte(tt.input), &i)
			if (err != nil) != tt.wantErr {
				t.Errorf("IntFlex.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && i != tt.expected {
				t.Errorf("IntFlex.UnmarshalJSON() = %v, want %v", i, tt.expected)
			}
		})
	}
}

func TestFloatFlex_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected FloatFlex
		wantErr  bool
	}{
		{"float", `123.45`, 123.45, false},
		{"integer", `100`, 100, false},
		{"string float", `"67.89"`, 67.89, false},
		{"string integer", `"200"`, 200, false},
		{"zero", `0`, 0, false},
		{"string zero", `"0"`, 0, false},
		{"empty string", `""`, 0, false},
		{"null", `null`, 0, false},
		{"invalid string", `"not_a_number"`, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f FloatFlex
			err := json.Unmarshal([]byte(tt.input), &f)
			if (err != nil) != tt.wantErr {
				t.Errorf("FloatFlex.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && f != tt.expected {
				t.Errorf("FloatFlex.UnmarshalJSON() = %v, want %v", f, tt.expected)
			}
		})
	}
}

func TestFlexTypes_MarshalJSON(t *testing.T) {
	t.Run("BoolFlex marshal", func(t *testing.T) {
		b := BoolFlex(true)
		data, err := json.Marshal(b)
		if err != nil {
			t.Fatalf("BoolFlex.MarshalJSON() error = %v", err)
		}
		if string(data) != "true" {
			t.Errorf("BoolFlex.MarshalJSON() = %s, want true", data)
		}
	})

	t.Run("IntFlex marshal", func(t *testing.T) {
		i := IntFlex(42)
		data, err := json.Marshal(i)
		if err != nil {
			t.Fatalf("IntFlex.MarshalJSON() error = %v", err)
		}
		if string(data) != "42" {
			t.Errorf("IntFlex.MarshalJSON() = %s, want 42", data)
		}
	})

	t.Run("FloatFlex marshal", func(t *testing.T) {
		f := FloatFlex(123.45)
		data, err := json.Marshal(f)
		if err != nil {
			t.Fatalf("FloatFlex.MarshalJSON() error = %v", err)
		}
		if string(data) != "123.45" {
			t.Errorf("FloatFlex.MarshalJSON() = %s, want 123.45", data)
		}
	})
}
