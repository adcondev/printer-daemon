package models

// Impuesto representa un impuesto aplicado a un concepto o de manera global
type Impuesto struct {
	Factor  string    `json:"factor"`
	Base    FloatFlex `json:"base"`
	Importe FloatFlex `json:"importe"`
	Codigo  string    `json:"impuestos"`
	Tasa    FloatFlex `json:"tasa"`
	Entidad string    `json:"entidad"`
	Tipo    string    `json:"tipo"`
}
