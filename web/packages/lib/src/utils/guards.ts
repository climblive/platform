export function isDefined<T>(
  value: NonNullable<T> | undefined | null,
): value is NonNullable<T> {
  return value !== null && value !== undefined;
}
