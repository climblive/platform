export const localizedPointsValue = (
  points: number,
  short: boolean = false,
) => {
  const rounded = Math.round(points * 10) / 10;

  const options: Intl.NumberFormatOptions = {
    maximumFractionDigits: 1,
  };

  if (short && points > 9_999) {
    options.notation = "compact";
  }

  return rounded.toLocaleString(undefined, options) + "p";
};
