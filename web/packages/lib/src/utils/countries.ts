import countriesData from "../../../../../countries.json";

export interface Country {
  code: string; // ISO 3166-1 alpha-2 code
  name: string;
}

export const countries: Country[] = countriesData as Country[];

const REGIONAL_INDICATOR_OFFSET = 127397;

const countryMap = new Map(countries.map((c) => [c.code, c.name]));

export function getFlag(countryCode: string | undefined): string {
  if (!countryCode || countryCode.length !== 2) return "";
  const codePoints = countryCode
    .toUpperCase()
    .split("")
    .map((char) => REGIONAL_INDICATOR_OFFSET + char.charCodeAt(0));
  return String.fromCodePoint(...codePoints);
}

export function getCountryName(countryCode: string | undefined): string {
  if (!countryCode) return "";
  return countryMap.get(countryCode) || "";
}

/**
 * Attempts to guess the user's country based on their browser timezone.
 * Used to pre-select a default country when creating new contests.
 * Falls back to 'AQ' (Antarctica) if timezone is not recognized.
 */
export function guessCountryFromTimezone(): string {
  const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;

  const timezoneToCountry: Record<string, string> = {
    "Europe/Stockholm": "SE",
    "Europe/Oslo": "NO",
    "Europe/Copenhagen": "DK",
    "Europe/Helsinki": "FI",
    "Europe/Reykjavik": "IS",
    "Europe/London": "GB",
    "Europe/Dublin": "IE",
    "Europe/Paris": "FR",
    "Europe/Berlin": "DE",
    "Europe/Amsterdam": "NL",
    "Europe/Brussels": "BE",
    "Europe/Luxembourg": "LU",
    "Europe/Zurich": "CH",
    "Europe/Vienna": "AT",
    "Europe/Rome": "IT",
    "Europe/Madrid": "ES",
    "Europe/Lisbon": "PT",
    "Europe/Athens": "GR",
    "Europe/Warsaw": "PL",
    "Europe/Prague": "CZ",
    "Europe/Budapest": "HU",
    "Europe/Bucharest": "RO",
    "Europe/Sofia": "BG",
    "America/New_York": "US",
    "America/Chicago": "US",
    "America/Denver": "US",
    "America/Los_Angeles": "US",
    "America/Toronto": "CA",
    "America/Vancouver": "CA",
    "America/Mexico_City": "MX",
    "America/Sao_Paulo": "BR",
    "America/Buenos_Aires": "AR",
    "Asia/Tokyo": "JP",
    "Asia/Seoul": "KR",
    "Asia/Shanghai": "CN",
    "Asia/Hong_Kong": "HK",
    "Asia/Singapore": "SG",
    "Asia/Bangkok": "TH",
    "Asia/Dubai": "AE",
    "Asia/Kolkata": "IN",
    "Australia/Sydney": "AU",
    "Australia/Melbourne": "AU",
    "Pacific/Auckland": "NZ",
  };

  return timezoneToCountry[timezone] || "AQ";
}
