<script lang="ts">
  import type { ScorecardSession } from "@/types";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import "@awesome.me/webawesome/dist/components/switch/switch.js";
  import { checked, GenericForm, name, value } from "@climblive/lib/forms";
  import { type ContenderPatch } from "@climblive/lib/models";
  import { getCompClassesQuery } from "@climblive/lib/queries";
  import type { ContestState } from "@climblive/lib/types";
  import { z } from "@climblive/lib/utils";
  import { add, formatDistance, isAfter } from "date-fns";
  import { getContext, type Snippet } from "svelte";
  import type { Readable } from "svelte/store";

  export const nanosecondsInMinute = 60 * 1_000_000_000;

  const registrationFormSchema: z.ZodType<ContenderPatch> = z.object({
    name: z.string().min(1),
    compClassId: z.coerce.number().gt(0, { message: "" }),
    withdrawnFromFinals: z.coerce.boolean(),
  });

  interface Props {
    data: Partial<ContenderPatch>;
    nameRetentionTime: number;
    submit: (patch: ContenderPatch) => void;
    children?: Snippet;
    contestState: ContestState;
  }

  let { data, nameRetentionTime, submit, children, contestState }: Props =
    $props();

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");

  let compClassesQuery = $derived(getCompClassesQuery($session.contestId));

  const retentionDuration = $derived.by(() => {
    const base = new Date(0);
    return formatDistance(
      add(base, {
        minutes: nameRetentionTime / nanosecondsInMinute,
      }),
      base,
    );
  });
</script>

{#if compClassesQuery.data}
  {@const disabled = contestState === "ENDED"}

  <GenericForm schema={registrationFormSchema} {submit}>
    <fieldset>
      <wa-input
        size="small"
        {@attach name("name")}
        label="Name"
        type="text"
        required
        value={data.name}
        {disabled}
      ></wa-input>
      <wa-callout variant="neutral" size="small" open>
        <wa-icon slot="icon" name="circle-info"></wa-icon>
        Your name will be removed and your results anonymized
        {retentionDuration} after the contest ends.
      </wa-callout>
      <wa-select
        size="small"
        {@attach name("compClassId")}
        label="Competition class"
        required
        {@attach value(data.compClassId)}
        {disabled}
      >
        {#each compClassesQuery.data as compClass (compClass.id)}
          <wa-option
            value={compClass.id}
            disabled={isAfter(new Date(), compClass.timeEnd)}
            label={compClass.name}
          >
            {compClass.name}
            {#if compClass.description}
              <small>{compClass.description}</small>
            {/if}
          </wa-option>
        {/each}
      </wa-select>
      <wa-switch
        size="small"
        {@attach name("withdrawnFromFinals")}
        hint="If you do not wish to participate in the finals, you can give up your spot."
        {disabled}
        {@attach checked(data.withdrawnFromFinals)}>Opt out of finals</wa-switch
      >
      {@render children?.()}
    </fieldset>
  </GenericForm>
{/if}

<style>
  fieldset {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
    padding: var(--wa-space-m);
  }
</style>
