<script lang="ts">
  import CompClassForm, { formSchema } from "@/forms/CompClassForm.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import type { CompClassTemplate } from "@climblive/lib/models";
  import { createCompClassMutation } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { add, roundToNearestHours } from "date-fns";
  import { navigate } from "svelte-routing";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const createCompClass = $derived(createCompClassMutation(contestId));

  const handleSubmit = async (tmpl: CompClassTemplate) => {
    $createCompClass.mutate(tmpl, {
      onSuccess: () => navigate(`/admin/contests/${contestId}#comp-classes`),
      onError: () => toastError("Failed to create class."),
    });
  };
</script>

<CompClassForm
  submit={handleSubmit}
  data={{
    timeBegin: roundToNearestHours(add(new Date(), { hours: 1 })),
    timeEnd: roundToNearestHours(add(new Date(), { hours: 4 })),
  }}
  schema={formSchema}
>
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
      loading={$createCompClass.isPending}
      variant="neutral"
      appearance="accent"
      >Create
    </wa-button>
  </div>
</CompClassForm>

<style>
  .controls {
    display: flex;
    gap: var(--wa-space-xs);
  }
</style>
