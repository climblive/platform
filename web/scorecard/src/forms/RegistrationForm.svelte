<script lang="ts">
  import type { ScorecardSession } from "@/types";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import "@awesome.me/webawesome/dist/components/switch/switch.js";
  import { checked, GenericForm, name, value } from "@climblive/lib/forms";
  import { type ContenderPatch } from "@climblive/lib/models";
  import { getCompClassesQuery } from "@climblive/lib/queries";
  import { isAfter } from "date-fns";
  import { getContext, type Snippet } from "svelte";
  import type { Readable } from "svelte/store";
  import { z } from "@climblive/lib/utils";

  const registrationFormSchema: z.ZodType<ContenderPatch> = z.object({
    name: z.string().min(1),
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

{#if compClassesQuery.data}
  <GenericForm schema={registrationFormSchema} {submit}>
    <fieldset>
      <wa-input
        size="small"
        {@attach name("name")}
        label="Name"
        type="text"
        required
        value={data.name}
      ></wa-input>
      <wa-select
        size="small"
        {@attach name("compClassId")}
        label="Competition class"
        required
        {@attach value(data.compClassId)}
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
