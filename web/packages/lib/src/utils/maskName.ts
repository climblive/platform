export function maskScrubbedName(contenderId: number): string {
  let hash = contenderId;
  hash = ((hash >> 16) ^ hash) * 0x45d9f3b;
  hash = ((hash >> 16) ^ hash) * 0x45d9f3b;
  hash = (hash >> 16) ^ hash;

  const firstNameLength = 3 + (Math.abs(hash) % 8);
  const lastNameLength = 4 + (Math.abs(hash >> 8) % 10);

  const firstName = "*".repeat(firstNameLength);
  const lastName = "*".repeat(lastNameLength);

  return `${firstName} ${lastName}`;
}
