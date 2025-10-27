package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// BoolFlex permite deserializar valores booleanos flexibles
type BoolFlex bool

func (b *BoolFlex) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*b = false
		return nil
	}

	var boolVal bool
	if err := json.Unmarshal(data, &boolVal); err == nil {
		*b = BoolFlex(boolVal)
		return nil
	}

	var strVal string
	if err := json.Unmarshal(data, &strVal); err == nil {
		switch strVal {
		case "1", "true":
			*b = true
		case "0", "false", "":
			*b = false
		default:
			return fmt.Errorf("valor no soportado para BoolFlex: %s", strVal)
		}
		return nil
	}

	return fmt.Errorf("no se pudo deserializar BoolFlex: %s", string(data))
}

func (b *BoolFlex) MarshalJSON() ([]byte, error) {
	return json.Marshal(bool(*b))
}

// IntFlex permite deserializar valores enteros flexibles
type IntFlex int

func (i *IntFlex) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*i = 0
		return nil
	}

	var intVal int
	if err := json.Unmarshal(data, &intVal); err == nil {
		*i = IntFlex(intVal)
		return nil
	}

	var strVal string
	if err := json.Unmarshal(data, &strVal); err == nil {
		if strVal == "" {
			*i = 0
			return nil
		}

		parsedInt, err := strconv.Atoi(strVal)
		if err != nil {
			return fmt.Errorf("valor no soportado para IntFlex: %s", strVal)
		}
		*i = IntFlex(parsedInt)
		return nil
	}

	return fmt.Errorf("no se pudo deserializar IntFlex: %s", string(data))
}

func (i *IntFlex) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(*i))
}

// FloatFlex permite deserializar valores flotantes que pueden venir como números o strings
type FloatFlex float64

func (f *FloatFlex) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*f = 0
		return nil
	}

	// Intenta como float64
	var floatVal float64
	if err := json.Unmarshal(data, &floatVal); err == nil {
		*f = FloatFlex(floatVal)
		return nil
	}

	// Intenta como string
	var strVal string
	if err := json.Unmarshal(data, &strVal); err == nil {
		if strVal == "" {
			*f = 0
			return nil
		}

		parsedFloat, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return fmt.Errorf("valor no soportado para FloatFlex: %s", strVal)
		}
		*f = FloatFlex(parsedFloat)
		return nil
	}

	return fmt.Errorf("no se pudo deserializar FloatFlex: %s", string(data))
}

// MarshalJSON implementa la interfaz json.Marshaler para FloatFlex
func (f *FloatFlex) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(*f))
}
