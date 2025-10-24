package models

// ClienteData contiene los datos fiscales y de contacto del cliente
type ClienteData struct {
	ClienteNombre         string `json:"cliente"`
	ClienteRFC            string `json:"cliente_rfc"`
	ClienteCP             string `json:"cliente_cp"`
	ClienteUsoCFDI        string `json:"cliente_uso_cfdi"`
	ClienteRegimenFiscal  string `json:"cliente_regimen_fiscal"`
	ClienteCalle          string `json:"cliente_calle"`
	ClienteNumeroExterior string `json:"cliente_numero_exterior"`
	ClienteNumeroInterior string `json:"cliente_numero_interior"`
	ClienteColonia        string `json:"cliente_colonia"`
	ClienteLocalidad      string `json:"cliente_localidad"`
	ClienteMunicipio      string `json:"cliente_delegacion"`
	ClienteEstado         string `json:"cliente_estado"`
	ClientePais           string `json:"cliente_pais"`
	ClienteEmail          string `json:"cliente_emails"`
}
