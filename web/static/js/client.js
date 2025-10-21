// client.js - Cliente WebSocket para comunicación con el servidor de impresión

// Actualiza esto para que coincida con el endpoint WebSocket de tu servidor
const WS_URL = "ws://192.168.8.76:8766/ws";
let socket = null;
let isConnected = false;

// Elementos del DOM
const statusEl = document.getElementById('status');
const logEl = document.getElementById('log');

// Botones
const btnSendConfig = document.getElementById('btnSendConfig');
const btnSendTemplate = document.getElementById('btnSendTemplate');
const btnSendTicket = document.getElementById('btnSendTicket');

// Función para agregar mensajes al log
function log(message, type = 'info') {
    const entry = document.createElement('div');
    entry.className = `log-entry log-${type}`;
    const timestamp = new Date().toLocaleTimeString('es-MX');
    entry.textContent = `[${timestamp}] ${message}`;
    logEl.appendChild(entry);
    logEl.scrollTop = logEl.scrollHeight;
}

// Función para actualizar el estado de la conexión
function updateStatus(connected) {
    isConnected = connected;
    statusEl.textContent = connected ? 'Conectado' : 'Desconectado';
    statusEl.className = connected ? 'connected' : 'disconnected';

    // Habilitar/deshabilitar botones según el estado
    btnSendConfig.disabled = !connected;
    btnSendTemplate.disabled = !connected;
    btnSendTicket.disabled = !connected;
}

// Función para conectar al WebSocket
function connectWebSocket() {
    log('Conectando a ' + WS_URL + '...');
    socket = new WebSocket(WS_URL);

    socket.onopen = () => {
        log('Conexión establecida con el servidor', 'success');
        updateStatus(true);
    };

    socket.onmessage = (event) => {
        try {
            const response = JSON.parse(event.data);

            if (response.tipo === 'info') {
                log(response.message, 'info');
            } else if (response.tipo === 'ack') {
                // Restaurar estado normal después de procesar
                if (statusEl.className === 'printing') {
                    statusEl.className = 'connected';
                }

                if (response.success) {
                    log(response.message, 'success');
                } else {
                    log(response.message, 'error');
                }
            }
        } catch (e) {
            log('Respuesta del servidor: ' + event.data, 'info');
        }
    };

    socket.onclose = () => {
        log('Conexión cerrada', 'error');
        updateStatus(false);

        // Intentar reconectar después de 3 segundos
        setTimeout(() => {
            log('Reintentando conexión...', 'info');
            connectWebSocket();
        }, 3000);
    };

    socket.onerror = (error) => {
        log('Error de WebSocket', 'error');
        console.error('WebSocket error:', error);
    };
}

// Función para recolectar datos de configuración
function collectConfigData() {
    return {
        data: {
            printer: document.getElementById('configPrinter').value,
            debug_log: document.getElementById('configDebugLog').checked
        }
    };
}

