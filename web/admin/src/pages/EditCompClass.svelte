<script lang="ts">
  import CompClassForm, { formSchema } from "@/forms/CompClassForm.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import type { CompClassPatch } from "@climblive/lib/models";
  import {
    getCompClassQuery,
    patchCompClassMutation,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";

  interface Props {
    compClassId: number;
  }

  let { compClassId }: Props = $props();

  const compClassQuery = $derived(getCompClassQuery(compClassId));
  const patchCompClass = $derived(patchCompClassMutation(compClassId));

  const compClass = $derived($compClassQuery.data);

  const handleSubmit = async (patch: CompClassPatch) => {
    $patchCompClass.mutate(patch, {
      onSuccess: (compClass) =>
        navigate(`/admin/contests/${compClass.contestId}#comp-classes`),
      onError: () => toastError("Failed to save comp class."),
    });
  };
</script>

<CompClassForm
  submit={handleSubmit}
  data={{
    ...compClass,
  }}
  schema={formSchema}
>
  <div class="controls">
    <wa-button
      size="small"
      type="button"
      variant="text"
      onclick={history.back()}>Cancel</wa-button
    >
    <wa-button
      size="small"
      type="submit"
      loading={$patchCompClass.isPending}
      variant="brand"
      >Save
    </wa-button>
  </div>
</CompClassForm>
