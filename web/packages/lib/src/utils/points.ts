export const roundPoints = (points: number) => Math.floor(points * 10) / 10;

export const localizedPointsValue = (
  points: number,
  compactThreshold: number = Number.MAX_SAFE_INTEGER,
) => {
  const rounded = roundPoints(points);

  const options: Intl.NumberFormatOptions = {
    maximumFractionDigits: 1,
  };

  if (points > compactThreshold) {
    options.notation = "compact";
  }

  return rounded.toLocaleString(undefined, options) + "p";
};
