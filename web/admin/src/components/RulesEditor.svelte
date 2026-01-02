<script lang="ts">
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import "@awesome.me/webawesome/dist/components/switch/switch.js";
  import { name } from "@climblive/lib/forms";
  import RuleOptionCard from "./RuleOptionCard.svelte";

  interface Props {
    contest: {
      id: number;
    };
  }

  const { contest }: Props = $props();
</script>

<h3>Ranking Method</h3>
<section>
  <RuleOptionCard
    title="Points"
    description="Contenders are ranked based on the total points scored across all problems."
  >
    <wa-radio checked></wa-radio>
  </RuleOptionCard>

  <RuleOptionCard
    title="Attempts"
    description="Contenders are ranked based on the number of attempts needed to complete problems."
    disabled
    tag="Unavailable"
  >
    <wa-radio disabled></wa-radio>
  </RuleOptionCard>
</section>

<h3>Options</h3>

<section>
  <RuleOptionCard
    title="Problem limit"
    description="Only count a configurable number of the hardest problems towards each contenders total score."
  >
    <wa-checkbox></wa-checkbox>
    {#snippet footer()}
      <wa-input
        size="small"
        {@attach name("qualifyingProblems")}
        label="Limit"
        type="number"
        required
        min={0}
        max={65536}
        value={10}
      ></wa-input>
    {/snippet}
  </RuleOptionCard>

  <RuleOptionCard
    title="Pooled points"
    description="Points for each boulder are split by percentages among the ascensionists. A boulder worth 1000 points with two tops would give each contender 500 points. If a third contender also tops the boulder, then all three would receive 333 points instead."
    disabled
    tag="Unavailable"
  >
    <wa-checkbox disabled></wa-checkbox>
  </RuleOptionCard>

  <RuleOptionCard
    title="Finalists"
    description="Number of contenders that should proceed to the finals.
    There might be additional finalists in case of ties."
  >
    <wa-checkbox></wa-checkbox>
    {#snippet footer()}
      <wa-input
        size="small"
        {@attach name("finalists")}
        label="Finalists"
        type="number"
        required
        min={0}
        max={65536}
        value={7}
      ></wa-input>
    {/snippet}
  </RuleOptionCard>
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
