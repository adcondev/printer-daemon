package models

// SucursalData contiene los datos fiscales y de contacto de la sucursal
type SucursalData struct {
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
	SucursalEmails          string `json:"sucursal_emails"`
	SucursalTelefono        string `json:"sucursal_telefono"`
}
