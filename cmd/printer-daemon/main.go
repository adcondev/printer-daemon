package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

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
		Message: "Connected to printer server",
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
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
