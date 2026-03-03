<script lang="ts">
  import { add, isBefore } from "date-fns";
  import { onDestroy, type Snippet } from "svelte";
  import { getCompClassQuery, getContestQuery } from "../queries";
  import type { ContestState } from "../types";

  interface Props {
    contestId: number;
    compClassId?: number;
    children?: Snippet<[{ contestState: ContestState }]>;
  }

  const { contestId, compClassId = undefined, children }: Props = $props();

  const contestQuery = $derived(getContestQuery(contestId));
  const compClassQuery = $derived(
    compClassId ? getCompClassQuery(compClassId) : undefined,
  );

  const contest = $derived(contestQuery.data);
  const compClass = $derived(compClassQuery?.data);

  const startTime = $derived(
    compClass?.timeBegin ?? contest?.timeBegin ?? new Date(8640000000000000),
  );
  const endTime = $derived(
    compClass?.timeEnd ?? contest?.timeEnd ?? new Date(-8640000000000000),
  );
  const gracePeriodEndTime = $derived(
    contest
      ? add(endTime, {
          minutes: (contest.gracePeriod ?? 0) / (1_000_000_000 * 60),
        })
      : undefined,
  );

  let contestState: ContestState = $state("NOT_STARTED");
  let intervalTimerId: number;

  $effect(() => {
    if (startTime && endTime) {
      clearTimeout(intervalTimerId);
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
