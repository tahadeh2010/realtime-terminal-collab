document.addEventListener('DOMContentLoaded', function() {
    const sessionIdElement = document.getElementById('sessionId');
    const urlParams = new URLSearchParams(window.location.search);
    const sessionId = urlParams.get('sessionID');

    if (sessionId) {
        sessionIdElement.textContent = sessionId;
        const terminal = new Terminal('terminal');
        terminal.connect(sessionId);
    } else {
        sessionIdElement.textContent = 'No session ID provided';
    }
});
