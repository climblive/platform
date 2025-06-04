<script lang="ts">
  import type { ScorecardSession } from "@/types";
  import { checked, GenericForm, name, value } from "@climblive/lib/forms";
  import { type ContenderPatch } from "@climblive/lib/models";
  import { getCompClassesQuery } from "@climblive/lib/queries";
  import "@shoelace-style/shoelace/dist/components/input/input.js";
  import "@shoelace-style/shoelace/dist/components/option/option.js";
  import "@shoelace-style/shoelace/dist/components/select/select.js";
  import "@shoelace-style/shoelace/dist/components/switch/switch.js";
  import { isAfter } from "date-fns";
  import { getContext, type Snippet } from "svelte";
  import type { Readable } from "svelte/store";
  import * as z from "zod/v4";

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
      <sl-input
        size="small"
        use:name={"name"}
        label="Full name"
        type="text"
        required
        use:value={data.name}
      ></sl-input>
      <sl-input
        size="small"
        use:name={"clubName"}
        label="Club name"
        type="text"
        use:value={data.clubName}
      ></sl-input>
      <sl-select
        size="small"
        use:name={"compClassId"}
        label="Competition class"
        required
        use:value={data.compClassId}
      >
        {#each $compClassesQuery.data as compClass (compClass.id)}
          <sl-option
            value={compClass.id}
            disabled={isAfter(new Date(), compClass.timeEnd)}
            >{compClass.name}</sl-option
          >
        {/each}
      </sl-select>
      <sl-switch
        size="small"
        use:name={"withdrawnFromFinals"}
        help-text="If you end up in the finals, you'll give up your spot."
        use:checked={data.withdrawnFromFinals}>Opt out of finals</sl-switch
      >
      {@render children?.()}
    </fieldset>
  </GenericForm>
{/if}

<style>
  fieldset {
    display: flex;
    flex-direction: column;
    gap: var(--sl-spacing-small);
    padding: var(--sl-spacing-medium);
  }
</style>
