<script lang="ts" module>
  export const doSubmit = (
    mutation: ReturnType<typeof patchContestMutation>,
    patch: ContestPatch,
  ) => {
    mutation.mutate(patch, {
      onError: () => toastError("Failed to update contest rules."),
    });
  };
</script>

<script lang="ts">
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import "@awesome.me/webawesome/dist/components/switch/switch.js";
  import type { Contest, ContestPatch } from "@climblive/lib/models";
  import { patchContestMutation } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import RuleOptionCard from "./RuleOptionCard.svelte";
  import Finalists from "./rules/Finalists.svelte";
  import ProblemLimit from "./rules/ProblemLimit.svelte";

  interface Props {
    contest: Contest;
  }

  const { contest }: Props = $props();
</script>

<h3>Ranking Method</h3>
<section>
  <RuleOptionCard
    title="Points"
    description="Contenders are ranked based on the total points scored across all problems."
  >
    {#snippet header()}
      <wa-radio size="small" checked></wa-radio>
    {/snippet}
  </RuleOptionCard>

  <RuleOptionCard
    title="Attempts"
    description="Contenders are ranked based on the number of attempts needed to complete problems."
    disabled
    tag="Upcoming"
  >
    {#snippet header()}
      <wa-radio size="small" disabled></wa-radio>
    {/snippet}
  </RuleOptionCard>
</section>

<h3>Options</h3>

<section>
  <ProblemLimit {contest} />

  <RuleOptionCard
    title="Pooled points"
    description="Points for each boulder are split by percentages among the ascensionists. A boulder worth 1000 points with two tops would give each contender 500 points. If a third contender also tops the boulder, then all three would receive 333 points instead."
    disabled
    tag="Upcoming"
  >
    {#snippet header()}
      <wa-checkbox size="small" disabled></wa-checkbox>
    {/snippet}
  </RuleOptionCard>

  <Finalists {contest} />
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
