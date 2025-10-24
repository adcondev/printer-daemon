package models

import (
	"fmt"
	"regexp"
	"strings"
)

// TODO: Use lazy initialization with sync.Once if it's rarely used to avoid unnecessary compilation overhead at startup.

var (
	// RFCPattern patrón para validar RFC mexicano
	RFCPattern = regexp.MustCompile(`^[A-ZÑ&]{3,4}[0-9]{6}[A-Z0-9]{3}$`)
	// CPPattern patrón para validar código postal
	CPPattern = regexp.MustCompile(`^[0-9]{5}$`)
)

// Validate valida los datos del ticket
func (t *TicketData) Validate() error {
	if t.Identificador == "" {
		return fmt.Errorf("identificador es requerido")
	}

	if t.Folio == "" {
		return fmt.Errorf("folio es requerido")
	}

	// Validar RFC si no es público general
	if t.ClienteData.ClienteRFC != "XAXX010101000" && !RFCPattern.MatchString(t.ClienteData.ClienteRFC) {
		return fmt.Errorf("RFC del cliente inválido: %s", t.ClienteData.ClienteRFC)
	}

	// Validar conceptos
	if len(t.Conceptos) == 0 {
		return fmt.Errorf("debe incluir al menos un concepto")
	}

	for i, concepto := range t.Conceptos {
		if err := concepto.Validate(); err != nil {
			return fmt.Errorf("error en concepto %d: %w", i+1, err)
		}
	}

	return nil
}

// Validate valida un concepto
func (c *Concepto) Validate() error {
	if c.Descripcion == "" {
		return fmt.Errorf("descripción es requerida")
	}

	if float64(c.Cantidad) <= 0 {
		return fmt.Errorf("cantidad debe ser mayor a cero")
	}

	if float64(c.PrecioVenta) < 0 {
		return fmt.Errorf("precio no puede ser negativo")
	}

	return nil
}

// Validate valida la configuración
func (c *ConfigData) Validate() error {
	if c.Printer == "" {
		return fmt.Errorf("impresora es requerida")
	}

	return nil
}

// Validate valida el template
func (t *TicketTemplateData) Validate() error {
	if int(t.TicketWidth) <= 0 {
		return fmt.Errorf("ancho del ticket debe ser mayor a cero")
	}

	// Validar alineación del logo si existe
	if t.Logo.Path != "" {
		alignment := strings.ToLower(t.Logo.Alignment)
		if alignment != "" && alignment != "left" && alignment != "center" && alignment != "right" {
			return fmt.Errorf("alineación del logo inválida: %s", t.Logo.Alignment)
		}
	}

	return nil
}
