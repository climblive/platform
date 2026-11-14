<script lang="ts">
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import WaCheckbox from "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import "@awesome.me/webawesome/dist/components/number-input/number-input.js";
  import { SaveIndicator } from "@climblive/lib/components";
  import { GenericForm, name } from "@climblive/lib/forms";
  import type { Competition, CompetitionPatch } from "@climblive/lib/models";
  import { patchCompetitionMutation } from "@climblive/lib/queries";
  import { debounce, z } from "@climblive/lib/utils";
  import RuleOptionCard from "../RuleOptionCard.svelte";
  import { doSubmit } from "../RulesEditor.svelte";

  interface Props {
    competition: Competition;
  }

  const { competition }: Props = $props();

  const competitionId = $derived(competition.id);
  const patchCompetition = $derived(patchCompetitionMutation(competitionId));

  let enabled = $derived(competition.maxAttempts > 0);

  const formSchema = z.object({
    maxAttempts: z.coerce.number().int().min(1).max(999).optional(),
  });

  const debouncedSubmit = debounce(
    (form: HTMLFormElement) => form.requestSubmit(),
    1000,
  );

  const handleSubmit = (value: Partial<CompetitionPatch>) =>
    doSubmit(patchCompetition, {
      maxAttempts: value.maxAttempts ?? 0,
    });
</script>

<GenericForm schema={formSchema} submit={handleSubmit}>
  {#snippet children(form)}
    <RuleOptionCard
      title="Max attempts"
      disabled={!competition.usePoints}
      description="Each competitor has a fixed number of attempts per problem. Attempts above that limit yield 0 points for that problem."
      tag="Beta"
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
            {@attach name("maxAttempts")}
            label="Attempts"
            required
            min={1}
            step={1}
            max={999}
            defaultValue={competition.maxAttempts || 10}
            oninput={() => debouncedSubmit(form)}
          ></wa-number-input>
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
