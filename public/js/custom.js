function startInvitation() {
    fetch('./irma_auth/start', { credentials: 'include' })
        .then((response) => {
            return response.json();
        })
        .then((qr) => {
            return irma.handleSession(qr, { language: 'nl' });
        })
        .then(() => {
            window.location.replace('./irma_session/create');
        })
        .catch((error) => {
            document.getElementById('errorboxtext').innerText = error;
            document.getElementById('errorbox').style.opacity = "1";
        });
}
