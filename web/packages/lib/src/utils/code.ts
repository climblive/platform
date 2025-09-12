export const extractCodeFromPath = () => {
  const match = window.location.pathname.match(/\/([A-Z0-9]{8})/i);
  return match ? match[1] : null;
};
