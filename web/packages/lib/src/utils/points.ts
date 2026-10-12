export const localizedPointsValue = (points: number) => {
  const rounded = Math.round(points * 10) / 10;
  return rounded.toLocaleString();
};
