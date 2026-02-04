<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/checkbox/checkbox.js";
  import "@awesome.me/webawesome/dist/components/input/input.js";
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
      description="Points for completed problems are split by percentages. A boulder worth 1000 points with two tops will give each contender 500 points. If a third contender also tops the boulder, then all three will receive 333 points instead."
      tag="New"
    >
      {#snippet header()}
        <wa-checkbox
          size="small"
          {@attach name("pooledPoints")}
          onchange={() => setTimeout(() => form.requestSubmit())}
          {@attach checked(contest.pooledPoints)}
        ></wa-checkbox>
      {/snippet}
    </RuleOptionCard>
  {/snippet}
</GenericForm>
