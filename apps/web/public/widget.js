(function () {
  const scripts = document.getElementsByTagName("script");
  const currentScript = scripts[scripts.length - 1];

  const siteId = currentScript.getAttribute("data-site-id");
  const amount = currentScript.getAttribute("data-amount") || "5";

  if (!siteId) {
    console.error("SatGo: Missing data-site-id attribute.");
    return;
  }

  const container = document.createElement("div");
  container.style.width = "100%";
  container.style.maxWidth = "420px";
  container.style.margin = "0 auto";

  const iframe = document.createElement("iframe");
  const host = window.SatGoHost || "http://localhost:3000";
  iframe.src = `${host}/embed/${siteId}?amount=${amount}`;
  iframe.style.width = "100%";
  iframe.style.height = "620px";
  iframe.style.border = "none";
  iframe.style.borderRadius = "12px";
  iframe.style.boxShadow = "0 10px 25px rgba(0,0,0,0.1)";
  iframe.style.overflow = "hidden";
  iframe.setAttribute("sandbox", "allow-scripts allow-same-origin allow-popups allow-forms");

  container.appendChild(iframe);
  currentScript.parentNode.insertBefore(container, currentScript.nextSibling);
})();
