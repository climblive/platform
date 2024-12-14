<script lang="ts">
  import { isBefore } from "date-fns";
  import { onDestroy, type Snippet } from "svelte";
  import type { ContestState } from "../types";

  interface Props {
    startTime: Date;
    endTime: Date;
    gracePeriodEndTime: Date | undefined;
    children?: Snippet<[{ contestState: ContestState }]>;
  }

  let { startTime, endTime, gracePeriodEndTime, children }: Props = $props();

  let contestState: ContestState = $state("NOT_STARTED");
  let intervalTimerId: number;

  $effect(() => {
    if (startTime && endTime) {
      clearInterval(intervalTimerId);
      tick();
    }
  });

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
    contestState = computeState();

    const firefoxEarlyWakeUpCompensation = 1;

    const drift = Date.now() % 1_000;
    const next = 1_000 - drift + firefoxEarlyWakeUpCompensation;

    intervalTimerId = setTimeout(tick, next);
  };

  onDestroy(() => {
    clearTimeout(intervalTimerId);
  });
</script>

{@render children?.({ contestState })}
