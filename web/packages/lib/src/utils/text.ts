export const asOrdinal = (value: number): string => {
  switch (value) {
    case 1:
      return "1st";
    case 2:
      return "2nd";
    case 3:
      return "3rd";
    default:
      return `${value}th`;
  }
};

