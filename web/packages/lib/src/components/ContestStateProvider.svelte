<script lang="ts">
  import { isBefore } from "date-fns";
  import { onMount, type Snippet } from "svelte";
  import type { ContestState } from "../types";
  import { SyncedTime } from "../utils";

  interface Props {
    startTime: Date;
    endTime: Date;
    gracePeriodEndTime?: Date;
    children?: Snippet<[{ contestState: ContestState; progress: number }]>;
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

  const progress = $derived.by(() => {
    const total = endTime.getTime() - startTime.getTime();
    if (total <= 0) {
      return 0;
    }

    const elapsed = time.current.getTime() - startTime.getTime();
    return Math.min(100, Math.max(0, (elapsed / total) * 100));
  });
</script>

{@render children?.({ contestState, progress })}
