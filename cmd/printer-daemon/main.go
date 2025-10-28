// Package principal del servidor WebSocket para el daemon de impresión de tickets POS.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"

	"github.com/adcondev/printer-daemon/internal/models"
	"github.com/adcondev/printer-daemon/internal/service"
)

// Configuración global
var (
	listenAddr = ":8766"

	upgrader = websocket.Upgrader{
		CheckOrigin: func(_ *http.Request) bool {
			return true // Permitir todos los orígenes en desarrollo
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	clients      = make(map[*websocket.Conn]bool)
	clientsMutex sync.Mutex

	// Variable global para el servicio de impresión
	ticketPrinter *service.TicketPrinter
)

// Message representa el mensaje recibido del cliente WebSocket
type Message struct {
	Tipo  string          `json:"tipo"`  // "config", "template", "ticket"
	Datos json.RawMessage `json:"datos"` // Datos JSON sin procesar
}

// Response representa la respuesta enviada al cliente
type Response struct {
	Tipo    string `json:"tipo"` // "ack" o "info"
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// printSeparator imprime un separador visual en la consola
func printSeparator(char string, _ int) {
	fmt.Println(strings.Repeat(char, 80))
}

// printHeader imprime un encabezado para los mensajes
func printHeader(messageType string) {
	printSeparator("=", 80)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("📦 MENSAJE RECIBIDO: %s [%s]\n", strings.ToUpper(messageType), timestamp)
	printSeparator("=", 80)
}

// processConfig procesa y muestra los datos de configuración
func processConfig(rawData json.RawMessage) error {
	printHeader("CONFIG")

	// Mostrar JSON crudo
	fmt.Println("\n🔸 JSON CRUDO:")
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, rawData, "", "  "); err != nil {
		return fmt.Errorf("error formateando JSON: %w", err)
	}
	fmt.Println(prettyJSON.String())

	// Parsear usando el modelo
	config, err := models.BytesToConfig(rawData)
	if err != nil {
		return fmt.Errorf("error parseando config: %w", err)
	}

	// Inicializar servicio de impresión con el nombre de la impresora
	if ticketPrinter != nil {
		err := ticketPrinter.Close()
		if err != nil {
			return err
		} // Cerrar conexión anterior si existe
	}

	ticketPrinter, err = service.NewTicketPrinter(config.Printer)
	if err != nil {
		return fmt.Errorf("error inicializando impresora: %w", err)
	}

	log.Printf("✅ Impresora %s configurada correctamente", config.Printer)

	// Mostrar datos parseados
	fmt.Println("\n🔹 DATOS PARSEADOS:")
	fmt.Printf("  Impresora:    %s\n", config.Printer)
	fmt.Printf("  Debug Log:    %v\n", config.DebugLog)

	printSeparator("-", 80)

	return nil
}

// processTemplate procesa y muestra los datos de plantilla
func processTemplate(rawData json.RawMessage) error {
	printHeader("TEMPLATE")

	// Mostrar JSON crudo
	fmt.Println("\n🔸 JSON CRUDO:")
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, rawData, "", "  "); err != nil {
		return fmt.Errorf("error formateando JSON: %w", err)
	}
	fmt.Println(prettyJSON.String())

	// Parsear usando el modelo
	template, err := models.BytesToTicketTemplate(rawData)
	if err != nil {
		return fmt.Errorf("error parseando template: %w", err)
	}

	// Configurar plantilla en el servicio
	if ticketPrinter != nil {
		ticketPrinter.SetTemplate(template)
		log.Println("✅ Plantilla configurada en el servicio de impresión")
	} else {
		return fmt.Errorf("servicio de impresión no inicializado")
	}

	// Mostrar datos parseados
	fmt.Println("\n🔹 DATOS PARSEADOS:")
	fmt.Printf("  Ancho Ticket:         %d caracteres\n", template.TicketWidth)
	fmt.Printf("  Tamaño Razón Social:  %d\n", template.RazonSocialSize)
	fmt.Printf("  Tamaño Datos:         %d\n", template.DatosSize)
	fmt.Printf("  Ancho Logo:           %d\n", template.LogoWidth)
	fmt.Println("\n  Elementos Visibles:")
	fmt.Printf("    - Logotipo:           %v\n", template.VerLogotipo)
	fmt.Printf("    - Nombre:             %v\n", template.VerNombre)
	fmt.Printf("    - Nombre Comercial:   %v\n", template.VerNombreC)
	fmt.Printf("    - RFC:                %v\n", template.VerRFC)
	fmt.Printf("    - Domicilio:          %v\n", template.VerDom)
	fmt.Printf("    - Régimen:            %v\n", template.VerRegimen)
	fmt.Printf("    - Email:              %v\n", template.VerEmail)
	fmt.Printf("    - Teléfono:           %v\n", template.VerTelefono)
	fmt.Printf("    - Nombre Cliente:     %v\n", template.VerNombreCliente)
	fmt.Printf("    - Folio:              %v\n", template.VerFolio)
	fmt.Printf("    - Fecha:              %v\n", template.VerFecha)
	fmt.Printf("    - Tienda:             %v\n", template.VerTienda)
	fmt.Printf("    - Precio Unitario:    %v\n", template.VerPrecioU)
	fmt.Printf("    - Cant. Productos:    %v\n", template.VerCantProductos)
	fmt.Printf("    - Impuestos:          %v\n", template.VerImpuestos)
	fmt.Printf("    - Total Impuestos:    %v\n", template.VerImpuestosTotal)
	fmt.Printf("    - Series:             %v\n", template.VerSeries)
	fmt.Println("\n  Textos Personalizados:")
	fmt.Printf("    - Cabecera:     %q\n", template.CambiarCabecera)
	fmt.Printf("    - Reclamación:  %q\n", template.CambiarReclamacion)
	fmt.Printf("    - Pie:          %q\n", template.CambiarPie)

	printSeparator("-", 80)

	return nil
}

// processTicket procesa y muestra los datos del ticket
func processTicket(rawData json.RawMessage) error {
	printHeader("TICKET")

	// Mostrar JSON crudo
	fmt.Println("\n🔸 JSON CRUDO:")
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, rawData, "", "  "); err != nil {
		return fmt.Errorf("error formateando JSON: %w", err)
	}
	fmt.Println(prettyJSON.String())

	// Parsear usando el modelo
	ticket, err := models.BytesToTicket(rawData)
	if err != nil {
		return fmt.Errorf("error parseando ticket: %w", err)
	}

	// IMPRIMIR EL TICKET
	if ticketPrinter != nil {
		log.Println("🖨️ Enviando ticket a la impresora...")

		if err := ticketPrinter.PrintTicket(ticket); err != nil {
			return fmt.Errorf("error imprimiendo ticket: %w", err)
		}

		log.Println("✅ Ticket impreso correctamente")
	} else {
		return fmt.Errorf("servicio de impresión no inicializado")
	}

	// Mostrar datos parseados
	fmt.Println("\n🔹 DATOS PARSEADOS:")
	fmt.Printf("  Identificador:  %s\n", ticket.Identificador)
	fmt.Printf("  Folio:          %s\n", ticket.Folio)
	fmt.Printf("  Serie:          %s\n", ticket.Serie)
	fmt.Printf("  Fecha:          %s\n", ticket.FechaSistema)
	fmt.Printf("  Vendedor:       %s\n", ticket.Vendedor)
	fmt.Printf("  Tipo Operación: %s\n", ticket.TipoOperacion)
	fmt.Printf("  Anulada:        %v\n", ticket.Anulada)

	fmt.Println("\n  Cliente:")
	fmt.Printf("    Nombre:       %s\n", ticket.Cliente)
	fmt.Printf("    RFC:          %s\n", ticket.ClienteRFC)
	fmt.Printf("    CP:           %s\n", ticket.ClienteCP)
	fmt.Printf("    Email:        %s\n", ticket.ClienteEmail)

	fmt.Println("\n  Sucursal:")
	fmt.Printf("    Nombre:       %s\n", ticket.SucursalNombre)
	fmt.Printf("    Comercial:    %s\n", ticket.SucursalNombreComercial)
	fmt.Printf("    RFC:          %s\n", ticket.SucursalRFC)
	fmt.Printf("    Tienda:       %s\n", ticket.SucursalTienda)
	fmt.Printf("    Teléfono:     %s\n", ticket.SucursalTelefono)
	fmt.Printf("    Email:        %s\n", ticket.SucursalEmail)

	fmt.Println("\n  Montos:")
	fmt.Printf("    Total:        $%.2f\n", float64(ticket.Total))
	fmt.Printf("    Descuento:    $%.2f\n", float64(ticket.Descuento))
	fmt.Printf("    Pagado:       $%.2f\n", float64(ticket.Pagado))
	fmt.Printf("    Cambio:       $%.2f\n", float64(ticket.Cambio))
	fmt.Printf("    Saldo:        $%.2f\n", float64(ticket.Saldo))

	fmt.Printf("\n  Conceptos (%d items):\n", len(ticket.Conceptos))
	for i, concepto := range ticket.Conceptos {
		fmt.Printf("    [%d] %s\n", i+1, concepto.Descripcion)
		fmt.Printf("        Clave:    %s\n", concepto.Clave)
		fmt.Printf("        Cantidad: %.2f %s\n", float64(concepto.Cantidad), concepto.Unidad)
		fmt.Printf("        Precio:   $%.2f\n", float64(concepto.PrecioVenta))
		fmt.Printf("        Total:    $%.2f\n", float64(concepto.Total))
		if len(concepto.Series) > 0 {
			fmt.Printf("        Series:   %v\n", concepto.Series)
		}
		if len(concepto.Impuestos) > 0 {
			fmt.Printf("        Impuestos: %d\n", len(concepto.Impuestos))
		}
	}

	fmt.Printf("\n  Formas de Pago (%d):\n", len(ticket.Pagos))
	for i, pago := range ticket.Pagos {
		fmt.Printf("    [%d] %s: $%.2f\n", i+1, pago.FormaPago, float64(pago.Cantidad))
	}

	if len(ticket.DocumentosPago) > 0 {
		fmt.Printf("\n  Documentos de Pago: %d\n", len(ticket.DocumentosPago))
	}

	fmt.Println("\n  Leyendas:")
	fmt.Printf("    1: %s\n", ticket.SucursalLeyenda1)
	fmt.Printf("    2: %s\n", ticket.SucursalLeyenda2)

	if ticket.AutofacturaLink != "" {
		fmt.Printf("\n  Autofactura: %s\n", ticket.AutofacturaLink)
	}

	printSeparator("-", 80)

	return nil
}

