package models

// Config es la estructura raíz para la configuración local
type Config struct {
	Data ConfigData `json:"data"`
}

// ConfigData contiene la configuración local de la aplicación
type ConfigData struct {
	// Configuración general
	Printer  string `json:"printer"`   // Nombre de la impresora a utilizar
	DebugLog bool   `json:"debug_log"` // Habilitar logs de depuración
}
