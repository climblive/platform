<script lang="ts">
  import type { ScorecardSession } from "@/types";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import "@awesome.me/webawesome/dist/components/switch/switch.js";
  import { checked, GenericForm, name, value } from "@climblive/lib/forms";
  import { type CompetitorPatch } from "@climblive/lib/models";
  import { getCategoriesQuery } from "@climblive/lib/queries";
  import type { CompetitionState } from "@climblive/lib/types";
  import { z } from "@climblive/lib/utils";
  import { isAfter } from "date-fns";
  import { getContext, type Snippet } from "svelte";
  import type { Readable } from "svelte/store";

  const registrationFormSchema: z.ZodType<CompetitorPatch> = z.object({
    name: z.string().min(1),
    categoryId: z.coerce.number().gt(0, { message: "" }),
    withdrawnFromFinals: z.coerce.boolean().default(false),
  });

  interface Props {
    data: Partial<CompetitorPatch>;
    submit: (patch: CompetitorPatch) => void;
    callout?: Snippet;
    children?: Snippet;
    competitionState: CompetitionState;
    disabled?: boolean;
  }

  const { data, submit, callout, children, competitionState, ...rest }: Props =
    $props();

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");

  let categoriesQuery = $derived(getCategoriesQuery($session.competitionId));

  const disabled = $derived(rest.disabled || competitionState === "ENDED");
</script>

{#if categoriesQuery.data}
  <GenericForm schema={registrationFormSchema} {submit}>
    <fieldset>
      <wa-input
        size="s"
        {@attach name("name")}
        label="Name"
        type="text"
        required
        value={data.name}
        {disabled}
      ></wa-input>
      {@render callout?.()}
      <wa-select
        size="s"
        {@attach name("categoryId")}
        label="Category"
        required
        {@attach value(data.categoryId)}
        {disabled}
      >
        {#each categoriesQuery.data as category (category.id)}
          <wa-option
            value={category.id}
            disabled={isAfter(new Date(), category.timeEnd)}
            label={category.name}
          >
            {category.name}
            {#if category.description}
              <small>{category.description}</small>
            {/if}
          </wa-option>
        {/each}
      </wa-select>
      <wa-switch
        size="s"
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
