<script lang="ts">
  import { differenceInMilliseconds, isBefore } from "date-fns";
  import { onDestroy } from "svelte";
  import type { ContestState } from "../types";

  export let startTime: Date;
  export let endTime: Date;
  export let gracePeriodEndTime: Date | undefined = undefined;

  let state: ContestState = "NOT_STARTED";
  let intervalTimerId: number;

  $: {
    if (startTime && endTime) {
      clearInterval(intervalTimerId);
      computeState();
    }
  }

  const computeState = () => {
    const now = new Date();
    let durationUntilNextState: number = NaN;

    switch (true) {
      case isBefore(now, startTime):
        state = "NOT_STARTED";
        durationUntilNextState = differenceInMilliseconds(startTime, now);

        break;
      case isBefore(now, endTime):
        state = "RUNNING";
        durationUntilNextState = differenceInMilliseconds(endTime, now);

        break;
      case gracePeriodEndTime && isBefore(now, gracePeriodEndTime):
        state = "GRACE_PERIOD";
        durationUntilNextState = differenceInMilliseconds(
          gracePeriodEndTime,
          now,
        );

        break;
      default:
        state = "ENDED";
    }

    if (durationUntilNextState) {
      intervalTimerId = setTimeout(computeState, durationUntilNextState);
    }
  };

  onDestroy(() => {
    clearTimeout(intervalTimerId);
  });
</script>

<slot {state} />
