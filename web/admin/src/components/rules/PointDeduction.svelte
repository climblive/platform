<script lang="ts">
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import WaCheckbox from "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import "@awesome.me/webawesome/dist/components/number-input/number-input.js";
  import { SaveIndicator } from "@climblive/lib/components";
  import { GenericForm, name } from "@climblive/lib/forms";
  import type { Competition, CompetitionPatch } from "@climblive/lib/models";
  import { patchCompetitionMutation } from "@climblive/lib/queries";
  import { debounce, roundPoints, z } from "@climblive/lib/utils";
  import RuleOptionCard from "../RuleOptionCard.svelte";
  import { doSubmit } from "../RulesEditor.svelte";

  interface Props {
    competition: Competition;
  }

  const { competition }: Props = $props();

  const competitionId = $derived(competition.id);
  const patchCompetition = $derived(patchCompetitionMutation(competitionId));

  let enabled = $derived(competition.pointDeduction > 0);

  const formSchema = z.object({
    pointDeduction: z.coerce
      .number()
      .min(0)
      .max(2 ** 31 - 1)
      .transform(roundPoints)
      .optional(),
  });

  const debouncedSubmit = debounce(
    (form: HTMLFormElement) => form.requestSubmit(),
    1000,
  );

  const handleSubmit = (value: Partial<CompetitionPatch>) =>
    doSubmit(patchCompetition, {
      pointDeduction: value.pointDeduction ?? 0,
    });
</script>

<GenericForm schema={formSchema} submit={handleSubmit}>
  {#snippet children(form)}
    <RuleOptionCard
      title="Point deduction"
      disabled={!competition.usePoints}
      description="Deduct points for each failed attempt before reaching the scored top or zone. The score cannot fall below 0."
      tags={["IFSC", "New"]}
    >
      {#snippet header()}
        <wa-checkbox
          size="s"
          disabled={!competition.usePoints}
          onchange={(event: InputEvent) => {
            const checkbox = event.target as WaCheckbox;
            enabled = checkbox.checked;

            setTimeout(() => form.requestSubmit());
          }}
          checked={enabled}
        ></wa-checkbox>
      {/snippet}
      {#snippet indicator()}
        {#if patchCompetition.isSuccess}
          <SaveIndicator />
        {/if}
      {/snippet}
      {#snippet footer()}
        {#if enabled}
          <wa-number-input
            size="s"
            disabled={!competition.usePoints}
            {@attach name("pointDeduction")}
            label="Deduction per failed attempt"
            required
            min={0}
            step={0.1}
            max={2 ** 31 - 1}
            defaultValue={competition.pointDeduction || 0.1}
            oninput={() => debouncedSubmit(form)}
          >
            <span slot="end">p</span>
          </wa-number-input>
        {/if}
      {/snippet}
    </RuleOptionCard>
  {/snippet}
</GenericForm>

<style>
  wa-number-input {
    width: 100%;
  }
</style>
