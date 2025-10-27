package models

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRealFiles_Integration(t *testing.T) {
	// Skip if test files don't exist
	if _, err := os.Stat("./config_ticket.json"); os.IsNotExist(err) {
		t.Skip("Skipping integration test: test files not found")
	}

	t.Run("config file", func(t *testing.T) {
		data, err := os.ReadFile("./config_ticket.json")
		if err != nil {
			t.Fatalf("Failed to read config file: %v", err)
		}

		config, err := BytesToConfig(data)
		if err != nil {
			t.Fatalf("BytesToConfig() error = %v", err)
		}

		if config.Printer == "" {
			t.Error("Config printer should not be empty")
		}

		// Test roundtrip
		parser := ConfigParser{}
		marshaled, err := parser.ToBytes(config)
		if err != nil {
			t.Fatalf("ToBytes() error = %v", err)
		}

		var roundtrip ConfigData
		if err := json.Unmarshal(marshaled, &struct {
			Data *ConfigData `json:"data"`
		}{Data: &roundtrip}); err != nil {
			t.Fatalf("Roundtrip unmarshal error = %v", err)
		}

		if roundtrip.Printer != config.Printer {
			t.Errorf("Roundtrip printer = %v, want %v", roundtrip.Printer, config.Printer)
		}
	})

	t.Run("template file", func(t *testing.T) {
		data, err := os.ReadFile("./template_ticket.json")
		if err != nil {
			t.Fatalf("Failed to read template file: %v", err)
		}

		template, err := BytesToTicketTemplate(data)
		if err != nil {
			t.Fatalf("BytesToTicketTemplate() error = %v", err)
		}

		// Verify some flex type fields
		if template.TicketWidth != 80 {
			t.Errorf("TicketWidth = %v, want 80", template.TicketWidth)
		}
		if !template.VerLogotipo {
			t.Error("VerLogotipo should be true")
		}

		// Test roundtrip
		parser := TemplateParser{}
		marshaled, err := parser.ToBytes(template)
		if err != nil {
			t.Fatalf("ToBytes() error = %v", err)
		}

		roundtrip, err := parser.Parse(marshaled)
		if err != nil {
			t.Fatalf("Parse() roundtrip error = %v", err)
		}

		if roundtrip.TicketWidth != template.TicketWidth {
			t.Errorf("Roundtrip TicketWidth = %v, want %v",
				roundtrip.TicketWidth, template.TicketWidth)
		}
	})

	t.Run("ticket file", func(t *testing.T) {
		data, err := os.ReadFile("./data_ticket.json")
		if err != nil {
			t.Fatalf("Failed to read ticket file: %v", err)
		}

		ticket, err := BytesToTicket(data)
		if err != nil {
			t.Fatalf("BytesToTicket() error = %v", err)
		}

		// Verify basic fields
		if ticket.Identificador == "" {
			t.Error("Ticket identificador should not be empty")
		}
		if len(ticket.Conceptos) == 0 {
			t.Error("Ticket should have conceptos")
		}

		// Test flex types in conceptos
		for i, concepto := range ticket.Conceptos {
			if concepto.Clave == "" {
				t.Errorf("Concepto[%d] clave should not be empty", i)
			}
			if concepto.Total == 0 {
				t.Errorf("Concepto[%d] total should not be zero", i)
			}
		}

		// Test roundtrip
		parser := TicketParser{}
		marshaled, err := parser.ToBytes(ticket)
		if err != nil {
			t.Fatalf("ToBytes() error = %v", err)
		}

		roundtrip, err := parser.Parse(marshaled)
		if err != nil {
			t.Fatalf("Parse() roundtrip error = %v", err)
		}

		if roundtrip.Identificador != ticket.Identificador {
			t.Errorf("Roundtrip identificador = %v, want %v",
				roundtrip.Identificador, ticket.Identificador)
		}
		if len(roundtrip.Conceptos) != len(ticket.Conceptos) {
			t.Errorf("Roundtrip conceptos length = %v, want %v",
				len(roundtrip.Conceptos), len(ticket.Conceptos))
		}
	})
}

func TestFileReader_Security(t *testing.T) {
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs")
	err := os.MkdirAll(docsDir, 0755)
	if err != nil {
		return
	}

	// Create a large file
	largeFile := filepath.Join(docsDir, "large.json")
	largeData := make([]byte, MaxFileSize+1)
	for i := range largeData {
		largeData[i] = 'a'
	}
	err = os.WriteFile(largeFile, largeData, 0644)
	if err != nil {
		return
	}

	reader := &DefaultFileReader{
		AllowedDir: docsDir,
		MaxSize:    MaxFileSize,
	}

	t.Run("reject large file", func(t *testing.T) {
		_, err := reader.ReadJSONFile(largeFile)
		if err == nil {
			t.Error("Expected error for large file")
		}
	})

	t.Run("reject path traversal", func(t *testing.T) {
		maliciousPath := filepath.Join(docsDir, "..", "..", "etc", "passwd.json")
		_, err := reader.ReadJSONFile(maliciousPath)
		if err == nil {
			t.Error("Expected error for path traversal")
		}
	})
}

func TestCompatibilityFunctions(t *testing.T) {
	ticketJSON := []byte(`{
		"data": {
			"identificador": "TEST",
			"total": "100",
			"conceptos": [],
			"impuestos": [],
			"documentos_pago": [],
			"pago": []
		}
	}`)

	t.Run("BytesToNewTicket", func(t *testing.T) {
		ticket1, err1 := BytesToTicket(ticketJSON)
		ticket2, err2 := BytesToNewTicket(ticketJSON)

		if err1 != nil || err2 != nil {
			t.Fatalf("Unexpected errors: %v, %v", err1, err2)
		}

		if ticket1.Identificador != ticket2.Identificador {
			t.Error("BytesToNewTicket should produce same result as BytesToTicket")
		}
	})
}
