var irmaFrontend;

function startInvitation(authUrl, redirectUrl) {
  document.getElementById("card-image").style.display = "none";
  document.getElementsByClassName("irma-web-form")[0].removeAttribute("style");
  if (irmaFrontend) irmaFrontend.abort();

  irmaFrontend = irma.newWeb({
    debugging: false,
    language: LANGUAGE,
    element: ".irma-web-form",

    session: {
      url: authUrl,

      start: {
        url: o => o.url,
        method: "GET",
        headers: {},
        credentials: "include"
      },

      mapping: {
        sessionPtr: r => r
      },

      result: false
    }
  });

  irmaFrontend
    .start()
    .then(() => {
      window.location.replace(redirectUrl);
    })
    .catch(error => {
      if (error !== "Aborted") {
        document.getElementById("errorboxtext").innerText = error;
        document.getElementById("errorbox").style.opacity = "1";
      }
    });
  window.scrollTo(0, 0);
}

function copy(text) {
  var input = document.createElement("textarea");
  input.innerHTML = text;
  document.body.appendChild(input);
  input.select();
  var result = document.execCommand("copy");
  document.body.removeChild(input);
  return result;
}
