package service

import (
	"fmt"
	"strconv"

	"github.com/adcondev/pos-printer/connector"
	"github.com/adcondev/pos-printer/devices"
	"github.com/adcondev/pos-printer/escpos"
	"github.com/adcondev/pos-printer/escpos/character"
	"github.com/adcondev/pos-printer/profile"

	"github.com/adcondev/printer-daemon/internal/models"
)

// TicketPrinter maneja la impresión de tickets usando la nueva arquitectura
type TicketPrinter struct {
	printer  *devices.Printer
	template *models.TicketTemplateData
	ticket   *models.TicketData
}

// NewTicketPrinter crea una nueva instancia del servicio
func NewTicketPrinter(printerName string) (*TicketPrinter, error) {
	// Crear conexión
	conn, err := connector.NewWindowsPrintConnector(printerName)
	if err != nil {
		return nil, fmt.Errorf("error creando conector: %w", err)
	}

	// Crear perfil (80mm por defecto, se puede parametrizar)
	prof := profile.CreateProfile80mm()
	prof.Model = printerName

	// Crear comandos ESC/POS
	commands := escpos.NewEscposCommands()

	// Crear printer
	printer, err := devices.NewPrinter(commands, prof, conn)
	if err != nil {
		err := conn.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("error creando printer: %w", err)
	}

	return &TicketPrinter{
		printer: printer,
	}, nil
}

// SetTemplate establece la plantilla del ticket
func (tp *TicketPrinter) SetTemplate(template *models.TicketTemplateData) {
	tp.template = template

	// Ajustar perfil según el ancho del ticket
	if tp.template.TicketWidth == 58 {
		prof := profile.CreateProfile58mm()
		tp.printer.Profile = *prof
	}
}

// PrintTicket imprime el ticket completo
func (tp *TicketPrinter) PrintTicket(ticket *models.TicketData) error {
	tp.ticket = ticket

	if tp.template == nil {
		return fmt.Errorf("plantilla no configurada")
	}

	// Inicializar impresora
	if err := tp.printer.Initialize(); err != nil {
		return fmt.Errorf("error inicializando: %w", err)
	}

	// FIXME: Definir en libreria pos printer el conjunto de caracteres a usar, tomar en cuenta el profile
	cmd, _ := tp.printer.Protocol.Character.SelectCharacterCodeTable(character.PC850)
	_, err := tp.printer.Connection.Write(cmd)
	if err != nil {
		return fmt.Errorf("error seleccionando conjunto de caracteres: %w", err)
	}

	// Imprimir secciones del ticket
	if err := tp.printHeader(); err != nil {
		return fmt.Errorf("error en header: %w", err)
	}

	if err := tp.printCustomerInfo(); err != nil {
		return fmt.Errorf("error en info cliente: %w", err)
	}

	if err := tp.printTicketInfo(); err != nil {
		return fmt.Errorf("error en info ticket: %w", err)
	}

	err = tp.printer.Feed(1)
	if err != nil {
		return err
	}

	taxes, err := tp.printItems()
	if err != nil {
		return fmt.Errorf("error en items: %w", err)
	}

	if err := tp.printTaxes(taxes); err != nil {
		return fmt.Errorf("error en impuestos: %w", err)
	}

	if err := tp.printPaymentInfo(); err != nil {
		return fmt.Errorf("error en info pago: %w", err)
	}

	if err := tp.printFooter(); err != nil {
		return fmt.Errorf("error en footer: %w", err)
	}

	// Cortar papel
	err = tp.printer.Feed(3)
	if err != nil {
		return err
	}
	if err := tp.printer.PartialCut(); err != nil {
		return fmt.Errorf("error al cortar: %w", err)
	}

	// IMPORTANTE: Forzar el envío inmediato del trabajo
	if err := tp.forceFlush(); err != nil {
		return fmt.Errorf("error forzando impresión: %w", err)
	}

	return nil
}

// forceFlush fuerza el envío del trabajo de impresión
func (tp *TicketPrinter) forceFlush() error {
	// Guardar el nombre de la impresora
	printerName := tp.printer.Profile.Model

	// Cerrar la conexión actual (esto fuerza el envío)
	if err := tp.printer.Close(); err != nil {
		return err
	}

	// Reabrir la conexión para el siguiente ticket
	conn, err := connector.NewWindowsPrintConnector(printerName)
	if err != nil {
		return fmt.Errorf("error reabriendo conexión: %w", err)
	}

	// Recrear el printer
	commands := escpos.NewEscposCommands()
	tp.printer, err = devices.NewPrinter(commands, &tp.printer.Profile, conn)
	if err != nil {
		err := conn.Close()
		if err != nil {
			return err
		}
		return fmt.Errorf("error recreando printer: %w", err)
	}

	return nil
}

