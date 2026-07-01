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
  import { debounce, z } from "@climblive/lib/utils";
  import { onDestroy } from "svelte";
  import RuleOptionCard from "../RuleOptionCard.svelte";
  import { doSubmit } from "../RulesEditor.svelte";

  interface Props {
    contest: Contest;
  }

  const { contest }: Props = $props();

  const patchContest = patchContestMutation(contest.id);

  let enabled = $derived(contest.qualifyingProblems > 0);
  let saved = $state(false);
  let savedTimer: ReturnType<typeof setTimeout> | undefined;

  onDestroy(() => clearTimeout(savedTimer));

  const formSchema = z.object({
    qualifyingProblems: z.coerce.number().min(0).max(65536).optional(),
  });

  const debouncedSubmit = debounce(
    (form: HTMLFormElement) => form.requestSubmit(),
    1000,
  );

  const handleSubmit = (value: Partial<ContestPatch>) =>
    doSubmit(
      patchContest,
      {
        qualifyingProblems: value.qualifyingProblems ?? 0,
      },
      () => {
        saved = true;
        clearTimeout(savedTimer);
        savedTimer = setTimeout(() => (saved = false), 2_000);
      },
    );
</script>

<GenericForm schema={formSchema} submit={handleSubmit}>
  {#snippet children(form)}
    <RuleOptionCard
      title="Problem limit"
      description="Only count the hardest problems towards each contender's total score."
    >
      {#snippet header()}
        <wa-checkbox
          size="s"
          onchange={(event: InputEvent) => {
            const checkbox = event.target as WaCheckbox;
            enabled = checkbox.checked;

            setTimeout(() => form.requestSubmit());
          }}
          {@attach checked(enabled)}
        ></wa-checkbox>
      {/snippet}
      {#snippet indicator()}
        {#if saved}
          <div class="indicator">
            <wa-icon name="check"></wa-icon>
            Saved
          </div>
        {/if}
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
              size="s"
              {@attach name("qualifyingProblems")}
              label="Limit"
              required
              min={0}
              max={65536}
              defaultValue={contest.qualifyingProblems || 10}
              oninput={() => debouncedSubmit(form)}
            ></wa-number-input>
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

  .indicator {
    margin-inline-start: auto;
    display: flex;
    align-items: center;
    gap: var(--wa-space-2xs);
    font-size: var(--wa-font-size-s);
    color: var(--wa-color-success-fill-loud);
    font-weight: var(--wa-font-weight-bold);
  }
</style>
