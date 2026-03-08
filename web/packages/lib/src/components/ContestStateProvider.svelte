<script lang="ts">
  import { add, isBefore } from "date-fns";
  import { onMount, type Snippet } from "svelte";
  import { getCompClassQuery, getContestQuery } from "../queries";
  import type { ContestState } from "../types";
  import { SyncedTime } from "../utils";

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
