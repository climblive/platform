<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import CategoryForm, { formSchema } from "@/forms/CategoryForm.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import type { CategoryPatch } from "@climblive/lib/models";
  import {
    getCategoryQuery,
    patchCategoryMutation,
  } from "@climblive/lib/queries";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";

  interface Props {
    categoryId: number;
  }

  let { categoryId }: Props = $props();

  const categoryQuery = $derived(getCategoryQuery(categoryId));
  const patchCategory = $derived(patchCategoryMutation(categoryId));

  const category = $derived(categoryQuery.data);

  const handleSubmit = async (patch: CategoryPatch) => {
    patchCategory.mutate(patch, {
      onSuccess: (category) =>
        navigate(`/admin/competitions/${category.competitionId}#categories`),
      onError: () => toastUnexpectedError("Failed to save category."),
    });
  };
</script>

{#if category === undefined}
  <Loader />
{:else}
  <CategoryForm submit={handleSubmit} data={category} schema={formSchema}>
    <div class="controls">
      <wa-button
        size="s"
        type="button"
        appearance="plain"
        onclick={() =>
          navigate(`/admin/competitions/${category.competitionId}#categories`)}
        >Cancel</wa-button
      >
      <wa-button
        size="s"
        type="submit"
        loading={patchCategory.isPending}
        variant="neutral"
        >Save
      </wa-button>
    </div>
  </CategoryForm>
{/if}

<style>
  .controls {
    display: flex;
    gap: var(--wa-space-xs);
    justify-content: right;
  }
</style>
