package models

// Ticket representa la estructura raíz de un ticket de venta
type Ticket struct {
	Data TicketData `json:"data"`
}

// TicketData contiene todos los datos de un ticket de venta
type TicketData struct {
	// Metadatos del ticket
	Identificador string   `json:"identificador"`
	Vendedor      string   `json:"vendedor"`
	Folio         string   `json:"folio"`
	Serie         string   `json:"serie"`
	FechaSistema  string   `json:"fecha_sistema"`
	TipoOperacion string   `json:"tipo_operacion"`
	Anulada       BoolFlex `json:"anulada"`

	// Montos del ticket
	Descuento         FloatFlex `json:"descuento"`
	DescuentoNotaCred FloatFlex `json:"descuento_nota_credito,omitempty"`
	Total             FloatFlex `json:"total"`
	Saldo             FloatFlex `json:"saldo"`
	Pagado            FloatFlex `json:"pagado"`
	Cambio            FloatFlex `json:"cambio"`

	// Composición de entidades
	ClienteData  ClienteData  `json:",inline"`
	SucursalData SucursalData `json:",inline"`

	// Enlaces y códigos QR
	AutofacturaLink   string `json:"autofactura_link"`
	AutofacturaLinkQr string `json:"autofactura_link_qr"`

	// Conceptos y pagos
	Conceptos      []Concepto      `json:"conceptos"`
	DocumentosPago []DocumentoPago `json:"documentos_pago"`
	Pagos          []Pago          `json:"pago"`

	// Identificadores codificados
	SerieIdentificador    string `json:"serie_identificador"`
	SucursalIdentificador string `json:"sucursal"`
	ClienteIdentificador  string `json:"cliente_identificador"`
	VendedorIdentificador string `json:"vendedor_identificador"`

	// Estados
	Enviada BoolFlex `json:"enviada"`

	// Datos de la sucursal
	SucursalEmail    string `json:"sucursal_email"`
	SucursalLeyenda1 string `json:"sucursal_leyenda_1"` // Leyenda 1
	SucursalLeyenda2 string `json:"sucursal_leyenda_2"` // Leyenda 2

	// Metadatos adicionales
	Comentario            any    `json:"comentario"`
	ComentarioInterno     any    `json:"comentario_interno"`
	AlmacenID             string `json:"almacen_id"`
	TipoConversionFactura string `json:"tipo_conversion_factura"`

	// Montos económicos (strings en el JSON)
	Costo           string `json:"costo"`
	CostoBruto      string `json:"costo_bruto"`
	DescuentoMotivo string `json:"descuento_motivo"`
	MetodoPago      string `json:"metodo_pago"`
	Abonado         any    `json:"abonado"`

	// Datos de recetas
	Receta map[string]any `json:"receta"`

	// Impuestos globales
	Impuestos []Impuesto `json:"impuestos"`
}
