// client.js - Cliente WebSocket para comunicación con el servidor de impresión

// Construir la URL del WebSocket dinámicamente basándose en el host actual
const WS_URL = `ws://${window.location.hostname}:8766/ws`;
let socket = null;
let isConnected = false;

// Contadores para elementos dinámicos
let conceptoCounter = 0;
let pagoCounter = 0;
let impuestoGlobalCounter = 0;

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
                if (isConnected) {
                    statusEl.textContent = 'Conectado';
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

// Función para manejar pestañas
function initializeTabs() {
    const tabButtons = document.querySelectorAll('.tab-button');
    const tabPanes = document.querySelectorAll('.tab-pane');

    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            const targetTab = button.getAttribute('data-tab');

            // Remover clase activa de todos los botones y paneles
            tabButtons.forEach(btn => btn.classList.remove('active'));
            tabPanes.forEach(pane => pane.classList.remove('active'));

            // Agregar clase activa al botón y panel seleccionados
            button.classList.add('active');
            document.getElementById(targetTab).classList.add('active');
        });
    });
}

// Función para agregar un concepto
function addConcepto() {
    conceptoCounter++;
    const container = document.getElementById('conceptosContainer');
    const conceptoDiv = document.createElement('div');
    conceptoDiv.className = 'concepto-item';
    conceptoDiv.id = `concepto-${conceptoCounter}`;

    conceptoDiv.innerHTML = `
        <div class="item-header">
            <h4>Concepto ${conceptoCounter}</h4>
            <button type="button" class="btn-small btn-delete" onclick="removeConcepto(${conceptoCounter})">❌ Eliminar</button>
        </div>
        <div class="form-row">
            <div class="form-group">
                <label>Clave:</label>
                <input type="text" id="concepto-clave-${conceptoCounter}" value="PROD${String(conceptoCounter).padStart(3, '0')}">
            </div>
            <div class="form-group">
                <label>Descripción:</label>
                <input type="text" id="concepto-descripcion-${conceptoCounter}" value="">
            </div>
            <div class="form-group">
                <label>Cantidad:</label>
                <input type="text" id="concepto-cantidad-${conceptoCounter}" value="1">
            </div>
            <div class="form-group">
                <label>Unidad:</label>
                <input type="text" id="concepto-unidad-${conceptoCounter}" value="PIEZA">
            </div>
        </div>
        <div class="form-row">
            <div class="form-group">
                <label>Precio Unitario:</label>
                <input type="text" id="concepto-precio-${conceptoCounter}" value="0">
            </div>
            <div class="form-group">
                <label>Total:</label>
                <input type="text" id="concepto-total-${conceptoCounter}" value="0">
            </div>
            <div class="form-group">
                <label>Clave Producto SAT:</label>
                <input type="text" id="concepto-clave-sat-${conceptoCounter}" value="01010101">
            </div>
            <div class="form-group">
                <label>Clave Unidad SAT:</label>
                <input type="text" id="concepto-unidad-sat-${conceptoCounter}" value="H87">
            </div>
        </div>
        <div class="form-row">
            <div class="form-group">
                <label>Series (separadas por coma):</label>
                <input type="text" id="concepto-series-${conceptoCounter}" value="">
            </div>
            <div class="form-group checkbox-group">
                <label>
                    <input type="checkbox" id="concepto-granel-${conceptoCounter}">
                    Venta a Granel
                </label>
            </div>
        </div>
        <details class="impuestos-concepto">
            <summary>Impuestos del Concepto</summary>
            <div id="concepto-impuestos-${conceptoCounter}">
                <button type="button" class="btn-small btn-add" onclick="addImpuestoConcepto(${conceptoCounter})">➕ Agregar Impuesto</button>
                <div id="concepto-impuestos-container-${conceptoCounter}"></div>
            </div>
        </details>
    `;

    container.appendChild(conceptoDiv);
}

// Función para eliminar un concepto
function removeConcepto(id) {
    const element = document.getElementById(`concepto-${id}`);
    if (element) {
        element.remove();
    }
}

