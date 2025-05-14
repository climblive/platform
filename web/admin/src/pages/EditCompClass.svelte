<script lang="ts">
  import CompClassForm, { formSchema } from "@/forms/CompClassForm.svelte";
  import type { CompClassPatch } from "@climblive/lib/models";
  import {
    getCompClassQuery,
    patchCompClassMutation,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import { navigate } from "svelte-routing";

  interface Props {
    compClassId: number;
  }

  let { compClassId }: Props = $props();

  const compClassQuery = getCompClassQuery(compClassId);
  const patchCompClass = patchCompClassMutation(compClassId);

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
    <sl-button
      size="small"
      type="button"
      variant="text"
      onclick={history.back()}>Cancel</sl-button
    >
    <sl-button
      size="small"
      type="submit"
      loading={$patchCompClass.isPending}
      variant="primary"
      >Save
    </sl-button>
  </div>
</CompClassForm>
