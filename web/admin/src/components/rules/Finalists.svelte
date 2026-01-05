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

  let enableFinalists = $derived(contest.finalists > 0);

  const formSchema = z.object({
    finalists: z.coerce.number().min(0).max(65536).optional(),
  });

  const handleSubmit = (value: Partial<ContestPatch>) =>
    doSubmit(patchContest, { finalists: value.finalists ?? 0 });
</script>

<GenericForm submit={handleSubmit} schema={formSchema}>
  {#snippet children(form)}
    <RuleOptionCard
      title="Finalists"
      description="The number of contenders that will proceed to the finals.
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