// Función para agregar impuesto a un concepto
function addImpuestoConcepto(conceptoId) {
    const container = document.getElementById(`concepto-impuestos-container-${conceptoId}`);
    const impuestoId = `${conceptoId}-${Date.now()}`;

    const impuestoDiv = document.createElement('div');
    impuestoDiv.className = 'impuesto-item';
    impuestoDiv.id = `concepto-impuesto-${impuestoId}`;

    impuestoDiv.innerHTML = `
        <div class="form-row">
            <div class="form-group">
                <label>Factor:</label>
                <select id="impuesto-factor-${impuestoId}">
                    <option value="Tasa">Tasa</option>
                    <option value="Cuota">Cuota</option>
                    <option value="Exento">Exento</option>
                </select>
            </div>
            <div class="form-group">
                <label>Tipo:</label>
                <select id="impuesto-tipo-${impuestoId}">
                    <option value="T">Traslado</option>
                    <option value="R">Retención</option>
                </select>
            </div>
            <div class="form-group">
                <label>Impuesto:</label>
                <select id="impuesto-impuestos-${impuestoId}">
                    <option value="002">IVA</option>
                    <option value="003">IEPS</option>
                    <option value="001">ISR</option>
                </select>
            </div>
            <div class="form-group">
                <label>Tasa:</label>
                <input type="text" id="impuesto-tasa-${impuestoId}" value="0.16">
            </div>
        </div>
        <div class="form-row">
            <div class="form-group">
                <label>Base:</label>
                <input type="text" id="impuesto-base-${impuestoId}" value="0">
            </div>
            <div class="form-group">
                <label>Importe:</label>
                <input type="text" id="impuesto-importe-${impuestoId}" value="0">
            </div>
            <div class="form-group">
                <label>Entidad:</label>
                <input type="text" id="impuesto-entidad-${impuestoId}" value="Federal">
            </div>
            <button type="button" class="btn-small btn-delete" onclick="removeImpuestoConcepto('${impuestoId}')">❌</button>
        </div>
    `;

    container.appendChild(impuestoDiv);
}

// Función para eliminar impuesto de concepto
function removeImpuestoConcepto(id) {
    const element = document.getElementById(`concepto-impuesto-${id}`);
    if (element) {
        element.remove();
    }
}

// Función para agregar forma de pago
function addPago() {
    pagoCounter++;
    const container = document.getElementById('pagosContainer');
    const pagoDiv = document.createElement('div');
    pagoDiv.className = 'pago-item';
    pagoDiv.id = `pago-${pagoCounter}`;

    pagoDiv.innerHTML = `
        <div class="form-row">
            <div class="form-group">
                <label>Forma de Pago:</label>
                <select id="pago-forma-${pagoCounter}">
                    <option value="Efectivo">Efectivo</option>
                    <option value="Tarjeta de Crédito">Tarjeta de Crédito</option>
                    <option value="Tarjeta de Débito">Tarjeta de Débito</option>
                    <option value="Transferencia">Transferencia</option>
                    <option value="Cheque">Cheque</option>
                </select>
            </div>
            <div class="form-group">
                <label>Cantidad:</label>
                <input type="text" id="pago-cantidad-${pagoCounter}" value="0">
            </div>
            <div class="form-group">
                <label>Identificador:</label>
                <input type="text" id="pago-identificador-${pagoCounter}" value="">
            </div>
            <button type="button" class="btn-small btn-delete" onclick="removePago(${pagoCounter})">❌ Eliminar</button>
        </div>
    `;

    container.appendChild(pagoDiv);
}

// Función para eliminar forma de pago
function removePago(id) {
    const element = document.getElementById(`pago-${id}`);
    if (element) {
        element.remove();
    }
}

