<script lang="ts">
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import WaCheckbox from "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import { checked, GenericForm, name } from "@climblive/lib/forms";
  import type { Contest, ContestPatch } from "@climblive/lib/models";
  import { patchContestMutation } from "@climblive/lib/queries";
  import * as z from "zod/v4";
  import RuleOptionCard from "../RuleOptionCard.svelte";
  import { doSubmit } from "../RulesEditor.svelte";

  interface Props {
    contest: Contest;
  }

  const { contest }: Props = $props();

  const patchContest = $derived(patchContestMutation(contest.id));

  let enableProblemLimit = $derived(contest.qualifyingProblems > 0);

  const formSchema = z.object({
    qualifyingProblems: z.coerce.number().min(0).max(65536).optional(),
  });

  const handleSubmit = (value: Partial<ContestPatch>) =>
    doSubmit(patchContest, {
      qualifyingProblems: value.qualifyingProblems ?? 0,
    });
</script>

<GenericForm schema={formSchema} submit={handleSubmit}>
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

<style>
  wa-input {
    width: 100%;
  }

  .controls {
    display: flex;
    gap: var(--wa-space-xs);
    align-items: end;
  }
</style>