// printHeader imprime el encabezado
func (tp *TicketPrinter) printHeader() error {
	tmpl := tp.template
	ticket := tp.ticket

	// Centrar para el header
	if err := tp.printer.AlignCenter(); err != nil {
		return err
	}

	// Cabecera personalizada
	if tmpl.CambiarCabecera != "" {
		if err := tp.printer.Bold(true); err != nil {
			return err
		}
		if err := tp.printer.PrintLine(tmpl.CambiarCabecera); err != nil {
			return err
		}
		if err := tp.printer.Bold(false); err != nil {
			return err
		}
	}

	// Nombre comercial
	if tmpl.VerNombreC && ticket.SucursalNombreComercial != "" {
		if err := tp.printer.PrintTitle(ticket.SucursalNombreComercial); err != nil {
			return err
		}
	}

	// RFC
	if tmpl.VerRFC && ticket.SucursalRFC != "" {
		if err := tp.printer.PrintLine("RFC: " + ticket.SucursalRFC); err != nil {
			return err
		}
	}

	// Domicilio
	if tmpl.VerDom && ticket.SucursalCalle != "" {
		domicilio := fmt.Sprintf("%s %s, Col. %s",
			ticket.SucursalCalle,
			ticket.SucursalNumero,
			ticket.SucursalColonia)

		if err := tp.printer.PrintLine(domicilio); err != nil {
			return err
		}

		direccion2 := fmt.Sprintf("%s, %s, CP %s",
			ticket.SucursalLocalidad,
			ticket.SucursalEstado,
			ticket.SucursalCP)

		if err := tp.printer.PrintLine(direccion2); err != nil {
			return err
		}
	}

	// Teléfono
	if tmpl.VerTelefono && ticket.SucursalTelefono != "" {
		if err := tp.printer.PrintLine("Tel: " + ticket.SucursalTelefono); err != nil {
			return err
		}
	}

	// Email
	if tmpl.VerEmail && ticket.SucursalEmail != "" {
		if err := tp.printer.PrintLine(ticket.SucursalEmail); err != nil {
			return err
		}
	}

	return tp.printer.PrintSeparator("=", 32)
}

// printCustomerInfo imprime información del cliente
func (tp *TicketPrinter) printCustomerInfo() error {
	if !tp.template.VerNombreCliente || tp.ticket.Cliente == "" {
		return nil
	}

	if err := tp.printer.AlignLeft(); err != nil {
		return err
	}

	return tp.printer.PrintLine("Cliente: " + tp.ticket.Cliente)
}

// printTicketInfo imprime folio, fecha, etc.
func (tp *TicketPrinter) printTicketInfo() error {
	tmpl := tp.template
	ticket := tp.ticket

	if err := tp.printer.AlignLeft(); err != nil {
		return err
	}

	if tmpl.VerFolio && ticket.Folio != "" {
		if err := tp.printer.PrintLine("Folio: " + ticket.Folio); err != nil {
			return err
		}
	}

	if tmpl.VerFecha && ticket.FechaSistema != "" {
		if err := tp.printer.PrintLine("Fecha: " + ticket.FechaSistema); err != nil {
			return err
		}
	}

	if tmpl.VerTienda && ticket.SucursalTienda != "" {
		if err := tp.printer.PrintLine("Tienda: " + ticket.SucursalTienda); err != nil {
			return err
		}
	}

	return nil
}

