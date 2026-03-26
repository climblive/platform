export function maskScrubbedName(contenderId: number): string {
  let hash = contenderId;
  hash = ((hash >> 16) ^ hash) * 0x45d9f3b;
  hash = ((hash >> 16) ^ hash) * 0x45d9f3b;
  hash = (hash >> 16) ^ hash;

  return `Anon ${Math.abs(hash)}`;
}
