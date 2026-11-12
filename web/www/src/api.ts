export const getApiBaseUrl = (hostname = window.location.hostname) => {
  if (hostname === "climblive.com" || hostname === "labs.climblive.com") {
    return `https://${hostname.replace(/\.com$/, ".app")}/api`;
  }

  return "/api";
};
