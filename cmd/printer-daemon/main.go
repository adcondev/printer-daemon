package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/adcondev/printer-daemon/internal/devices"

	"github.com/gorilla/websocket"
)

// Configuración global
var (
	// Configuración del servidor
	listenAddr = ":8766"

	// Configuración de WebSocket
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Permitir todos los orígenes en desarrollo
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// Gestión de clientes
	clients      = make(map[*websocket.Conn]bool)
	clientsMutex sync.Mutex
)

// PrintJob representa un trabajo de impresión recibido a través de WebSocket
type PrintJob struct {
	Tipo   string `json:"tipo"`   // Tipo "print"
	Texto  string `json:"texto"`  // Texto a imprimir
	Cortar bool   `json:"cortar"` // Cortar después de imprimir
}

// Response representa la respuesta enviada al cliente
type Response struct {
	Tipo    string `json:"tipo"` // "ack" o "error"
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Registros en consola en esta versión simplificada
func processPrintJob(job PrintJob) error {
	log.Printf("[PRINT] Starting print job (cut=%v)", job.Cortar)

	// Imprime en la consola en lugar de una impresora real
	fmt.Println("\n===== PRINT JOB START =====")
	fmt.Println(job.Texto)

	if job.Cortar {
		fmt.Println("--------- CUT HERE ---------")
	}

	fmt.Print("===== PRINT JOB END =====\n\n")

	return nil
}

// Maneja el cliente WebSocket
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

	log.Printf("[CLIENT] New connection (total: %d)", clientCount)

	// Enviar mensaje de bienvenida
	welcome := Response{
		Tipo:    "info",
		Success: true,
		Message: "Connected to devices server",
	}
	err = conn.WriteJSON(welcome)
	if err != nil {
		log.Printf("[ERROR] Failed to send welcome message: %v", err)
		return
	}

	// Leer mensajes
	for {
		var job PrintJob
		err := conn.ReadJSON(&job)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[ERROR] WebSocket error: %v", err)
			}
			break
		}

		log.Printf("[CLIENT] Received job: tipo=%s, cortar=%v, len(texto)=%d",
			job.Tipo, job.Cortar, len(job.Texto))

		if job.Tipo != "print" {
			continue
		}

		// Procesar trabajo de impresión
		startTime := time.Now()
		log.Println("[INIT] Starting Printer Daemon")
		printer, err := devices.NewECPM80250()
		if err != nil {
			log.Printf("[ERROR] Failed to initialize devices: %v", err)
		}
		devices.PrintFromWs(job.Texto, job.Cortar, printer)
		err = processPrintJob(job)
		duration := time.Since(startTime)

		// Enviar respuesta
		response := Response{
			Tipo: "ack",
		}

		if err != nil {
			response.Success = false
			response.Message = fmt.Sprintf("Print failed: %v", err)
			log.Printf("[ERROR] Print job failed: %v", err)
		} else {
			response.Success = true
			response.Message = fmt.Sprintf("Printed successfully in %v", duration)
		}

		if err := conn.WriteJSON(response); err != nil {
			log.Printf("[ERROR] Failed to send response: %v", err)
			break
		}
	}

	// Anular registro del cliente
	clientsMutex.Lock()
	delete(clients, conn)
	clientCount = len(clients)
	clientsMutex.Unlock()

	log.Printf("[CLIENT] Disconnected (remaining: %d)", clientCount)
}

func main() {
	log.Printf("[START] Simple Text Printer WebSocket Server")
	log.Printf("[CONFIG] Listening on %s", listenAddr)

	// Configurar el servidor de archivos estáticos
	fileServer := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Ruta principal para el HTML
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Servir el HTML
		http.ServeFile(w, r, "web/index.html")
	})

	http.HandleFunc("/ws", handleWebSocket)

	log.Printf("[START] Server listening on %s", listenAddr)
	log.Printf("[INFO] Open http://localhost%s in your browser", listenAddr)

	server := &http.Server{
		Addr:    listenAddr,
		Handler: nil, // usa DefaultServeMux
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[FATAL] ListenAndServe: %v", err)
		}
	}()

	// esperar señal Ctrl+C
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Printf("[SHUTDOWN] Signal received, shutting down...")

	// contexto con timeout para graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("[ERROR] Server Shutdown Failed: %v", err)
	} else {
		log.Printf("[SHUTDOWN] HTTP server stopped")
	}

	// cerrar clientes WebSocket activos
	clientsMutex.Lock()
	for conn := range clients {
		_ = conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "server shutting down"),
		)
		_ = conn.Close()
		delete(clients, conn)
	}
	clientsMutex.Unlock()

	log.Printf("[SHUTDOWN] All websocket clients closed, exiting")
}
