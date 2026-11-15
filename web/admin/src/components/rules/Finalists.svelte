<script lang="ts">
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import WaCheckbox from "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import "@awesome.me/webawesome/dist/components/number-input/number-input.js";
  import { SaveIndicator } from "@climblive/lib/components";
  import { checked, GenericForm, name } from "@climblive/lib/forms";
  import type { Competition, CompetitionPatch } from "@climblive/lib/models";
  import { patchCompetitionMutation } from "@climblive/lib/queries";
  import { debounce, z } from "@climblive/lib/utils";
  import RuleOptionCard from "../RuleOptionCard.svelte";
  import { doSubmit } from "../RulesEditor.svelte";

  interface Props {
    competition: Competition;
  }

  const { competition }: Props = $props();

  const patchCompetition = patchCompetitionMutation(competition.id);

  let enabled = $derived(competition.finalists > 0);

  const formSchema = z.object({
    finalists: z.coerce.number().min(0).max(65536).optional(),
  });

  const debouncedSubmit = debounce(
    (form: HTMLFormElement) => form.requestSubmit(),
    1000,
  );

  const handleSubmit = (value: Partial<CompetitionPatch>) =>
    doSubmit(patchCompetition, { finalists: value.finalists ?? 0 });
</script>

<GenericForm submit={handleSubmit} schema={formSchema}>
  {#snippet children(form)}
    <RuleOptionCard
      title="Finalists"
      description="The number of competitors that will proceed to finals.
    There might be additional finalists in case of ties."
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
        {#if patchCompetition.isSuccess}
          <SaveIndicator />
        {/if}
      {/snippet}
      {#snippet footer()}
        <div class="controls">
          {#if enabled}
            <wa-number-input
              size="s"
              {@attach name("finalists")}
              label="Finalists"
              required
              min={0}
              max={65536}
              defaultValue={competition.finalists || 7}
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
    gap: var(--wa-space-xs);
    align-items: end;
  }
</style>
