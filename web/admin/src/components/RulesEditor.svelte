<script lang="ts" module>
  import * as z from "zod/v4";

  export const formSchema = z.object({
    qualifyingProblems: z.coerce.number().min(0).max(65536),
    finalists: z.coerce.number().min(0).max(65536),
  });
</script>

<script lang="ts">
  import { serialize } from "@awesome.me/webawesome";
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import WaCheckbox from "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import type WaInput from "@awesome.me/webawesome/dist/components/input/input.js";
  import "@awesome.me/webawesome/dist/components/switch/switch.js";
  import { checked, name } from "@climblive/lib/forms";
  import type { Contest } from "@climblive/lib/models";
  import { patchContestMutation } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import RuleOptionCard from "./RuleOptionCard.svelte";

  interface Props {
    contest: Contest;
  }

  const { contest }: Props = $props();

  const patchContest = $derived(patchContestMutation(contest.id));

  let enableProblemLimit = $derived(contest.qualifyingProblems > 0);
  let enableFinalists = $derived(contest.finalists > 0);

  const handleSubmit = (event: SubmitEvent) => {
    event.preventDefault();

    if (!form) {
      return;
    }

    const data = serialize(form);
    const result = formSchema.safeParse(data);

    if (result.success) {
      patchContest.mutate(
        {
          qualifyingProblems: result.data.qualifyingProblems,
          finalists: result.data.finalists,
        },
        {
          onError: () => toastError("Failed to update contest rules."),
        },
      );
    }
  };

  let form: HTMLFormElement | undefined = $state();
  let qualifyingProblemsInput: WaInput | undefined = $state();
  let finalistsInput: WaInput | undefined = $state();
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

<form bind:this={form} onsubmit={handleSubmit}>
  <section>
    <RuleOptionCard
      title="Problem limit"
      description="Only count a configurable number of the hardest problems towards each contenders total score."
    >
      {#snippet header()}
        <wa-checkbox
          size="small"
          onchange={(event: InputEvent) => {
            const checkbox = event.target as WaCheckbox;
            enableProblemLimit = checkbox.checked;

            if (qualifyingProblemsInput) {
              qualifyingProblemsInput.value = checkbox.checked ? "10" : "0";
            }
          }}
          {@attach checked(enableProblemLimit)}
        ></wa-checkbox>
      {/snippet}
      {#snippet footer()}
        {#if enableProblemLimit}
          <wa-input
            bind:this={qualifyingProblemsInput}
            size="small"
            {@attach name("qualifyingProblems")}
            label="Limit"
            type="number"
            required
            min={0}
            max={65536}
            value={10}
            disabled={!enableProblemLimit}
          ></wa-input>
        {:else}
          <input type="hidden" {@attach name("qualifyingProblems")} value="0" />
        {/if}
      {/snippet}
    </RuleOptionCard>

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

    <RuleOptionCard
      title="Finalists"
      description="Number of contenders that should proceed to the finals.
    There might be additional finalists in case of ties."
    >
      {#snippet header()}
        <wa-checkbox
          size="small"
          onchange={(event: InputEvent) => {
            const checkbox = event.target as WaCheckbox;
            enableFinalists = checkbox.checked;
            if (finalistsInput) {
              finalistsInput.value = checkbox.checked ? "7" : "0";
            }
          }}
          {@attach checked(enableFinalists)}
        ></wa-checkbox>
      {/snippet}
      {#snippet footer()}
        {#if enableFinalists}
          <wa-input
            bind:this={finalistsInput}
            size="small"
            {@attach name("finalists")}
            label="Finalists"
            type="number"
            required
            min={0}
            max={65536}
            value={7}
            disabled={!enableFinalists}
          ></wa-input>
        {:else}
          <input type="hidden" {@attach name("finalists")} value="0" />
        {/if}
      {/snippet}
    </RuleOptionCard>
  </section>

  <wa-button type="submit" variant="primary" loading={patchContest.isPending}
    >Save changes</wa-button
  >
</form>

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

  wa-button {
    margin-top: var(--wa-space-m);
  }
</style>
