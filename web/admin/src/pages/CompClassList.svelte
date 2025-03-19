<script lang="ts">
  import CompClassForm from "@/forms/CompClassForm.svelte";
  import type { CompClassTemplate } from "@climblive/lib/models";
  import {
    createCompClassMutation,
    getCompClassesQuery,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import "@shoelace-style/shoelace/dist/components/icon/icon.js";
  import "@shoelace-style/shoelace/dist/components/qr-code/qr-code.js";
  import { add, roundToNearestHours } from "date-fns";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const compClassesQuery = getCompClassesQuery(contestId);
  const createCompClass = createCompClassMutation(contestId);

  let compClasses = $derived($compClassesQuery.data);

  const handleSubmit = async (tmpl: CompClassTemplate) => {
    $createCompClass.mutate(tmpl, {
      onError: () => toastError("Failed to create class."),
    });
  };
</script>

<CompClassForm
  submit={handleSubmit}
  data={{
    name: "Males",
    timeBegin: roundToNearestHours(add(new Date(), { hours: 1 })),
    timeEnd: roundToNearestHours(add(new Date(), { hours: 4 })),
  }}
>
  <div class="controls">
    <sl-button size="small" type="button" variant="text">
      <sl-button
        size="small"
        type="submit"
        loading={$createCompClass.isPending}
        variant="primary"
        >Create
      </sl-button>
    </sl-button>
  </div>
</CompClassForm>

<style>
  section {
    display: flex;
    gap: var(--sl-spacing-x-small);
  }
</style>
