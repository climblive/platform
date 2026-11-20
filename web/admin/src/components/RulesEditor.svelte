<script lang="ts" module>
  export const doSubmit = (
    mutation: ReturnType<typeof patchCompetitionMutation>,
    patch: CompetitionPatch,
  ) => {
    mutation.mutate(patch, {
      onError: () => toastUnexpectedError("Failed to update rules."),
    });
  };
</script>

<script lang="ts">
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import "@awesome.me/webawesome/dist/components/radio/radio.js";
  import type { Competition, CompetitionPatch } from "@climblive/lib/models";
  import { patchCompetitionMutation } from "@climblive/lib/queries";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import RuleOptionCard from "./RuleOptionCard.svelte";
  import Finalists from "./rules/Finalists.svelte";
  import PooledPoints from "./rules/PooledPoints.svelte";
  import MaxAttempts from "./rules/MaxAttempts.svelte";
  import PointDeduction from "./rules/PointDeduction.svelte";
  import ProblemLimit from "./rules/ProblemLimit.svelte";

  interface Props {
    competition: Competition;
  }

  const { competition }: Props = $props();

  const patchCompetition = $derived(patchCompetitionMutation(competition.id));

  const handleUsePointsChange = (usePoints: boolean) => {
    const patch: CompetitionPatch = {
      usePoints,
    };

    if (!usePoints) {
      patch.qualifyingProblems = 0;
      patch.pooledPoints = false;
      patch.maxAttempts = 0;
      patch.pointDeduction = 0;
    }

    patchCompetition.mutate(patch, {
      onError: () => toastUnexpectedError("Failed to update rules."),
    });
  };
</script>

<h3>Ranking method</h3>
<section>
  <RuleOptionCard
    title="Points"
    description="Competitors are ranked based on the total points scored across all problems."
    tags={["Recommended"]}
  >
    {#snippet header()}
      <wa-radio
        onclick={() => handleUsePointsChange(true)}
        size="s"
        checked={competition.usePoints ? true : undefined}
      ></wa-radio>
    {/snippet}
  </RuleOptionCard>

  <RuleOptionCard
    title="Attempts"
    description="Competitors are ranked by most tops, then fewest attempts to tops, most zone 2s, fewest attempts to zone 2s, most zone 1s, and finally fewest attempts to zone 1s. Each criterion is used only to break ties in the preceding criteria, so fewer attempts to tops rank ahead of additional zones."
  >
    {#snippet header()}
      <wa-radio
        onclick={() => handleUsePointsChange(false)}
        size="s"
        checked={!competition.usePoints ? true : undefined}
      ></wa-radio>
    {/snippet}
  </RuleOptionCard>
</section>

<h3>Options</h3>
<section>
  <Finalists {competition} />

  <PooledPoints {competition} />

  <ProblemLimit {competition} />

  <MaxAttempts {competition} />

  <PointDeduction {competition} />
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
