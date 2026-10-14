export const localizedPointsValue = (points: number, short: boolean) => {
  const rounded = Math.round(points * 10) / 10;

  const options: Intl.NumberFormatOptions = {
    maximumFractionDigits: 1,
  };

  if (short) {
    options.notation = "compact";
  }

  return rounded.toLocaleString(undefined, options) + "p";
};
