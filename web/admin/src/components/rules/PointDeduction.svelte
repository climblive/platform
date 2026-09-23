<script lang="ts">
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import WaCheckbox from "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import "@awesome.me/webawesome/dist/components/number-input/number-input.js";
  import { SaveIndicator } from "@climblive/lib/components";
  import { GenericForm, name } from "@climblive/lib/forms";
  import type { Contest, ContestPatch } from "@climblive/lib/models";
  import { patchContestMutation } from "@climblive/lib/queries";
  import { debounce, z } from "@climblive/lib/utils";
  import RuleOptionCard from "../RuleOptionCard.svelte";
  import { doSubmit } from "../RulesEditor.svelte";

  interface Props {
    contest: Contest;
  }

  const { contest }: Props = $props();

  const patchContest = $derived(patchContestMutation(contest.id));

  let enabled = $derived(contest.pointDeduction > 0);

  const formSchema = z.object({
    pointDeduction: z.coerce
      .number()
      .int()
      .min(0)
      .max(2 ** 31 - 1)
      .optional(),
  });

  const debouncedSubmit = debounce(
    (form: HTMLFormElement) => form.requestSubmit(),
    1000,
  );

  const handleSubmit = (value: Partial<ContestPatch>) =>
    doSubmit(patchContest, {
      pointDeduction: value.pointDeduction ?? 0,
    });
</script>

<GenericForm schema={formSchema} submit={handleSubmit}>
  {#snippet children(form)}
    <RuleOptionCard
      title="Point deduction"
      disabled={!contest.usePoints}
      description="Deduct points for each failed attempt before reaching the scored top or zone. The score cannot fall below 0."
    >
      {#snippet header()}
        <wa-checkbox
          size="s"
          disabled={!contest.usePoints}
          onchange={(event: InputEvent) => {
            const checkbox = event.target as WaCheckbox;
            enabled = checkbox.checked;

            setTimeout(() => form.requestSubmit());
          }}
          checked={enabled}
        ></wa-checkbox>
      {/snippet}
      {#snippet indicator()}
        {#if patchContest.isSuccess}
          <SaveIndicator />
        {/if}
      {/snippet}
      {#snippet footer()}
        <div class="controls">
          {#if enabled}
            <wa-number-input
              size="s"
              disabled={!contest.usePoints}
              {@attach name("pointDeduction")}
              label="Points per failed attempt"
              required
              min={0}
              step={1}
              max={2 ** 31 - 1}
              defaultValue={contest.pointDeduction || 1}
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

  .controls {
    display: flex;
    flex-wrap: wrap;
    gap: var(--wa-space-xs);
    align-items: end;
  }
</style>
