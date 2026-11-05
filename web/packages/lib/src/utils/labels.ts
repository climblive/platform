export const compareLabels = (l1: string, l2: string) => {
  return l1.localeCompare(l2, undefined, {
    numeric: true,
    sensitivity: "base",
  });
};

export const renderLabel = (label: string) => {
  if (label.match(/^[0-9]+$/)) {
    return "#" + label;
  }

  return label;
};

export const labelRegex = /^(\p{L}*)([0-9]+)(\p{L}*)$/u;
