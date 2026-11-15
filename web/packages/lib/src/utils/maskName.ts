export function maskScrubbedName(competitorId: number): string {
  let hash = competitorId;
  hash = ((hash >> 16) ^ hash) * 0x45d9f3b;
  hash = ((hash >> 16) ^ hash) * 0x45d9f3b;
  hash = (hash >> 16) ^ hash;

  return `${Math.abs(hash)}`;
}