// handleWebSocket maneja la conexión WebSocket
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to upgrade connection: %v", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("[ERROR] Failed to close connection: %v", err)
		}
	}(conn)

	// Registrar cliente
	clientsMutex.Lock()
	clients[conn] = true
	clientCount := len(clients)
	clientsMutex.Unlock()

	log.Printf("[CLIENT] ✅ Nueva conexión establecida (total: %d)", clientCount)

	// Enviar mensaje de bienvenida
	welcome := Response{
		Tipo:    "info",
		Success: true,
		Message: "Conectado al servidor de impresión",
	}
	if err := conn.WriteJSON(welcome); err != nil {
		log.Printf("[ERROR] Failed to send welcome message: %v", err)
		return
	}

	// Leer mensajes
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[ERROR] WebSocket error: %v", err)
			}
			break
		}

		log.Printf("[CLIENT] 📨 Mensaje recibido: tipo=%s, tamaño=%d bytes", msg.Tipo, len(msg.Datos))

		// Procesar según el tipo
		var processErr error
		switch msg.Tipo {
		case "config":
			processErr = processConfig(msg.Datos)
		case "template":
			processErr = processTemplate(msg.Datos)
		case "ticket":
			processErr = processTicket(msg.Datos)
		default:
			processErr = fmt.Errorf("tipo de mensaje desconocido: %s", msg.Tipo)
		}

		// Enviar respuesta
		response := Response{
			Tipo: "ack",
		}

		if processErr != nil {
			response.Success = false
			response.Message = fmt.Sprintf("Error procesando %s: %v", msg.Tipo, processErr)
			log.Printf("[ERROR] ❌ %s", response.Message)
		} else {
			response.Success = true
			response.Message = fmt.Sprintf("✅ %s recibido y procesado correctamente", strings.ToUpper(msg.Tipo))
			log.Printf("[SUCCESS] %s", response.Message)
		}

		if err := conn.WriteJSON(response); err != nil {
			log.Printf("[ERROR] Failed to send response: %v", err)
			break
		}
	}

	// Desregistrar cliente
	clientsMutex.Lock()
	delete(clients, conn)
	clientCount = len(clients)
	clientsMutex.Unlock()

	log.Printf("[CLIENT] 👋 Cliente desconectado (restantes: %d)", clientCount)
}

