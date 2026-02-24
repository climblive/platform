export type ContestState = "NOT_STARTED" | "RUNNING" | "GRACE_PERIOD" | "ENDED";

export const contestStateToString = (state: ContestState): string => {
  switch (state) {
    case "NOT_STARTED":
      return "Not started";
    case "RUNNING":
      return "Running";
    case "GRACE_PERIOD":
    case "ENDED":
      return "Ended";
  }
};
