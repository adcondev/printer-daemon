package devices

import (
	"log"

	"github.com/adcondev/pos-printer/connector"
	"github.com/adcondev/pos-printer/escpos"
	"github.com/adcondev/pos-printer/pos"
	"github.com/adcondev/pos-printer/profile"
)

func NewECPM80250() (*pos.EscposPrinter, error) {
	// Crear conector de impresora
	// Crear perfil de impresora
	prof := profile.CreateECPM80250()
	log.Printf("Intentando conectar a la impresora: %s", prof.Model)
	conn, err := connector.NewWindowsPrintConnector(prof.Model)
	if err != nil {
		log.Fatalf("Error al crear el conector: %v", err)
		return nil, err
	}

	// Crear instancia de impresora genérica
	printer, err := pos.NewPrinter(pos.EscposProto, conn, prof)
	if err != nil {
		log.Printf("Error al crear la impresora: %v", err)
		return nil, err
	}
	return printer, err
}

func PrintFromWs(content string, cut bool, device *pos.EscposPrinter) {
	defer func(printer *pos.EscposPrinter) {
		err := printer.Close()
		if err != nil {
			log.Printf("Error al cerrar la impresora: %v", err)
		}
	}(device)

	// Inicializar impresora
	log.Println("Enviando comandos de prueba...")
	if err := device.Initialize(); err != nil {
		log.Printf("Error al inicializar: %v", err)
	}
	// Imprimir documento completo
	if err := device.Print(content); err != nil {
		log.Fatalf("Error al imprimir contenido: %v", err)
	}
	if err := device.Feed(3); err != nil {
		log.Fatalf("Error al alimentar papel: %v", err)
	}
	if cut {
		if err := device.Cut(escpos.PartialCut); err != nil {
			log.Fatalf("Error al cortar papel: %v", err)
		}
	}
	if err := device.Feed(3); err != nil {
		log.Fatalf("Error al alimentar papel: %v", err)
	}

	log.Println("Impresión completada!")
}
