import countriesData from "../../../../../countries.json";

export interface Country {
  code: string; // ISO 3166-1 alpha-2 code
  name: string;
}

export const countries: Country[] = countriesData as Country[];

const REGIONAL_INDICATOR_OFFSET = 127397;

const countryMap = new Map(countries.map((c) => [c.code, c.name]));

export function getFlag(countryCode: string): string {
  const codePoints = countryCode
    .toUpperCase()
    .split("")
    .map((char) => REGIONAL_INDICATOR_OFFSET + char.charCodeAt(0));

  return String.fromCodePoint(...codePoints);
}

export function getCountryName(countryCode: string): string {
  return countryMap.get(countryCode) || "";
}
