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

function toggleMoreInfo(meetingType) {
    var moreText = document.getElementById(meetingType);
    var toggleText = document.getElementById(meetingType+"Btn");
  
    if (toggleText.innerHTML == "Show less") {
      toggleText.innerHTML = "Show more"; 
      moreText.style.display = "none";
    } else {
      toggleText.innerHTML = "Show less"; 
      moreText.style.display = "inline";
    }
  }