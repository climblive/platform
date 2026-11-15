<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import CategoryForm, { formSchema } from "@/forms/CategoryForm.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import type { CategoryTemplate } from "@climblive/lib/models";
  import {
    createCategoryMutation,
    getCategoriesQuery,
  } from "@climblive/lib/queries";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import { add, roundToNearestHours } from "date-fns";
  import { navigate } from "svelte-routing";

  interface Props {
    competitionId: string;
  }

  let { competitionId }: Props = $props();

  const categoriesQuery = $derived(getCategoriesQuery(competitionId));
  const createCategory = $derived(createCategoryMutation(competitionId));

  const categories = $derived(categoriesQuery.data);

  const initialTimes = $derived.by(() => {
    if (!categories || categories.length === 0) {
      return {
        timeBegin: roundToNearestHours(add(new Date(), { hours: 1 })),
        timeEnd: roundToNearestHours(add(new Date(), { hours: 4 })),
      };
    }

    const lastCategory = categories[categories.length - 1];

    return {
      timeBegin: lastCategory.timeBegin,
      timeEnd: lastCategory.timeEnd,
    };
  });

  const handleSubmit = async (tmpl: CategoryTemplate) => {
    createCategory.mutate(tmpl, {
      onSuccess: () =>
        navigate(`/admin/competitions/${competitionId}#categories`),
      onError: () => toastUnexpectedError("Failed to create category."),
    });
  };
</script>

{#if categories === undefined}
  <Loader />
{:else}
  <CategoryForm submit={handleSubmit} data={initialTimes} schema={formSchema}>
    <div class="controls">
      <wa-button
        size="s"
        type="button"
        appearance="plain"
        onclick={() =>
          navigate(`/admin/competitions/${competitionId}#categories`)}
        >Cancel</wa-button
      >
      <wa-button
        size="s"
        type="submit"
        loading={createCategory.isPending}
        variant="neutral"
        appearance="accent"
        >Create
      </wa-button>
    </div>
  </CategoryForm>
{/if}

<style>
  .controls {
    display: flex;
    gap: var(--wa-space-xs);
    justify-content: end;
  }
</style>
