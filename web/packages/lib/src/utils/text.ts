export const ordinalSuperscript = (value: number): string => {
  switch (true) {
    case value === 1:
      return "st";
    case value === 2:
      return "nd";
    case value === 3:
      return "rd";
    case value <= 20:
      return `th`;
    case value % 10 === 1:
      return `st`;
    case value % 10 === 2:
      return `nd`;
    case value % 10 === 3:
      return `rd`;
    default:
      return `th`;
  }
};
