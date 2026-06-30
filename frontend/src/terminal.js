class Terminal {
    constructor(elementId) {
        this.element = document.getElementById(elementId);
        this.socket = null;
    }

    connect(sessionID) {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const host = window.location.hostname;
        const port = '8080';
        const url = `${protocol}//${host}:${port}/ws?sessionID=${sessionID}`;

        this.socket = new WebSocket(url);

        this.socket.onopen = () => {
            console.log('Connected to session');
        };

        this.socket.onmessage = (event) => {
            this.appendOutput(event.data);
        };

        this.socket.onclose = () => {
            console.log('Disconnected from session');
        };

        this.socket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
    }

    appendOutput(data) {
        this.element.textContent += data;
        this.element.scrollTop = this.element.scrollHeight;
    }

    send(data) {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(data);
        }
    }
}
