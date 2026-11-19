export const roundPoints = (points: number) => Math.round(points * 10) / 10;

let compactFormatter: Intl.NumberFormat | undefined;

export const pointsNumberFormat = (
  points: number,
  compactThreshold: number = Number.MAX_SAFE_INTEGER,
) => {
  const rounded = roundPoints(points);
  const compact = points > compactThreshold;
  let fractionDigits = Number.isInteger(rounded) ? 0 : 1;

  if (compact) {
    compactFormatter ??= new Intl.NumberFormat(undefined, {
      notation: "compact",
      maximumFractionDigits: 1,
    });

    fractionDigits =
      compactFormatter
        .formatToParts(rounded)
        .find((part) => part.type === "fraction")?.value.length ?? 0;
  }

  return new Intl.NumberFormat(undefined, {
    notation: compact ? "compact" : "standard",
    minimumFractionDigits: fractionDigits,
    maximumFractionDigits: fractionDigits,
  });
};

export const localizedPointsValue = (
  points: number,
  compactThreshold: number = Number.MAX_SAFE_INTEGER,
) =>
  pointsNumberFormat(points, compactThreshold).format(roundPoints(points)) +
  "p";
