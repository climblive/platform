<script lang="ts">
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import { SaveIndicator } from "@climblive/lib/components";
  import { GenericForm, name } from "@climblive/lib/forms";
  import type { Competition, CompetitionPatch } from "@climblive/lib/models";
  import { patchCompetitionMutation } from "@climblive/lib/queries";
  import { z } from "@climblive/lib/utils";
  import RuleOptionCard from "../RuleOptionCard.svelte";
  import { doSubmit } from "../RulesEditor.svelte";

  interface Props {
    competition: Competition;
  }

  const { competition }: Props = $props();

  const patchCompetition = patchCompetitionMutation(competition.id);

  const formSchema = z.object({
    pooledPoints: z.coerce.boolean().default(false),
  });

  const handleSubmit = (value: Partial<CompetitionPatch>) =>
    doSubmit(patchCompetition, { pooledPoints: value.pooledPoints ?? false });
</script>

<GenericForm submit={handleSubmit} schema={formSchema}>
  {#snippet children(form)}
    <RuleOptionCard
      title="Pooled points"
      disabled={!competition.usePoints}
      description="Points for completed problems and zones are split by percentages. A boulder worth 1000 points with two tops will give each competitor 500 points. If a third competitor also tops the boulder, then all three will receive 333.3 points instead. Minimum point value is 0.1 point."
      tags={["New"]}
    >
      {#snippet header()}
        <wa-checkbox
          size="small"
          disabled={!competition.usePoints}
          {@attach name("pooledPoints")}
          onchange={() => setTimeout(() => form.requestSubmit())}
          checked={competition.pooledPoints}
        ></wa-checkbox>
      {/snippet}
      {#snippet indicator()}
        {#if patchCompetition.isSuccess}
          <SaveIndicator />
        {/if}
      {/snippet}
    </RuleOptionCard>
  {/snippet}
</GenericForm>