// Función para recolectar datos de plantilla
function collectTemplateData() {
    return {
        data: {
            ticket_width: document.getElementById('templateTicketWidth').value,
            razon_social_size: document.getElementById('templateRazonSocialSize').value,
            datos_size: document.getElementById('templateDatosSize').value,
            logo_width: document.getElementById('templateLogoWidth').value,
            ver_logotipo: document.getElementById('templateVerLogotipo').checked ? "1" : "0",
            ver_nombre: document.getElementById('templateVerNombre').checked ? "1" : "0",
            ver_nombre_c: document.getElementById('templateVerNombreC').checked ? "1" : "0",
            ver_rfc: document.getElementById('templateVerRFC').checked ? "1" : "0",
            ver_dom: document.getElementById('templateVerDom').checked ? "1" : "0",
            ver_leyenda: document.getElementById('templateVerLeyenda').checked ? "1" : "0",
            ver_regimen: document.getElementById('templateVerRegimen').checked ? "1" : "0",
            ver_email: document.getElementById('templateVerEmail').checked ? "1" : "0",
            ver_telefono: document.getElementById('templateVerTelefono').checked ? "1" : "0",
            ver_nombre_cliente: document.getElementById('templateVerNombreCliente').checked ? "1" : "0",
            ver_folio: document.getElementById('templateVerFolio').checked ? "1" : "0",
            ver_fecha: document.getElementById('templateVerFecha').checked ? "1" : "0",
            ver_tienda: document.getElementById('templateVerTienda').checked ? "1" : "0",
            ver_precio_u: document.getElementById('templateVerPrecioU').checked ? "1" : "0",
            incluye_impuestos: document.getElementById('templateIncluyeImpuestos').checked ? "1" : "0",
            ver_cant_productos: document.getElementById('templateVerCantProductos').checked ? "1" : "0",
            ver_impuestos: document.getElementById('templateVerImpuestos').checked ? "1" : "0",
            ver_impuestos_total: document.getElementById('templateVerImpuestosTotal').checked ? "1" : "0",
            ver_series: document.getElementById('templateVerSeries').checked ? "1" : "0",
            cambiar_cabecera: document.getElementById('templateCambiarCabecera').value,
            cambiar_reclamacion: document.getElementById('templateCambiarReclamacion').value,
            cambiar_pie: document.getElementById('templateCambiarPie').value
        }
    };
}

