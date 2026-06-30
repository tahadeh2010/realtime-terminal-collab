document.addEventListener('DOMContentLoaded', function() {
    const joinForm = document.getElementById('joinForm');
    const sessionIdInput = document.getElementById('sessionIdInput');

    joinForm.addEventListener('submit', function(e) {
        e.preventDefault();
        const sessionId = sessionIdInput.value.trim();
        if (sessionId) {
            window.location.href = `session.html?sessionID=${encodeURIComponent(sessionId)}`;
        }
    });
});
