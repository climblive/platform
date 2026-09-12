<script lang="ts">
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import { SaveIndicator } from "@climblive/lib/components";
  import { GenericForm, name } from "@climblive/lib/forms";
  import type { Contest, ContestPatch } from "@climblive/lib/models";
  import { patchContestMutation } from "@climblive/lib/queries";
  import { z } from "@climblive/lib/utils";
  import RuleOptionCard from "../RuleOptionCard.svelte";
  import { doSubmit } from "../RulesEditor.svelte";

  interface Props {
    contest: Contest;
  }

  const { contest }: Props = $props();

  const patchContest = patchContestMutation(contest.id);

  const formSchema = z.object({
    pooledPoints: z.coerce.boolean(),
  });

  const handleSubmit = (value: Partial<ContestPatch>) =>
    doSubmit(patchContest, { pooledPoints: value.pooledPoints ?? false });
</script>

<GenericForm submit={handleSubmit} schema={formSchema}>
  {#snippet children(form)}
    <RuleOptionCard
      title="Pooled points"
      disabled={!contest.usePoints}
      description="Points for completed problems and zones are split by percentages. A boulder worth 1000 points with two tops will give each competitor 500 points. If a third competitor also tops the boulder, then all three will receive 333 points instead. Minimum point value is 1 point."
      tag="New"
    >
      {#snippet header()}
        <wa-checkbox
          size="small"
          disabled={!contest.usePoints}
          {@attach name("pooledPoints")}
          onchange={() => setTimeout(() => form.requestSubmit())}
          checked={contest.pooledPoints}
        ></wa-checkbox>
      {/snippet}
      {#snippet indicator()}
        {#if patchContest.isSuccess}
          <SaveIndicator />
        {/if}
      {/snippet}
    </RuleOptionCard>
  {/snippet}
</GenericForm>
