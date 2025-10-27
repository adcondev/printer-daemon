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
	DescuentoNotaCred FloatFlex `json:"descuento_nota_credito"`
	Total             FloatFlex `json:"total"`
	Saldo             FloatFlex `json:"saldo"`
	Pagado            FloatFlex `json:"pagado"`
	Cambio            FloatFlex `json:"cambio"`

	// Datos del cliente
	Cliente               string `json:"cliente"`
	ClienteRFC            string `json:"cliente_rfc"`
	ClienteCP             string `json:"cliente_cp"`
	ClienteUsoCFDI        string `json:"cliente_uso_cfdi"`
	ClienteRegimenFiscal  string `json:"cliente_regimen_fiscal"`
	ClienteCalle          string `json:"cliente_calle"`
	ClienteNumeroExterior string `json:"cliente_numero_exterior"`
	ClienteNumeroInterior string `json:"cliente_numero_interior"`
	ClienteColonia        string `json:"cliente_colonia"`
	ClienteLocalidad      string `json:"cliente_localidad"`
	ClienteDelegacion     string `json:"cliente_delegacion"`
	ClienteEstado         string `json:"cliente_estado"`
	ClientePais           string `json:"cliente_pais"`
	ClienteEmail          string `json:"cliente_emails"`

	// Datos de la sucursal
	SucursalRFC             string `json:"sucursal_rfc"`
	SucursalNombre          string `json:"sucursal_nombre"`
	SucursalNombreComercial string `json:"sucursal_nombre_comercial"`
	SucursalTienda          string `json:"sucursal_tienda"`
	SucursalRegimenClave    string `json:"sucursal_regimen_clave"`
	SucursalRegimen         string `json:"sucursal_regimen"`
	SucursalCalle           string `json:"sucursal_calle"`
	SucursalNumero          string `json:"sucursal_numero"`
	SucursalNumeroInt       string `json:"sucursal_numero_int"`
	SucursalColonia         string `json:"sucursal_colonia"`
	SucursalLocalidad       string `json:"sucursal_localidad"`
	SucursalMunicipio       string `json:"sucursal_municipio"`
	SucursalEstado          string `json:"sucursal_estado"`
	SucursalCP              string `json:"sucursal_cp"`
	SucursalPais            string `json:"sucursal_pais"`
	SucursalTelefono        string `json:"sucursal_telefono"`

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
	Comentario            *string `json:"comentario,omitempty"`
	ComentarioInterno     *string `json:"comentario_interno,omitempty"`
	AlmacenID             string  `json:"almacen_id"`
	TipoConversionFactura string  `json:"tipo_conversion_factura"`

	// Montos económicos (strings en el JSON)
	Costo           FloatFlex `json:"costo"`
	CostoBruto      FloatFlex `json:"costo_bruto"`
	DescuentoMotivo string    `json:"descuento_motivo"`
	MetodoPago      string    `json:"metodo_pago"`
	Abonado         FloatFlex `json:"abonado"`

	// Datos de recetas
	Receta map[string]any `json:"receta"`

	// Impuestos globales
	Impuestos []Impuesto `json:"impuestos"`
}
