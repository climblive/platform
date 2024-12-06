export const asOrdinal = (value: number): string => {
  switch (true) {
    case value === 1:
      return "1st";
    case value === 2:
      return "2nd";
    case value === 3:
      return "3rd";
    case value <= 20:
      return `${value}th`;
    case value % 10 === 1:
      return `${value}st`;
    case value % 10 === 2:
      return `${value}nd`;
    case value % 10 === 3:
      return `${value}rd`;
    default:
      return `${value}th`;
  }
};