// Función para recolectar datos del ticket
function collectTicketData() {
    return {
        data: {
            identificador: "NTQ3",
            serie_identificador: "MA==",
            sucursal: "S0001",
            cliente_identificador: "MTU=",
            vendedor_identificador: "Ng==",
            vendedor: document.getElementById('ticketVendedor').value,
            cliente: document.getElementById('ticketClienteNombre').value,
            cliente_rfc: document.getElementById('ticketClienteRFC').value,
            cliente_cp: document.getElementById('ticketClienteCP').value,
            cliente_uso_cfdi: "G03",
            cliente_regimen_fiscal: "616",
            cliente_calle: "",
            cliente_numero_exterior: "",
            cliente_numero_interior: "",
            cliente_estado: "Sinaloa",
            cliente_colonia: "",
            cliente_localidad: "MAZATLÁN",
            cliente_delegacion: "MAZATLÁN",
            cliente_pais: "MÉXICO",
            sucursal_rfc: document.getElementById('ticketSucursalRFC').value,
            sucursal_nombre: document.getElementById('ticketSucursalNombre').value,
            sucursal_cp: "82050",
            sucursal_regimen_clave: "601",
            sucursal_regimen: "RÉGIMEN ACTIVIDAD EMPRESARIAL Y PROFESIONAL PERSONA FÍSICA",
            sucursal_nombre_comercial: document.getElementById('ticketSucursalNombreComercial').value,
            sucursal_tienda: document.getElementById('ticketSucursalTienda').value,
            sucursal_email: document.getElementById('ticketSucursalEmail').value,
            sucursal_telefono: document.getElementById('ticketSucursalTelefono').value,
            sucursal_calle: "Ejemplo 31",
            sucursal_numero: "123",
            sucursal_numero_int: "111",
            sucursal_colonia: "Ejemplo 2",
            sucursal_estado: "Sinaloa",
            sucursal_localidad: "MAZATLÁN",
            sucursal_municipio: "",
            sucursal_pais: "MÉXICO",
            sucursal_leyenda_1: document.getElementById('ticketLeyenda1').value,
            sucursal_leyenda_2: document.getElementById('ticketLeyenda2').value,
            comentario: null,
            comentario_interno: null,
            costo: "500",
            costo_bruto: document.getElementById('ticketTotal').value,
            descuento_motivo: "",
            descuento: document.getElementById('ticketDescuento').value,
            metodo_pago: "PUE",
            total: document.getElementById('ticketTotal').value,
            saldo: "0",
            abonado: null,
            pagado: document.getElementById('ticketPagado').value,
            cambio: document.getElementById('ticketCambio').value,
            descuento_nota_credito: null,
            tipo_operacion: "NOTA_VENTA",
            anulada: "0",
            fecha_sistema: document.getElementById('ticketFecha').value,
            folio: document.getElementById('ticketFolio').value,
            serie: document.getElementById('ticketSerie').value,
            enviada: "0",
            cliente_emails: document.getElementById('ticketClienteEmail').value,
            almacen_id: "5",
            tipo_conversion_factura: "0",
            autofactura_link: "https://af.capacita.edu.mx/hola-mundo",
            autofactura_link_qr: "https://af.capacita.edu.mx/hola-mundo?total=" +
                document.getElementById('ticketTotal').value + "&fecha=2025-07-16&folio=" +
                document.getElementById('ticketFolio').value,
            receta: {},
            conceptos: [
                {
                    clave: "PRO000029",
                    descripcion: "Producto con Series 2",
                    precio_venta: "78",
                    cantidad: "3",
                    total: "234",
                    unidad: "Pieza",
                    clave_producto_servicio: "27112309",
                    clave_unidad_sat: "H87",
                    venta_granel: "0",
                    series: ["155548830", "155548834", "155548835"],
                    impuestos: []
                },
                {
                    clave: "0001",
                    descripcion: "MANTENIMIENTO OTROS CORRECTIVOS",
                    precio_venta: "37041.8",
                    cantidad: "1",
                    total: "37041.8",
                    unidad: "PIEZA",
                    clave_producto_servicio: "52101502",
                    clave_unidad_sat: "MTK",
                    venta_granel: "0",
                    impuestos: []
                }
            ],
            impuestos: [],
            documentos_pago: [
                {
                    total: document.getElementById('ticketTotal').value + ".000000",
                    tipo_cambio: "1.000000",
                    saldo: "0.000000",
                    nota: "",
                    sistema: document.getElementById('ticketFecha').value,
                    anulado: "0",
                    cambio: document.getElementById('ticketCambio').value + ".000000",
                    fecha_pago: document.getElementById('ticketFecha').value,
                    formas_pago: [
                        {
                            forma_pago: "Efectivo",
                            cantidad: document.getElementById('ticketPagado').value,
                            forma_pago_identificador: "MQ=="
                        }
                    ]
                }
            ],
            pago: [
                {
                    forma_pago: "Efectivo",
                    cantidad: document.getElementById('ticketPagado').value,
                    forma_pago_identificador: "eyJpdiI6IjNEK0s4NU0zaE9LdEdSVUtISGVMQ3c9PSIsInZhbHVlIjoiY01ISU41ckdkMHFnWTdLRVdlTjBwQT09IiwibWFjIjoiODE5ODNlMDEzMTMwN2FiNTgzOGY1OTU4OWEzNDczNmY0YjFiNDA0ZDMwZjliYzRiZTYwNmQwMDExODA1NjhiZiIsInRhZyI6IiJ9"
                }
            ]
        }
    };
}

// Función para enviar datos al servidor
function sendData(type, data) {
    if (!isConnected || !socket) {
        log('No hay conexión con el servidor', 'error');
        return;
    }

    const message = {
        tipo: type,
        datos: data
    };

    // Cambiar estado a "procesando"
    statusEl.className = 'printing';
    statusEl.textContent = 'Procesando...';

    socket.send(JSON.stringify(message));
    log(`Enviando ${type} al servidor...`, 'info');
}

// Event listeners para los botones
btnSendConfig.addEventListener('click', () => {
    const configData = collectConfigData();
    sendData('config', configData);
});

btnSendTemplate.addEventListener('click', () => {
    const templateData = collectTemplateData();
    sendData('template', templateData);
});

btnSendTicket.addEventListener('click', () => {
    const ticketData = collectTicketData();
    sendData('ticket', ticketData);
});

// Iniciar conexión cuando se carga la página
window.addEventListener('load', () => {
    connectWebSocket();
});