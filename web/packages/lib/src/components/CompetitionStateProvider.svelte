<script lang="ts">
  import { add, isBefore } from "date-fns";
  import { onMount, type Snippet } from "svelte";
  import { getCompClassQuery, getCompetitionQuery } from "../queries";
  import type { CompetitionState } from "../types";
  import { SyncedTime } from "../utils";

  interface Props {
    competitionId: number;
    compClassId?: number;
    children?: Snippet<
      [{ competitionState: CompetitionState; progress: number }]
    >;
  }

  const { competitionId, compClassId, children }: Props = $props();

  const competitionQuery = $derived(getCompetitionQuery(competitionId));
  const compClassQuery = $derived(
    compClassId ? getCompClassQuery(compClassId) : undefined,
  );

  const competition = $derived(competitionQuery.data);
  const compClass = $derived(compClassQuery?.data);

  const startTime = $derived(
    compClass?.timeBegin ??
      competition?.timeBegin ??
      new Date(8640000000000000),
  );
  const endTime = $derived(
    compClass?.timeEnd ?? competition?.timeEnd ?? new Date(-8640000000000000),
  );
  const gracePeriodEndTime = $derived(
    competition
      ? add(endTime, {
          minutes: (competition.gracePeriod ?? 0) / (1_000_000_000 * 60),
        })
      : undefined,
  );

  const time = new SyncedTime(1_000);

  onMount(() => {
    time.start();

    return () => time.stop();
  });

  const computeState = (): CompetitionState => {
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

  const competitionState: CompetitionState = $derived(computeState());

  const progress = $derived.by(() => {
    const total = endTime.getTime() - startTime.getTime();
    if (total <= 0) {
      return 0;
    }

    const elapsed = time.current.getTime() - startTime.getTime();
    return Math.min(100, Math.max(0, (elapsed / total) * 100));
  });
</script>

{@render children?.({ competitionState, progress })}
