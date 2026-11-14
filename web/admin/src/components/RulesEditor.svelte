<script lang="ts" module>
  export const doSubmit = (
    mutation: ReturnType<typeof patchContestMutation>,
    patch: ContestPatch,
  ) => {
    mutation.mutate(patch, {
      onError: () => toastUnexpectedError("Failed to update rules."),
    });
  };
</script>

<script lang="ts">
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import "@awesome.me/webawesome/dist/components/radio/radio.js";
  import type { Contest, ContestPatch } from "@climblive/lib/models";
  import { patchContestMutation } from "@climblive/lib/queries";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import RuleOptionCard from "./RuleOptionCard.svelte";
  import Finalists from "./rules/Finalists.svelte";
  import PooledPoints from "./rules/PooledPoints.svelte";
  import MaxAttempts from "./rules/MaxAttempts.svelte";
  import PointDeduction from "./rules/PointDeduction.svelte";
  import ProblemLimit from "./rules/ProblemLimit.svelte";

  interface Props {
    contest: Contest;
  }

  const { contest }: Props = $props();

  const patchContest = $derived(patchContestMutation(contest.id));

  const handleUsePointsChange = (usePoints: boolean) => {
    const patch: ContestPatch = {
      usePoints,
    };

    if (!usePoints) {
      patch.qualifyingProblems = 0;
      patch.pooledPoints = false;
      patch.maxAttempts = 0;
      patch.pointDeduction = 0;
    }

    patchContest.mutate(patch, {
      onError: () => toastUnexpectedError("Failed to update rules."),
    });
  };
</script>

<h3>Ranking method</h3>
<section>
  <RuleOptionCard
    title="Points"
    description="Competitors are ranked based on the total points scored across all problems."
  >
    {#snippet header()}
      <wa-radio
        onclick={() => handleUsePointsChange(true)}
        size="s"
        checked={contest.usePoints ? true : undefined}
      ></wa-radio>
    {/snippet}
  </RuleOptionCard>

  <RuleOptionCard
    title="Attempts"
    description="Competitors are ranked based on the number of attempts needed to complete problems."
    tag="Beta"
  >
    {#snippet header()}
      <wa-radio
        onclick={() => handleUsePointsChange(false)}
        size="s"
        checked={!contest.usePoints ? true : undefined}
      ></wa-radio>
    {/snippet}
  </RuleOptionCard>
</section>

<h3>Options</h3>
<section>
  <Finalists {contest} />

  <PooledPoints {contest} />

  <ProblemLimit {contest} />

  <MaxAttempts {contest} />

  <PointDeduction {contest} />
</section>

<style>
  section {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--wa-space-m);
  }

  @media screen and (max-width: 768px) {
    section {
      grid-template-columns: 1fr;
    }
  }
</style>
