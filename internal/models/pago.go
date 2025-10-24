package models

// DocumentoPago representa un documento de pago asociado al ticket
type DocumentoPago struct {
	Total      FloatFlex `json:"total"`
	TipoCambio FloatFlex `json:"tipo_cambio"`
	Saldo      FloatFlex `json:"saldo"`
	Nota       string    `json:"nota"`
	Sistema    string    `json:"sistema"`
	Anulado    BoolFlex  `json:"anulado"`
	Cambio     FloatFlex `json:"cambio"`
	FechaPago  string    `json:"fecha_pago"`
	FormasPago []Pago    `json:"formas_pago"`
}

// Pago representa una forma de pago utilizada en el ticket
type Pago struct {
	FormaPago     string    `json:"forma_pago"`
	Cantidad      FloatFlex `json:"cantidad"`
	Identificador string    `json:"forma_pago_identificador"`
}
