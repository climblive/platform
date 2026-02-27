<script lang="ts">
  import { isBefore } from "date-fns";
  import { onMount, type Snippet } from "svelte";
  import type { ContestState } from "../types";
  import { SyncedTime } from "../utils";

  interface Props {
    startTime: Date;
    endTime: Date;
    gracePeriodEndTime?: Date;
    children?: Snippet<[{ contestState: ContestState }]>;
  }

  const {
    startTime,
    endTime,
    gracePeriodEndTime = undefined,
    children,
  }: Props = $props();

  const time = new SyncedTime(1_000);

  onMount(() => {
    time.start();

    return () => time.stop();
  });

  const computeState = (): ContestState => {
    const now = new Date(time.current);
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

  const contestState: ContestState = $derived(computeState());
</script>

{@render children?.({ contestState })}
