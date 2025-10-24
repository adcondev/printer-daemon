package models

// Concepto representa un producto o servicio en el ticket
type Concepto struct {
	Clave                 string     `json:"clave"`
	Descripcion           string     `json:"descripcion"`
	Cantidad              FloatFlex  `json:"cantidad"`
	Unidad                string     `json:"unidad"`
	PrecioVenta           FloatFlex  `json:"precio_venta"`
	Total                 FloatFlex  `json:"total"`
	ClaveProductoServicio string     `json:"clave_producto_servicio"`
	ClaveUnidadSAT        string     `json:"clave_unidad_sat"`
	VentaGranel           BoolFlex   `json:"venta_granel"`
	Impuestos             []Impuesto `json:"impuestos"`
	Series                []string   `json:"series,omitempty"`
}
