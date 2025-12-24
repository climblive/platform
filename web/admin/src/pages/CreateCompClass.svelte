<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import CompClassForm, { formSchema } from "@/forms/CompClassForm.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import type { CompClassTemplate } from "@climblive/lib/models";
  import {
    createCompClassMutation,
    getCompClassesQuery,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { add, roundToNearestHours } from "date-fns";
  import { navigate } from "svelte-routing";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const compClassesQuery = $derived(getCompClassesQuery(contestId));
  const createCompClass = $derived(createCompClassMutation(contestId));

  const compClasses = $derived(compClassesQuery.data);

  const initialTimes = $derived.by(() => {
    if (!compClasses || compClasses.length === 0) {
      return {
        timeBegin: roundToNearestHours(add(new Date(), { hours: 1 })),
        timeEnd: roundToNearestHours(add(new Date(), { hours: 4 })),
      };
    }

    const lastCompClass = compClasses[compClasses.length - 1];

    return {
      timeBegin: lastCompClass.timeBegin,
      timeEnd: lastCompClass.timeEnd,
    };
  });

  const handleSubmit = async (tmpl: CompClassTemplate) => {
    createCompClass.mutate(tmpl, {
      onSuccess: () => navigate(`/admin/contests/${contestId}#comp-classes`),
      onError: () => toastError("Failed to create class."),
    });
  };
</script>

{#if compClasses === undefined}
  <Loader />
{:else}
  <CompClassForm submit={handleSubmit} data={initialTimes} schema={formSchema}>
    <div class="controls">
      <wa-button
        size="small"
        type="button"
        appearance="plain"
        onclick={() => navigate(`/admin/contests/${contestId}#comp-classes`)}
        >Cancel</wa-button
      >
      <wa-button
        size="small"
        type="submit"
        loading={createCompClass.isPending}
        variant="neutral"
        appearance="accent"
        >Create
      </wa-button>
    </div>
  </CompClassForm>
{/if}

<style>
  .controls {
    display: flex;
    gap: var(--wa-space-xs);
    justify-content: end;
  }
</style>
