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
      // "Aborted" is fired when we cancel a running session ourselves (e.g. the
      // user starts a new one), so there is nothing to report to the user.
      if (error !== "Aborted") {
        showSessionError(error);
      }
    });
  window.scrollTo(0, 0);
}

// Turn a failed Yivi session into descriptive feedback instead of showing the
// raw status string (such as "CANCELLED") that the Yivi frontend rejects with.
function showSessionError(status) {
  var messages = window.yiviSessionMessages || {};
  var key = String(status == null ? "" : status).toLowerCase();

  var text;
  if (key.indexOf("cancel") !== -1) {
    text = messages.cancelled;
  } else if (key.indexOf("time") !== -1) {
    text = messages.timeout;
  } else if (key === "error" || key === "fail") {
    text = messages.error;
  } else {
    text = messages.fallback;
  }
  // Ultimate fallback if no translations were injected.
  if (!text) text = messages.fallback || String(status);

  var html = escapeHtml(text);
  if (messages.docsUrl && messages.docsLabel) {
    html +=
      ' <a href="' +
      encodeURI(messages.docsUrl) +
      '" target="_blank" rel="noopener noreferrer">' +
      escapeHtml(messages.docsLabel) +
      "</a>";
  }

  var errorboxtext = document.getElementById("errorboxtext");
  var errorbox = document.getElementById("errorbox");
  if (errorboxtext) errorboxtext.innerHTML = html;
  if (errorbox) {
    errorbox.classList.add("show");
    errorbox.style.opacity = "1";
  }
}

function escapeHtml(value) {
  return String(value == null ? "" : value)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#39;");
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
