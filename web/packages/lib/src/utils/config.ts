import configData from "../config.json";

export const getApiUrl = () => {
  return configData.API_URL || "/api";
};
