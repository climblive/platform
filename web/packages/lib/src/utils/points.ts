export const localizedPointsValue = (
  points: number,
  compactThreshold: number = Number.MAX_SAFE_INTEGER,
) => {
  const rounded = Math.round(points * 10) / 10;

  const options: Intl.NumberFormatOptions = {
    maximumFractionDigits: 1,
  };

  if (points > compactThreshold) {
    options.notation = "compact";
  }

  return rounded.toLocaleString(undefined, options) + "p";
};