// Función para agregar impuesto global
function addImpuestoGlobal() {
    impuestoGlobalCounter++;
    const container = document.getElementById('impuestosGlobalContainer');
    const impuestoDiv = document.createElement('div');
    impuestoDiv.className = 'impuesto-global-item';
    impuestoDiv.id = `impuesto-global-${impuestoGlobalCounter}`;

    impuestoDiv.innerHTML = `
        <div class="form-row">
            <div class="form-group">
                <label>Factor:</label>
                <select id="impuesto-global-factor-${impuestoGlobalCounter}">
                    <option value="Tasa">Tasa</option>
                    <option value="Cuota">Cuota</option>
                    <option value="Exento">Exento</option>
                </select>
            </div>
            <div class="form-group">
                <label>Tipo:</label>
                <select id="impuesto-global-tipo-${impuestoGlobalCounter}">
                    <option value="T">Traslado</option>
                    <option value="R">Retención</option>
                </select>
            </div>
            <div class="form-group">
                <label>Impuesto:</label>
                <select id="impuesto-global-impuestos-${impuestoGlobalCounter}">
                    <option value="002">IVA</option>
                    <option value="003">IEPS</option>
                    <option value="001">ISR</option>
                </select>
            </div>
            <div class="form-group">
                <label>Tasa:</label>
                <input type="text" id="impuesto-global-tasa-${impuestoGlobalCounter}" value="0.16">
            </div>
        </div>
        <div class="form-row">
            <div class="form-group">
                <label>Base:</label>
                <input type="text" id="impuesto-global-base-${impuestoGlobalCounter}" value="0">
            </div>
            <div class="form-group">
                <label>Importe:</label>
                <input type="text" id="impuesto-global-importe-${impuestoGlobalCounter}" value="0">
            </div>
            <div class="form-group">
                <label>Entidad:</label>
                <input type="text" id="impuesto-global-entidad-${impuestoGlobalCounter}" value="Federal">
            </div>
            <button type="button" class="btn-small btn-delete" onclick="removeImpuestoGlobal(${impuestoGlobalCounter})">❌ Eliminar</button>
        </div>
    `;

    container.appendChild(impuestoDiv);
}

// Función para eliminar impuesto global
function removeImpuestoGlobal(id) {
    const element = document.getElementById(`impuesto-global-${id}`);
    if (element) {
        element.remove();
    }
}

// Función para limpiar todos los conceptos
function clearConceptos() {
    document.getElementById('conceptosContainer').innerHTML = '';
    conceptoCounter = 0;
}

// Función para limpiar todos los pagos
function clearPagos() {
    document.getElementById('pagosContainer').innerHTML = '';
    pagoCounter = 0;
}

