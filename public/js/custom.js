function startInvitation(authUrl, redirectUrl) {
    fetch(authUrl, { credentials: 'include' })
        .then((response) => {
            return response.json();
        })
        .then((qr) => {
            return irma.handleSession(qr, { language: LANGUAGE });
        })
        .then(() => {
            window.location.replace(redirectUrl);
        })
        .catch((error) => {
            document.getElementById('errorboxtext').innerText = error;
            document.getElementById('errorbox').style.opacity = "1";
        });
}
