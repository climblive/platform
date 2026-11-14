export const roundPoints = (points: number) => Math.floor(points * 10) / 10;

export const pointsNumberFormat = (
  points: number,
  compactThreshold: number = Number.MAX_SAFE_INTEGER,
) => {
  const options: Intl.NumberFormatOptions = {
    maximumFractionDigits: 1,
  };

  if (points > compactThreshold) {
    options.notation = "compact";
  }

  return new Intl.NumberFormat(undefined, options);
};

export const localizedPointsValue = (
  points: number,
  compactThreshold: number = Number.MAX_SAFE_INTEGER,
) =>
  pointsNumberFormat(points, compactThreshold).format(roundPoints(points)) +
  "p";
