package models

import (
	"encoding/json"
	"testing"
)

func TestConcepto_Serialization(t *testing.T) {
	concepto := Concepto{
		Clave:                 "TEST001",
		Descripcion:           "Test Product",
		Cantidad:              FloatFlex(2.5),
		Unidad:                "PIEZA",
		PrecioVenta:           FloatFlex(100.50),
		Total:                 FloatFlex(251.25),
		ClaveProductoServicio: "12345678",
		ClaveUnidadSAT:        "H87",
		VentaGranel:           BoolFlex(false),
		Series:                []string{"S1", "S2"},
		Impuestos: []Impuesto{
			{
				Factor:  "Tasa",
				Base:    FloatFlex(100),
				Importe: FloatFlex(16),
				Tasa:    FloatFlex(0.16),
				Tipo:    "T",
			},
		},
	}

	// Marshal
	data, err := json.Marshal(concepto)
	if err != nil {
		t.Fatalf("Failed to marshal Concepto: %v", err)
	}

	// Unmarshal
	var decoded Concepto
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal Concepto: %v", err)
	}

	// Verify
	if decoded.Clave != concepto.Clave {
		t.Errorf("Clave = %v, want %v", decoded.Clave, concepto.Clave)
	}
	if decoded.Cantidad != concepto.Cantidad {
		t.Errorf("Cantidad = %v, want %v", decoded.Cantidad, concepto.Cantidad)
	}
	if len(decoded.Series) != len(concepto.Series) {
		t.Errorf("Series length = %v, want %v", len(decoded.Series), len(concepto.Series))
	}
	if len(decoded.Impuestos) != 1 {
		t.Errorf("Impuestos length = %v, want 1", len(decoded.Impuestos))
	}
}

func TestImpuesto_Serialization(t *testing.T) {
	impuesto := Impuesto{
		Factor:    "Tasa",
		Base:      FloatFlex(1000),
		Importe:   FloatFlex(160),
		Impuestos: "002",
		Tasa:      FloatFlex(0.16),
		Entidad:   "Federal",
		Tipo:      "T",
	}

	data, err := json.Marshal(impuesto)
	if err != nil {
		t.Fatalf("Failed to marshal Impuesto: %v", err)
	}

	var decoded Impuesto
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal Impuesto: %v", err)
	}

	if decoded.Factor != impuesto.Factor {
		t.Errorf("Factor = %v, want %v", decoded.Factor, impuesto.Factor)
	}
	if decoded.Base != impuesto.Base {
		t.Errorf("Base = %v, want %v", decoded.Base, impuesto.Base)
	}
}

func TestPago_Serialization(t *testing.T) {
	pago := Pago{
		FormaPago:              "Efectivo",
		Cantidad:               FloatFlex(500),
		FormaPagoIdentificador: "EFE001",
	}

	data, err := json.Marshal(pago)
	if err != nil {
		t.Fatalf("Failed to marshal Pago: %v", err)
	}

	var decoded Pago
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal Pago: %v", err)
	}

	if decoded.FormaPago != pago.FormaPago {
		t.Errorf("FormaPago = %v, want %v", decoded.FormaPago, pago.FormaPago)
	}
	if decoded.Cantidad != pago.Cantidad {
		t.Errorf("Cantidad = %v, want %v", decoded.Cantidad, pago.Cantidad)
	}
}

func TestDocumentoPago_ComplexNesting(t *testing.T) {
	doc := DocumentoPago{
		Total:      FloatFlex(1000),
		TipoCambio: FloatFlex(1),
		Saldo:      FloatFlex(0),
		Anulado:    BoolFlex(false),
		Cambio:     FloatFlex(50),
		FechaPago:  "2025-01-01 12:00:00",
		FormasPago: []Pago{
			{
				FormaPago: "Efectivo",
				Cantidad:  FloatFlex(1050),
			},
		},
	}

	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("Failed to marshal DocumentoPago: %v", err)
	}

	var decoded DocumentoPago
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal DocumentoPago: %v", err)
	}

	if decoded.Total != doc.Total {
		t.Errorf("Total = %v, want %v", decoded.Total, doc.Total)
	}
	if len(decoded.FormasPago) != 1 {
		t.Errorf("FormasPago length = %v, want 1", len(decoded.FormasPago))
	}
}

func TestTicketData_PartialFields(t *testing.T) {
	// Test with minimal required fields
	minimalJSON := []byte(`{
		"identificador": "MIN001",
		"total": "100",
		"anulada": "0",
		"conceptos": [],
		"impuestos": [],
		"documentos_pago": [],
		"pago": []
	}`)

	var ticket TicketData
	if err := json.Unmarshal(minimalJSON, &ticket); err != nil {
		t.Fatalf("Failed to unmarshal minimal ticket: %v", err)
	}

	if ticket.Identificador != "MIN001" {
		t.Errorf("Identificador = %v, want MIN001", ticket.Identificador)
	}
	if float64(ticket.Total) != 100 {
		t.Errorf("Total = %v, want 100", ticket.Total)
	}
	if ticket.Anulada {
		t.Error("Anulada should be false")
	}
}
