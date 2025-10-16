// Actualiza esto para que coincida con el endpoint WebSocket de tu servidor
const WS_URL = 'ws://localhost:8766/ws';
let socket = null;
let isConnected = false;

const statusEl = document.getElementById('status');
const textEl = document.getElementById('textContent');
const btnPrint = document.getElementById('btnPrint');
const btnPrintAndCut = document.getElementById('btnPrintAndCut');
const logEl = document.getElementById('log');

function log(message, type = 'info') {
    const entry = document.createElement('div');
    entry.className = `log-entry log-${type}`;
    entry.textContent = `[${new Date().toLocaleTimeString()}] ${message}`;
    logEl.appendChild(entry);
    logEl.scrollTop = logEl.scrollHeight;
}

function updateStatus(connected) {
    isConnected = connected;
    statusEl.textContent = connected ? 'Conectado' : 'Desconectado';
    statusEl.className = connected ? 'connected' : 'disconnected';
    btnPrint.disabled = !connected;
    btnPrintAndCut.disabled = !connected;
}

function connectWebSocket() {
    log('Conectando a ' + WS_URL + '...');
    socket = new WebSocket(WS_URL);

    socket.onopen = () => {
        log('Conexión establecida', 'success');
        updateStatus(true);
    };

    socket.onmessage = (event) => {
        try {
            const response = JSON.parse(event.data);

            if (response.tipo === 'ack') {
                statusEl.className = 'connected';
                if (response.success) {
                    log(`Impresión completada: ${response.message}`, 'success');
                } else {
                    log(`Error en impresión: ${response.message}`, 'error');
                }
            } else if (response.tipo === 'info') {
                log(`Servidor: ${response.message}`, 'info');
            }
        } catch (e) {
            log('Respuesta del servidor: ' + event.data, 'info');
        }
    };

    socket.onclose = () => {
        log('Conexión cerrada', 'error');
        updateStatus(false);
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

function sendPrintJob(cutAfter = false) {
    if (!isConnected || !socket) {
        log('No hay conexión activa', 'error');
        return;
    }

    const text = textEl.value.trim();
    if (!text) {
        log('No hay texto para imprimir', 'error');
        return;
    }

    const job = {
        tipo: 'print',
        texto: text,
        cortar: cutAfter
    };

    statusEl.className = 'printing';
    log(`Enviando trabajo de impresión${cutAfter ? ' con corte' : ''}...`, 'info');
    socket.send(JSON.stringify(job));
}

btnPrint.addEventListener('click', () => sendPrintJob(false));
btnPrintAndCut.addEventListener('click', () => sendPrintJob(true));

// Iniciar conexión
connectWebSocket();
