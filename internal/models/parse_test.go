package models

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultFileReader_ReadJSONFile(t *testing.T) {
	// Create temp directory for test files
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs")
	err := os.MkdirAll(docsDir, 0755)
	if err != nil {
		return
	}

	// Create test files
	validJSON := filepath.Join(docsDir, "valid.json")
	err = os.WriteFile(validJSON, []byte(`{"test": "data"}`), 0644)
	if err != nil {
		return
	}

	invalidJSON := filepath.Join(docsDir, "invalid.json")
	err = os.WriteFile(invalidJSON, []byte(`{invalid json}`), 0644)
	if err != nil {
		return
	}

	notJSON := filepath.Join(docsDir, "test.txt")
	err = os.WriteFile(notJSON, []byte(`text file`), 0644)
	if err != nil {
		return
	}

	tests := []struct {
		name       string
		filepath   string
		allowedDir string
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "valid JSON file",
			filepath:   validJSON,
			allowedDir: docsDir,
			wantErr:    false,
		},
		{
			name:       "invalid JSON content",
			filepath:   invalidJSON,
			allowedDir: docsDir,
			wantErr:    true,
			errMsg:     "archivo JSON inválido",
		},
		{
			name:       "non-JSON extension",
			filepath:   notJSON,
			allowedDir: docsDir,
			wantErr:    true,
			errMsg:     "solo se permiten archivos JSON",
		},
		{
			name:       "file outside allowed directory",
			filepath:   "/etc/passwd",
			allowedDir: docsDir,
			wantErr:    true,
			errMsg:     "acceso denegado",
		},
		{
			name:       "non-existent file",
			filepath:   filepath.Join(docsDir, "nonexistent.json"),
			allowedDir: docsDir,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &DefaultFileReader{
				AllowedDir: tt.allowedDir,
				MaxSize:    MaxFileSize,
			}
			_, err := reader.ReadJSONFile(tt.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadJSONFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigParser(t *testing.T) {
	parser := ConfigParser{}

	t.Run("parse valid config", func(t *testing.T) {
		input := []byte(`{
			"data": {
				"printer": "TestPrinter",
				"debug_log": true
			}
		}`)

		config, err := parser.Parse(input)
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		if config.Printer != "TestPrinter" {
			t.Errorf("Parse() Printer = %v, want TestPrinter", config.Printer)
		}
		if !config.DebugLog {
			t.Errorf("Parse() DebugLog = %v, want true", config.DebugLog)
		}
	})

	t.Run("marshal config", func(t *testing.T) {
		config := &ConfigData{
			Printer:  "TestPrinter",
			DebugLog: true,
		}

		data, err := parser.ToBytes(config)
		if err != nil {
			t.Fatalf("ToBytes() error = %v", err)
		}

		// Parse back to verify
		parsed, err := parser.Parse(data)
		if err != nil {
			t.Fatalf("Parse() after ToBytes() error = %v", err)
		}

		if parsed.Printer != config.Printer {
			t.Errorf("Roundtrip Printer = %v, want %v", parsed.Printer, config.Printer)
		}
	})
}

func TestTemplateParser(t *testing.T) {
	parser := TemplateParser{}

	t.Run("parse template with flex types", func(t *testing.T) {
		input := []byte(`{
			"data": {
				"ticket_width": "80",
				"ver_logotipo": "1",
				"ver_nombre": false,
				"logo_width": 120
			}
		}`)

		template, err := parser.Parse(input)
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		if template.TicketWidth != 80 {
			t.Errorf("Parse() TicketWidth = %v, want 80", template.TicketWidth)
		}
		if !template.VerLogotipo {
			t.Errorf("Parse() VerLogotipo = %v, want true", template.VerLogotipo)
		}
		if template.VerNombre {
			t.Errorf("Parse() VerNombre = %v, want false", template.VerNombre)
		}
		if template.LogoWidth != 120 {
			t.Errorf("Parse() LogoWidth = %v, want 120", template.LogoWidth)
		}
	})
}

func TestTicketParser(t *testing.T) {
	parser := TicketParser{}

	t.Run("parse ticket with conceptos", func(t *testing.T) {
		input := []byte(`{
			"data": {
				"identificador": "TEST123",
				"total": "1234.56",
				"anulada": "0",
				"conceptos": [
					{
						"clave": "PROD001",
						"descripcion": "Producto Test",
						"cantidad": "2",
						"precio_venta": "100.50",
						"total": 201.00,
						"unidad": "PIEZA",
						"clave_producto_servicio": "12345678",
						"clave_unidad_sat": "H87",
						"venta_granel": false,
						"series": ["SERIE1", "SERIE2"],
						"impuestos": []
					}
				],
				"impuestos": [],
				"documentos_pago": [],
				"pago": []
			}
		}`)

		ticket, err := parser.Parse(input)
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		if ticket.Identificador != "TEST123" {
			t.Errorf("Parse() Identificador = %v, want TEST123", ticket.Identificador)
		}
		if float64(ticket.Total) != 1234.56 {
			t.Errorf("Parse() Total = %v, want 1234.56", ticket.Total)
		}
		if ticket.Anulada {
			t.Errorf("Parse() Anulada = %v, want false", ticket.Anulada)
		}
		if len(ticket.Conceptos) != 1 {
			t.Fatalf("Parse() Conceptos length = %v, want 1", len(ticket.Conceptos))
		}

		concepto := ticket.Conceptos[0]
		if concepto.Clave != "PROD001" {
			t.Errorf("Parse() Concepto.Clave = %v, want PROD001", concepto.Clave)
		}
		if float64(concepto.Cantidad) != 2 {
			t.Errorf("Parse() Concepto.Cantidad = %v, want 2", concepto.Cantidad)
		}
		if len(concepto.Series) != 2 {
			t.Errorf("Parse() Concepto.Series length = %v, want 2", len(concepto.Series))
		}
	})
}