// Función para limpiar todos los impuestos globales
function clearImpuestosGlobal() {
    document.getElementById('impuestosGlobalContainer').innerHTML = '';
    impuestoGlobalCounter = 0;
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

// Función para recolectar conceptos
function collectConceptos() {
    const conceptos = [];
    const conceptItems = document.querySelectorAll('.concepto-item');

    conceptItems.forEach(item => {
        const id = item.id.replace('concepto-', '');

        // Recolectar impuestos del concepto
        const impuestos = [];
        const impuestoItems = item.querySelectorAll('.impuesto-item');
        impuestoItems.forEach(impItem => {
            const impId = impItem.id.replace('concepto-impuesto-', '');
            impuestos.push({
                factor: document.getElementById(`impuesto-factor-${impId}`).value,
                base: document.getElementById(`impuesto-base-${impId}`).value,
                importe: document.getElementById(`impuesto-importe-${impId}`).value,
                impuestos: document.getElementById(`impuesto-impuestos-${impId}`).value,
                tasa: document.getElementById(`impuesto-tasa-${impId}`).value,
                entidad: document.getElementById(`impuesto-entidad-${impId}`).value,
                tipo: document.getElementById(`impuesto-tipo-${impId}`).value
            });
        });

        // Procesar series
        const seriesValue = document.getElementById(`concepto-series-${id}`).value;
        const series = seriesValue ? seriesValue.split(',').map(s => s.trim()).filter(s => s) : [];

        conceptos.push({
            clave: document.getElementById(`concepto-clave-${id}`).value,
            descripcion: document.getElementById(`concepto-descripcion-${id}`).value,
            cantidad: document.getElementById(`concepto-cantidad-${id}`).value,
            unidad: document.getElementById(`concepto-unidad-${id}`).value,
            precio_venta: document.getElementById(`concepto-precio-${id}`).value,
            total: document.getElementById(`concepto-total-${id}`).value,
            clave_producto_servicio: document.getElementById(`concepto-clave-sat-${id}`).value,
            clave_unidad_sat: document.getElementById(`concepto-unidad-sat-${id}`).value,
            venta_granel: document.getElementById(`concepto-granel-${id}`).checked ? "1" : "0",
            series: series,
            impuestos: impuestos
        });
    });

    return conceptos;
}

// Función para recolectar pagos
function collectPagos() {
    const pagos = [];
    const pagoItems = document.querySelectorAll('.pago-item');

    pagoItems.forEach(item => {
        const id = item.id.replace('pago-', '');
        pagos.push({
            forma_pago: document.getElementById(`pago-forma-${id}`).value,
            cantidad: document.getElementById(`pago-cantidad-${id}`).value,
            forma_pago_identificador: document.getElementById(`pago-identificador-${id}`).value
        });
    });

    return pagos;
}

// Función para recolectar impuestos globales
function collectImpuestosGlobal() {
    const impuestos = [];
    const impuestoItems = document.querySelectorAll('.impuesto-global-item');

    impuestoItems.forEach(item => {
        const id = item.id.replace('impuesto-global-', '');
        impuestos.push({
            factor: document.getElementById(`impuesto-global-factor-${id}`).value,
            base: document.getElementById(`impuesto-global-base-${id}`).value,
            importe: document.getElementById(`impuesto-global-importe-${id}`).value,
            impuestos: document.getElementById(`impuesto-global-impuestos-${id}`).value,
            tasa: document.getElementById(`impuesto-global-tasa-${id}`).value,
            entidad: document.getElementById(`impuesto-global-entidad-${id}`).value,
            tipo: document.getElementById(`impuesto-global-tipo-${id}`).value
        });
    });

    return impuestos;
}

// Función para recolectar datos del ticket
function collectTicketData() {
    // Recolectar comentarios (pueden ser null)
    const comentario = document.getElementById('ticketComentario').value.trim();
    const comentarioInterno = document.getElementById('ticketComentarioInterno').value.trim();
    const abonado = document.getElementById('ticketAbonado').value.trim();
    const descuentoNotaCredito = document.getElementById('ticketDescuentoNotaCredito').value.trim();

    return {
        data: {
            // Identificadores
            identificador: document.getElementById('ticketIdentificador').value,
            serie_identificador: document.getElementById('ticketSerieIdentificador').value,
            sucursal: document.getElementById('ticketSucursal').value,
            cliente_identificador: document.getElementById('ticketClienteIdentificador').value,
            vendedor_identificador: document.getElementById('ticketVendedorIdentificador').value,

            // Información general
            vendedor: document.getElementById('ticketVendedor').value,
            folio: document.getElementById('ticketFolio').value,
            serie: document.getElementById('ticketSerie').value,
            fecha_sistema: document.getElementById('ticketFecha').value,
            tipo_operacion: document.getElementById('ticketTipoOperacion').value,
            anulada: document.getElementById('ticketAnulada').checked ? "1" : "0",
            enviada: document.getElementById('ticketEnviada').checked ? "1" : "0",
            almacen_id: document.getElementById('ticketAlmacenId').value,
            tipo_conversion_factura: document.getElementById('ticketTipoConversionFactura').value,

            // Cliente
            cliente: document.getElementById('ticketClienteNombre').value,
            cliente_rfc: document.getElementById('ticketClienteRFC').value,
            cliente_cp: document.getElementById('ticketClienteCP').value,
            cliente_emails: document.getElementById('ticketClienteEmail').value,
            cliente_uso_cfdi: document.getElementById('ticketClienteUsoCFDI').value,
            cliente_regimen_fiscal: document.getElementById('ticketClienteRegimenFiscal').value,
            cliente_calle: document.getElementById('ticketClienteCalle').value,
            cliente_numero_exterior: document.getElementById('ticketClienteNumeroExterior').value,
            cliente_numero_interior: document.getElementById('ticketClienteNumeroInterior').value,
            cliente_colonia: document.getElementById('ticketClienteColonia').value,
            cliente_localidad: document.getElementById('ticketClienteLocalidad').value,
            cliente_delegacion: document.getElementById('ticketClienteDelegacion').value,
            cliente_estado: document.getElementById('ticketClienteEstado').value,
            cliente_pais: document.getElementById('ticketClientePais').value,

            // Sucursal
            sucursal_rfc: document.getElementById('ticketSucursalRFC').value,
            sucursal_nombre: document.getElementById('ticketSucursalNombre').value,
            sucursal_nombre_comercial: document.getElementById('ticketSucursalNombreComercial').value,
            sucursal_cp: document.getElementById('ticketSucursalCP').value,
            sucursal_regimen_clave: document.getElementById('ticketSucursalRegimenClave').value,
            sucursal_regimen: document.getElementById('ticketSucursalRegimen').value,
            sucursal_tienda: document.getElementById('ticketSucursalTienda').value,
            sucursal_email: document.getElementById('ticketSucursalEmail').value,
            sucursal_telefono: document.getElementById('ticketSucursalTelefono').value,
            sucursal_calle: document.getElementById('ticketSucursalCalle').value,
            sucursal_numero: document.getElementById('ticketSucursalNumero').value,
            sucursal_numero_int: document.getElementById('ticketSucursalNumeroInt').value,
            sucursal_colonia: document.getElementById('ticketSucursalColonia').value,
            sucursal_localidad: document.getElementById('ticketSucursalLocalidad').value,
            sucursal_municipio: document.getElementById('ticketSucursalMunicipio').value,
            sucursal_estado: document.getElementById('ticketSucursalEstado').value,
            sucursal_pais: document.getElementById('ticketSucursalPais').value,
            sucursal_leyenda_1: document.getElementById('ticketLeyenda1').value,
            sucursal_leyenda_2: document.getElementById('ticketLeyenda2').value,

            // Montos
            costo: document.getElementById('ticketCosto').value,
            costo_bruto: document.getElementById('ticketCostoBruto').value,
            total: document.getElementById('ticketTotal').value,
            saldo: document.getElementById('ticketSaldo').value,
            pagado: document.getElementById('ticketPagado').value,
            cambio: document.getElementById('ticketCambio').value,
            abonado: abonado || null,
            descuento: document.getElementById('ticketDescuento').value,
            descuento_motivo: document.getElementById('ticketDescuentoMotivo').value,
            descuento_nota_credito: descuentoNotaCredito || null,
            metodo_pago: document.getElementById('ticketMetodoPago').value,

            // Adicionales
            comentario: comentario || null,
            comentario_interno: comentarioInterno || null,
            autofactura_link: document.getElementById('ticketAutofacturaLink').value,
            autofactura_link_qr: document.getElementById('ticketAutofacturaLinkQr').value,

            // Receta (objeto vacío por defecto)
            receta: {},

            // Arrays
            conceptos: collectConceptos(),
            impuestos: collectImpuestosGlobal(),
            pago: collectPagos(),

            // Documentos de pago
            documentos_pago: [
                {
                    total: document.getElementById('ticketDocPagoTotal').value + ".000000",
                    tipo_cambio: document.getElementById('ticketDocPagoTipoCambio').value + ".000000",
                    saldo: document.getElementById('ticketDocPagoSaldo').value + ".000000",
                    nota: document.getElementById('ticketDocPagoNota').value,
                    sistema: document.getElementById('ticketDocPagoFecha').value,
                    anulado: document.getElementById('ticketDocPagoAnulado').checked ? "1" : "0",
                    cambio: document.getElementById('ticketCambio').value + ".000000",
                    fecha_pago: document.getElementById('ticketDocPagoFecha').value,
                    formas_pago: collectPagos()
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

// Event listeners para botones de agregar/limpiar
document.getElementById('btnAddConcepto').addEventListener('click', addConcepto);
document.getElementById('btnClearConceptos').addEventListener('click', clearConceptos);
document.getElementById('btnAddPago').addEventListener('click', addPago);
document.getElementById('btnClearPagos').addEventListener('click', clearPagos);
document.getElementById('btnAddImpuestoGlobal').addEventListener('click', addImpuestoGlobal);
document.getElementById('btnClearImpuestosGlobal').addEventListener('click', clearImpuestosGlobal);

// Inicializar con algunos elementos por defecto
window.addEventListener('load', () => {
    // Initialize tabs
    initializeTabs();

    // Connect WebSocket
    connectWebSocket();

    // Fill forms with JSON data - comment out these lines when you don't need test data
    fillConfigData();
    fillTemplateData();
    fillTicketData();
});

// Add these data filling functions after the collect functions

function fillConfigData() {
    document.getElementById('configPrinter').value = '80mm EC-PM-80250';
    document.getElementById('configDebugLog').checked = true;
}

function fillTemplateData() {
    // Exact values from template_ticket.json
    document.getElementById('templateTicketWidth').value = '80';
    document.getElementById('templateRazonSocialSize').value = '10';
    document.getElementById('templateDatosSize').value = '10';
    document.getElementById('templateLogoWidth').value = '120';

    // Checkboxes
    document.getElementById('templateVerLogotipo').checked = true;
    document.getElementById('templateVerNombre').checked = true;
    document.getElementById('templateVerNombreC').checked = true;
    document.getElementById('templateVerRFC').checked = true;
    document.getElementById('templateVerDom').checked = true;
    document.getElementById('templateVerRegimen').checked = true;
    document.getElementById('templateVerEmail').checked = true;
    document.getElementById('templateVerTelefono').checked = true;
    document.getElementById('templateVerNombreCliente').checked = true;
    document.getElementById('templateVerFolio').checked = true;
    document.getElementById('templateVerFecha').checked = true;
    document.getElementById('templateVerTienda').checked = true;
    document.getElementById('templateVerPrecioU').checked = true;
    document.getElementById('templateIncluyeImpuestos').checked = true;
    document.getElementById('templateVerCantProductos').checked = true;
    document.getElementById('templateVerImpuestos').checked = true;
    document.getElementById('templateVerImpuestosTotal').checked = true;
    document.getElementById('templateVerSeries').checked = true;
    document.getElementById('templateVerLeyenda').checked = true;

    // Text fields
    document.getElementById('templateCambiarCabecera').value = 'Ejemplo Cabecera';
    document.getElementById('templateCambiarReclamacion').value = 'PARA CUALQUIER RECLAMACION ES NECESARIO\r\nPRESENTAR SU TICKET DE COMPRAS';
    document.getElementById('templateCambiarPie').value = 'Ejemplo Pie';
}

function fillTicketData() {
    // General tab - exact values from data_ticket.json
    document.getElementById('ticketIdentificador').value = 'NTQ3';
    document.getElementById('ticketSerieIdentificador').value = 'MA==';
    document.getElementById('ticketFolio').value = '326';
    document.getElementById('ticketSerie').value = 'ABC1';
    document.getElementById('ticketFecha').value = '16/07/2025 12:18:18';
    document.getElementById('ticketVendedor').value = 'S0001';
    document.getElementById('ticketVendedorIdentificador').value = 'Ng==';
    document.getElementById('ticketSucursal').value = 'S0001';
    document.getElementById('ticketTipoOperacion').value = 'NOTA_VENTA';
    document.getElementById('ticketMetodoPago').value = 'PUE';
    document.getElementById('ticketAlmacenId').value = '5';
    document.getElementById('ticketTipoConversionFactura').value = '0';
    document.getElementById('ticketAnulada').checked = false;
    document.getElementById('ticketEnviada').checked = false;

    // Cliente tab
    document.getElementById('ticketClienteIdentificador').value = 'MTU=';
    document.getElementById('ticketClienteNombre').value = 'PUBLICO EN GENERAL';
    document.getElementById('ticketClienteRFC').value = 'XAXX010101000';
    document.getElementById('ticketClienteCP').value = '82000';
    document.getElementById('ticketClienteEmail').value = 'jrodriguez@red2000.mx';
    document.getElementById('ticketClienteUsoCFDI').value = 'G03';
    document.getElementById('ticketClienteRegimenFiscal').value = '616';
    document.getElementById('ticketClienteLocalidad').value = 'MAZATLAN';
    document.getElementById('ticketClienteDelegacion').value = 'MAZATLAN';
    document.getElementById('ticketClienteEstado').value = 'Sinaloa';
    document.getElementById('ticketClientePais').value = 'MEXICO';

    // Sucursal tab
    document.getElementById('ticketSucursalNombre').value = 'ESCUELA KEMPER URGATE';
    document.getElementById('ticketSucursalNombreComercial').value = 'LA RAZON';
    document.getElementById('ticketSucursalRFC').value = 'EKU9003173C9';
    document.getElementById('ticketSucursalCP').value = '82050';
    document.getElementById('ticketSucursalTienda').value = 'Almacen Principal';
    document.getElementById('ticketSucursalRegimenClave').value = '601';
    document.getElementById('ticketSucursalRegimen').value = 'REGIMEN ACTIVIDAD EMPRESARIAL Y PROFESIONAL PERSONA FISICA';
    document.getElementById('ticketSucursalTelefono').value = '982-66-09';
    document.getElementById('ticketSucursalEmail').value = 'lugoman58@gmail.com';
    document.getElementById('ticketSucursalCalle').value = 'Ejemplo 31';
    document.getElementById('ticketSucursalNumero').value = '123';
    document.getElementById('ticketSucursalNumeroInt').value = '111';
    document.getElementById('ticketSucursalColonia').value = 'Ejemplo 2';
    document.getElementById('ticketSucursalLocalidad').value = 'MAZATLAN';
    document.getElementById('ticketSucursalEstado').value = 'Sinaloa';
    document.getElementById('ticketSucursalPais').value = 'MEXICO';

    // Montos tab
    document.getElementById('ticketCosto').value = '500';
    document.getElementById('ticketCostoBruto').value = '234';
    document.getElementById('ticketTotal').value = '234';
    document.getElementById('ticketSaldo').value = '0';
    document.getElementById('ticketPagado').value = '234';
    document.getElementById('ticketCambio').value = '0';
    document.getElementById('ticketDescuento').value = '0';

    // Leyendas
    document.getElementById('ticketLeyenda1').value = '¡GRACIAS POR SU COMPRA!';
    document.getElementById('ticketLeyenda2').value = 'Desarrollado por Red2000';

    // Links
    document.getElementById('ticketAutofacturaLink').value = 'https://af.capacita.edu.mx/hola-mundo';
    document.getElementById('ticketAutofacturaLinkQr').value = 'https://af.capacita.edu.mx/hola-mundo?total=234&fecha=2025-07-16&folio=326';

    // Clear dynamic items first
    clearConceptos();
    clearPagos();

    // Add conceptos from JSON
    addConceptoFromData('PRO000029', 'Producto con Series 2', '3', 'Pieza', '78', '234', '27112309', 'H87', false, ['155548830', '155548834', '155548835']);

    // Add pago from JSON
    addPagoFromData('Efectivo', '234', 'eyJpdiI6IjNEK0s4NU0zaE9LdEdSVUtISGVMQ3c9PSIsInZhbHVlIjoiY01ISU41ckdkMHFnWTdLRVdlTjBwQT09IiwibWFjIjoiODE5ODNlMDEzMTMwN2FiNTgzOGY1OTU4OWEzNDczNmY0YjFiNDA0ZDMwZjliYzRiZTYwNmQwMDExODA1NjhiZiIsInRhZyI6IiJ9');

    // Documentos de pago
    document.getElementById('ticketDocPagoTotal').value = '234';
    document.getElementById('ticketDocPagoTipoCambio').value = '1';
    document.getElementById('ticketDocPagoSaldo').value = '0';
    document.getElementById('ticketDocPagoFecha').value = '2025-07-16 12:18:18';
}

// Helper functions to add dynamic items with specific data
function addConceptoFromData(clave, descripcion, cantidad, unidad, precio, total, claveSat, unidadSat, granel, series = []) {
    conceptoCounter++;
    const container = document.getElementById('conceptosContainer');
    const conceptoDiv = document.createElement('div');
    conceptoDiv.className = 'concepto-item';
    conceptoDiv.id = `concepto-${conceptoCounter}`;

    conceptoDiv.innerHTML = `
        <div class="item-header">
            <h4>Concepto ${conceptoCounter}</h4>
            <button type="button" class="btn-small btn-delete" onclick="removeConcepto(${conceptoCounter})">❌ Eliminar</button>
        </div>
        <div class="form-row">
            <div class="form-group">
                <label>Clave:</label>
                <input type="text" id="concepto-clave-${conceptoCounter}" value="${clave}">
            </div>
            <div class="form-group">
                <label>Descripción:</label>
                <input type="text" id="concepto-descripcion-${conceptoCounter}" value="${descripcion}">
            </div>
            <div class="form-group">
                <label>Cantidad:</label>
                <input type="text" id="concepto-cantidad-${conceptoCounter}" value="${cantidad}">
            </div>
            <div class="form-group">
                <label>Unidad:</label>
                <input type="text" id="concepto-unidad-${conceptoCounter}" value="${unidad}">
            </div>
        </div>
        <div class="form-row">
            <div class="form-group">
                <label>Precio Unitario:</label>
                <input type="text" id="concepto-precio-${conceptoCounter}" value="${precio}">
            </div>
            <div class="form-group">
                <label>Total:</label>
                <input type="text" id="concepto-total-${conceptoCounter}" value="${total}">
            </div>
            <div class="form-group">
                <label>Clave Producto SAT:</label>
                <input type="text" id="concepto-clave-sat-${conceptoCounter}" value="${claveSat}">
            </div>
            <div class="form-group">
                <label>Clave Unidad SAT:</label>
                <input type="text" id="concepto-unidad-sat-${conceptoCounter}" value="${unidadSat}">
            </div>
        </div>
        <div class="form-row">
            <div class="form-group">
                <label>Series (separadas por coma):</label>
                <input type="text" id="concepto-series-${conceptoCounter}" value="${series.join(', ')}">
            </div>
            <div class="form-group checkbox-group">
                <label>
                    <input type="checkbox" id="concepto-granel-${conceptoCounter}" ${granel ? 'checked' : ''}>
                    Venta a Granel
                </label>
            </div>
        </div>
        <details class="impuestos-concepto">
            <summary>Impuestos del Concepto</summary>
            <div id="concepto-impuestos-${conceptoCounter}">
                <button type="button" class="btn-small btn-add" onclick="addImpuestoConcepto(${conceptoCounter})">➕ Agregar Impuesto</button>
                <div id="concepto-impuestos-container-${conceptoCounter}"></div>
            </div>
        </details>
    `;

    container.appendChild(conceptoDiv);
}

function addPagoFromData(formaPago, cantidad, identificador) {
    pagoCounter++;
    const container = document.getElementById('pagosContainer');
    const pagoDiv = document.createElement('div');
    pagoDiv.className = 'pago-item';
    pagoDiv.id = `pago-${pagoCounter}`;

    pagoDiv.innerHTML = `
        <div class="form-row">
            <div class="form-group">
                <label>Forma de Pago:</label>
                <select id="pago-forma-${pagoCounter}">
                    <option value="Efectivo" ${formaPago === 'Efectivo' ? 'selected' : ''}>Efectivo</option>
                    <option value="Tarjeta de Crédito" ${formaPago === 'Tarjeta de Crédito' ? 'selected' : ''}>Tarjeta de Crédito</option>
                    <option value="Tarjeta de Débito" ${formaPago === 'Tarjeta de Débito' ? 'selected' : ''}>Tarjeta de Débito</option>
                    <option value="Transferencia" ${formaPago === 'Transferencia' ? 'selected' : ''}>Transferencia</option>
                    <option value="Cheque" ${formaPago === 'Cheque' ? 'selected' : ''}>Cheque</option>
                </select>
            </div>
            <div class="form-group">
                <label>Cantidad:</label>
                <input type="text" id="pago-cantidad-${pagoCounter}" value="${cantidad}">
            </div>
            <div class="form-group">
                <label>Identificador:</label>
                <input type="text" id="pago-identificador-${pagoCounter}" value="${identificador}">
            </div>
            <button type="button" class="btn-small btn-delete" onclick="removePago(${pagoCounter})">❌ Eliminar</button>
        </div>
    `;

    container.appendChild(pagoDiv);
}