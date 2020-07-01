function startInvitation(authUrl, redirectUrl) {
    fetch(authUrl, { credentials: 'include' })
        .then((response) => {
            return response.json();
        })
        .then((qr) => {
            return irma.handleSession(qr, { language: 'nl' });
        })
        .then(() => {
            window.location.replace(redirectUrl);
        })
        .catch((error) => {
            document.getElementById('errorboxtext').innerText = error;
            document.getElementById('errorbox').style.opacity = "1";
        });
}

function toggleMoreOrLessInfo() {
    var dots = document.getElementById("dots");
    var moreText = document.getElementById("more");
    var btnText = document.getElementById("myBtn");
  
    if (dots.style.display === "none") {
      dots.style.display = "inline";
      btnText.innerHTML = "Read more"; 
      moreText.style.display = "none";
    } else {
      dots.style.display = "none";
      btnText.innerHTML = "Read less"; 
      moreText.style.display = "inline";
    }
  }