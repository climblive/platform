<script lang="ts">
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import WaCheckbox from "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import "@awesome.me/webawesome/dist/components/switch/switch.js";
  import { checked, GenericForm, name } from "@climblive/lib/forms";
  import type { Contest, ContestPatch } from "@climblive/lib/models";
  import { patchContestMutation } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import * as z from "zod/v4";
  import RuleOptionCard from "./RuleOptionCard.svelte";

  interface Props {
    contest: Contest;
  }

  const { contest }: Props = $props();

  const patchContest = $derived(patchContestMutation(contest.id));

  let enableProblemLimit = $derived(contest.qualifyingProblems > 0);
  let enableFinalists = $derived(contest.finalists > 0);

  const doSubmit = (patch: ContestPatch) => {
    patchContest.mutate(patch, {
      onError: () => toastError("Failed to update contest rules."),
    });
  };
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
  <GenericForm
    submit={(value) =>
      doSubmit({ qualifyingProblems: value.qualifyingProblems ?? 0 })}
    schema={z.object({
      qualifyingProblems: z.coerce.number().min(0).max(65536).optional(),
    })}
  >
    {#snippet children(form)}
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

              setTimeout(() => form.requestSubmit());
            }}
            {@attach checked(enableProblemLimit)}
          ></wa-checkbox>
        {/snippet}
        {#snippet footer()}
          <div class="controls">
            {#if enableProblemLimit}
              <wa-input
                size="small"
                {@attach name("qualifyingProblems")}
                label="Limit"
                type="number"
                required
                min={0}
                max={65536}
                defaultValue={contest.qualifyingProblems || 10}
              ></wa-input>

              <wa-button
                type="submit"
                size="small"
                appearance="outlined"
                loading={patchContest.isPending}>Save</wa-button
              >
            {/if}
          </div>
        {/snippet}
      </RuleOptionCard>
    {/snippet}
  </GenericForm>

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

  <GenericForm
    submit={(value) => doSubmit({ finalists: value.finalists ?? 0 })}
    schema={z.object({
      finalists: z.coerce.number().min(0).max(65536).optional(),
    })}
  >
    {#snippet children(form)}
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

              setTimeout(() => form.requestSubmit());
            }}
            {@attach checked(enableFinalists)}
          ></wa-checkbox>
        {/snippet}
        {#snippet footer()}
          <div class="controls">
            {#if enableFinalists}
              <wa-input
                size="small"
                {@attach name("finalists")}
                label="Finalists"
                type="number"
                required
                min={0}
                max={65536}
                defaultValue={contest.finalists || 7}
                disabled={!enableFinalists}
              ></wa-input>

              <wa-button
                type="submit"
                size="small"
                appearance="outlined"
                loading={patchContest.isPending}>Save</wa-button
              >
            {/if}
          </div>
        {/snippet}
      </RuleOptionCard>
    {/snippet}
  </GenericForm>
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

  wa-input {
    width: 100%;
  }

  .controls {
    display: flex;
    gap: var(--wa-space-xs);
    align-items: end;
  }
</style>
