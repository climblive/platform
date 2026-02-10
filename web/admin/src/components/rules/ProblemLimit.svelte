<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import WaCheckbox from "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import "@awesome.me/webawesome/dist/components/number-input/number-input.js";
  import { checked, GenericForm, name } from "@climblive/lib/forms";
  import type { Contest, ContestPatch } from "@climblive/lib/models";
  import { patchContestMutation } from "@climblive/lib/queries";
  import { z } from "@climblive/lib/utils";
  import RuleOptionCard from "../RuleOptionCard.svelte";
  import { doSubmit } from "../RulesEditor.svelte";

  interface Props {
    contest: Contest;
  }

  const { contest }: Props = $props();

  const patchContest = $derived(patchContestMutation(contest.id));

  let enabled = $derived(contest.qualifyingProblems > 0);

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
      description="Only count a configurable number of the hardest problems towards each contender's total score."
    >
      {#snippet header()}
        <wa-checkbox
          size="small"
          onchange={(event: InputEvent) => {
            const checkbox = event.target as WaCheckbox;
            enabled = checkbox.checked;

            setTimeout(() => form.requestSubmit());
          }}
          {@attach checked(enabled)}
        ></wa-checkbox>
      {/snippet}
      {#snippet footer()}
        <div class="controls">
          {#if enabled}
            {#if contest.pooledPoints}
              <wa-callout variant="warning" size="small">
                <wa-icon slot="icon" name="triangle-exclamation"></wa-icon>
                This setting in combination with pooled points may lead to unintuitive
                scoring results.
              </wa-callout>
            {/if}

            <wa-number-input
              size="small"
              {@attach name("qualifyingProblems")}
              label="Limit"
              required
              min={0}
              max={65536}
              defaultValue={contest.qualifyingProblems || 10}
            ></wa-number-input>

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
  wa-number-input {
    width: 100%;
  }

  wa-callout {
    width: 100%;
    margin-bottom: var(--wa-space-xs);
  }

  .controls {
    display: flex;
    flex-wrap: wrap;
    gap: var(--wa-space-xs);
    align-items: end;
  }
</style>
