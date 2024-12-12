<!-- @migration-task Error while migrating Svelte code: can't migrate `let state: ContestState = "NOT_STARTED";` to `$state` because there's a variable named state.
     Rename the variable and try again or migrate by hand. -->
<script lang="ts">
  import { isBefore } from "date-fns";
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
      tick();
    }
  }

  const computeState = (): ContestState => {
    const now = new Date();
    now.setMilliseconds(0);

    switch (true) {
      case isBefore(now, startTime):
        return "NOT_STARTED";
      case isBefore(now, endTime):
        return "RUNNING";
      case gracePeriodEndTime && isBefore(now, gracePeriodEndTime):
        return "GRACE_PERIOD";
      default:
        return "ENDED";
    }
  };

  const tick = () => {
    state = computeState();

    const firefoxEarlyWakeUpCompensation = 1;

    const drift = Date.now() % 1_000;
    const next = 1_000 - drift + firefoxEarlyWakeUpCompensation;

    intervalTimerId = setTimeout(tick, next);
  };

  onDestroy(() => {
    clearTimeout(intervalTimerId);
  });
</script>

<slot {state} />