// printItems imprime los conceptos/productos
func (tp *TicketPrinter) printItems() (map[string]float64, error) {
	// Calcular anchos según papel
	var lenCant, lenDesc, lenPrecio, lenTotal int

	if tp.printer.Profile.PaperWidth == 80 {
		lenCant = 6
		lenDesc = 20
		lenPrecio = 11
		lenTotal = 11
	} else {
		lenCant = 4
		lenDesc = 18
		lenPrecio = 9
		lenTotal = 9
	}

	// Encabezados de columnas
	if err := tp.printer.Bold(true); err != nil {
		return nil, err
	}

	header := PadCenter("CANT", lenCant, ' ') +
		PadCenter("PRODUCTO", lenDesc, ' ') +
		PadCenter("PRECIO", lenPrecio, ' ') +
		PadLeft("TOTAL", lenTotal, ' ')

	if err := tp.printer.PrintLine(header); err != nil {
		return nil, err
	}

	if err := tp.printer.Bold(false); err != nil {
		return nil, err
	}

	// Imprimir cada concepto
	var subtotal float64
	taxes := make(map[string]float64)

	for _, concepto := range tp.ticket.Conceptos {
		cant := PadCenter(strconv.FormatFloat(float64(concepto.Cantidad), 'f', 0, 64), lenCant, ' ')
		desc := PadCenter(Substr(concepto.Descripcion, lenDesc), lenDesc, ' ')
		precio := PadCenter(FormatFloat(float64(concepto.PrecioVenta), 2), lenPrecio, ' ')
		total := PadLeft(FormatFloat(float64(concepto.Total), 2), lenTotal, ' ')

		line := cant + desc + precio + total
		if err := tp.printer.PrintLine(line); err != nil {
			return nil, err
		}

		subtotal += float64(concepto.Total)

		// Procesar impuestos (simplificado)
		for _, imp := range concepto.Impuestos {
			if imp.Tipo == "T" {
				taxes["traslado"] += float64(imp.Importe)
			} else {
				taxes["retencion"] += float64(imp.Importe)
			}
		}
	}

	// Línea separadora
	if err := tp.printer.PrintSeparator("-", 32); err != nil {
		return nil, err
	}

	// Subtotal
	if err := tp.printer.AlignRight(); err != nil {
		return nil, err
	}

	if err := tp.printer.PrintLine(fmt.Sprintf("Subtotal: $%.2f", subtotal)); err != nil {
		return nil, err
	}

	return taxes, nil
}

// printTaxes imprime los impuestos
func (tp *TicketPrinter) printTaxes(taxes map[string]float64) error {
	if !tp.template.VerImpuestos {
		return nil
	}

	if taxes["traslado"] > 0 {
		if err := tp.printer.PrintLine(fmt.Sprintf("IVA: $%.2f", taxes["traslado"])); err != nil {
			return err
		}
	}

	if taxes["retencion"] > 0 {
		if err := tp.printer.PrintLine(fmt.Sprintf("Retención: $%.2f", taxes["retencion"])); err != nil {
			return err
		}
	}

	return nil
}

// printPaymentInfo imprime información de pago
func (tp *TicketPrinter) printPaymentInfo() error {
	// Total
	if err := tp.printer.Bold(true); err != nil {
		return err
	}

	if err := tp.printer.Size(2, 2); err != nil {
		return err
	}

	total := float64(tp.ticket.Total)
	if err := tp.printer.PrintLine(fmt.Sprintf("TOTAL: $%.2f", total)); err != nil {
		return err
	}

	if err := tp.printer.NormalSize(); err != nil {
		return err
	}

	if err := tp.printer.Bold(false); err != nil {
		return err
	}

	// Formas de pago
	for _, pago := range tp.ticket.Pagos {
		cantidad := float64(pago.Cantidad)
		if err := tp.printer.PrintLine(fmt.Sprintf("%s: $%.2f", pago.FormaPago, cantidad)); err != nil {
			return err
		}
	}

	// Cambio
	cambio := float64(tp.ticket.Cambio)
	if cambio > 0 {
		if err := tp.printer.PrintLine(fmt.Sprintf("Cambio: $%.2f", cambio)); err != nil {
			return err
		}
	}

	return nil
}

// printFooter imprime el pie del ticket
func (tp *TicketPrinter) printFooter() error {
	if err := tp.printer.AlignCenter(); err != nil {
		return err
	}

	// Texto pagado
	err := tp.printer.Feed(1)
	if err != nil {
		return err
	}
	if err := tp.printer.PrintTitle("PAGADO"); err != nil {
		return err
	}

	// Leyendas
	if tp.ticket.SucursalLeyenda1 != "" {
		if err := tp.printer.PrintLine(tp.ticket.SucursalLeyenda1); err != nil {
			return err
		}
	}

	// Pie personalizado
	if tp.template.CambiarPie != "" {
		if err := tp.printer.PrintLine(tp.template.CambiarPie); err != nil {
			return err
		}
	}

	return nil
}

// Close cierra la conexión
func (tp *TicketPrinter) Close() error {
	return tp.printer.Close()
}
