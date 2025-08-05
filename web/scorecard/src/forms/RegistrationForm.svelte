<script lang="ts">
  import type { ScorecardSession } from "@/types";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import "@awesome.me/webawesome/dist/components/switch/switch.js";
  import { GenericForm, name } from "@climblive/lib/forms";
  import { type ContenderPatch } from "@climblive/lib/models";
  import { getCompClassesQuery } from "@climblive/lib/queries";
  import { isAfter } from "date-fns";
  import { getContext, type Snippet } from "svelte";
  import type { Readable } from "svelte/store";
  import * as z from "zod";

  const registrationFormSchema: z.ZodType<ContenderPatch> = z.object({
    name: z.string().min(1),
    clubName: z.string().optional(),
    compClassId: z.coerce.number().gt(0, { message: "" }),
    withdrawnFromFinals: z.coerce.boolean(),
  });

  interface Props {
    data: Partial<ContenderPatch>;
    submit: (patch: ContenderPatch) => void;
    children?: Snippet;
  }

  let { data, submit, children }: Props = $props();

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");

  let compClassesQuery = $derived(getCompClassesQuery($session.contestId));
</script>

{#if $compClassesQuery.data}
  <GenericForm schema={registrationFormSchema} {submit}>
    <fieldset>
      <wa-input
        size="small"
        {@attach name("name")}
        label="Full name"
        type="text"
        required
        value={data.name}
      ></wa-input>
      <wa-input
        size="small"
        {@attach name("clubName")}
        label="Club name"
        type="text"
        value={data.clubName}
      ></wa-input>
      <wa-select
        size="small"
        {@attach name("compClassId")}
        label="Competition class"
        required
        value={data.compClassId}
      >
        {#each $compClassesQuery.data as compClass (compClass.id)}
          <wa-option
            value={compClass.id}
            disabled={isAfter(new Date(), compClass.timeEnd)}
            >{compClass.name}</wa-option
          >
        {/each}
      </wa-select>
      <wa-switch
        size="small"
        {@attach name("withdrawnFromFinals")}
        hint="If you end up in the finals, you'll give up your spot."
        checked={data.withdrawnFromFinals}>Opt out of finals</wa-switch
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