func main() {
	log.Println("╔════════════════════════════════════════════════════════════╗")
	log.Println("║   PRINTER DAEMON - WebSocket Server                       ║")
	log.Println("║   Servidor de Impresión de Tickets POS                    ║")
	log.Println("╚════════════════════════════════════════════════════════════╝")
	log.Printf("\n🚀 Iniciando servidor en %s\n", listenAddr)

	// Agregar defer para cerrar la impresora
	defer func() {
		if ticketPrinter != nil {
			err := ticketPrinter.Close()
			if err != nil {
				return
			}
		}
	}()

	// Configurar servidor de archivos estáticos
	fileServer := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Ruta principal para el HTML
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "web/index.html")
	})

	// Endpoint WebSocket
	http.HandleFunc("/ws", handleWebSocket)

	server := &http.Server{
		Addr:         listenAddr,
		Handler:      nil,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Iniciar servidor en goroutine
	go func() {
		log.Printf("📡 Servidor HTTP iniciado en http://localhost%s", listenAddr)
		log.Printf("🌐 WebSocket disponible en ws://localhost%s/ws", listenAddr)
		log.Println("\n💡 Abre tu navegador en: http://192.168.8.82" + listenAddr)
		log.Println("⚡ Presiona Ctrl+C para detener el servidor")
		printSeparator("=", 80)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[FATAL] Error en ListenAndServe: %v", err)
		}
	}()

	// Esperar señal de interrupción
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("\n\n🛑 Señal de apagado recibida...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("[ERROR] Error durante el apagado: %v", err)
	}

	// Cerrar clientes WebSocket
	clientsMutex.Lock()
	for conn := range clients {
		_ = conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Servidor apagándose"),
		)
		_ = conn.Close()
	}
	clientsMutex.Unlock()

	log.Println("✅ Servidor detenido correctamente")
	log.Println("👋 ¡Hasta luego!")
}
