package models

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	fp "path/filepath"
	"strings"
)

const (
	// MaxFileSize límite de tamaño del archivo JSON (10MB)
	MaxFileSize = 10 * 1024 * 1024
	// DefaultAllowedDir directorio por defecto permitido
	DefaultAllowedDir = "./docs"
)

// FileReader interfaz para permitir testing
type FileReader interface {
	ReadJSONFile(filepath string) ([]byte, error)
}

// DefaultFileReader implementación por defecto
type DefaultFileReader struct {
	AllowedDir string
	MaxSize    int64
}

// NewFileReader crea un nuevo lector de archivos con configuración por defecto
func NewFileReader() *DefaultFileReader {
	return &DefaultFileReader{
		AllowedDir: DefaultAllowedDir,
		MaxSize:    MaxFileSize,
	}
}

// ReadJSONFile lee un archivo JSON con validaciones de seguridad
func (r *DefaultFileReader) ReadJSONFile(filepath string) ([]byte, error) {
	// Normalizar la ruta
	filepath = fp.Clean(filepath)

	// Validar extensión
	if !strings.HasSuffix(strings.ToLower(filepath), ".json") {
		return nil, fmt.Errorf("solo se permiten archivos JSON")
	}

	// Validar que esté dentro del directorio permitido
	absPath, err := fp.Abs(filepath)
	if err != nil {
		return nil, fmt.Errorf("error al resolver la ruta: %w", err)
	}

	absPath, err = fp.EvalSymlinks(absPath)
	if err != nil {
		return nil, fmt.Errorf("error al evaluar enlaces simbólicos: %w", err)
	}

	absAllowedDir, err := fp.Abs(r.AllowedDir)
	if err != nil {
		return nil, fmt.Errorf("error al resolver directorio permitido: %w", err)
	}

	absAllowedDir, err = fp.EvalSymlinks(absAllowedDir)
	if err != nil {
		return nil, fmt.Errorf("error al evaluar directorio permitido: %w", err)
	}

	if !strings.HasPrefix(absPath, absAllowedDir) {
		return nil, fmt.Errorf("acceso denegado: archivo fuera del directorio permitido")
	}

	// Verificar el tamaño antes de leer
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return nil, fmt.Errorf("error al obtener información del archivo: %w", err)
	}

	if fileInfo.Size() > r.MaxSize {
		return nil, fmt.Errorf("archivo muy grande: %d bytes (máximo: %d)", fileInfo.Size(), r.MaxSize)
	}

	// Abrir y leer el archivo
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir archivo: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error al cerrar archivo JSON: %v", err)
		}
	}()

	// Usar LimitReader como medida adicional de seguridad
	limitedReader := io.LimitReader(file, r.MaxSize)
	content, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("error al leer archivo: %w", err)
	}

	// Validar que sea JSON válido
	var js json.RawMessage
	if err := json.Unmarshal(content, &js); err != nil {
		return nil, fmt.Errorf("archivo JSON inválido: %w", err)
	}

	return content, nil
}

// JSONFileToBytes función de compatibilidad
func JSONFileToBytes(filepath string) ([]byte, error) {
	reader := NewFileReader()
	return reader.ReadJSONFile(filepath)
}

// Parser genérico para estructuras
type Parser[T any] interface {
	Parse([]byte) (*T, error)
	ToBytes(*T) ([]byte, error)
}

// ConfigParser parser para configuración
type ConfigParser struct{}

func (p ConfigParser) Parse(b []byte) (*ConfigData, error) {
	var config Config
	if err := json.Unmarshal(b, &config); err != nil {
		return nil, fmt.Errorf("error al parsear configuración: %w", err)
	}
	return &config.Data, nil
}

func (p ConfigParser) ToBytes(data *ConfigData) ([]byte, error) {
	return json.Marshal(Config{Data: *data})
}

// TicketParser parser para tickets
type TicketParser struct{}

func (p TicketParser) Parse(b []byte) (*TicketData, error) {
	var ticket Ticket
	if err := json.Unmarshal(b, &ticket); err != nil {
		return nil, fmt.Errorf("error al parsear ticket: %w", err)
	}
	return &ticket.Data, nil
}

func (p TicketParser) ToBytes(data *TicketData) ([]byte, error) {
	return json.Marshal(Ticket{Data: *data})
}

// TemplateParser parser para templates
type TemplateParser struct{}

func (p TemplateParser) Parse(b []byte) (*TicketTemplateData, error) {
	var template TicketTemplate
	if err := json.Unmarshal(b, &template); err != nil {
		return nil, fmt.Errorf("error al parsear template: %w", err)
	}
	return &template.Data, nil
}

func (p TemplateParser) ToBytes(data *TicketTemplateData) ([]byte, error) {
	return json.Marshal(TicketTemplate{Data: *data})
}

// Funciones de compatibilidad
func BytesToConfig(b []byte) (*ConfigData, error) {
	parser := ConfigParser{}
	return parser.Parse(b)
}

func BytesToTicket(b []byte) (*TicketData, error) {
	parser := TicketParser{}
	return parser.Parse(b)
}

func BytesToTicketTemplate(b []byte) (*TicketTemplateData, error) {
	parser := TemplateParser{}
	return parser.Parse(b)
}

// BytesToNewTicket alias para BytesToTicket
func BytesToNewTicket(b []byte) (*TicketData, error) {
	return BytesToTicket(b)
}
