<script lang="ts">
  import CompClassForm from "@/forms/CompClassForm.svelte";
  import type { CompClassTemplate } from "@climblive/lib/models";
  import { createCompClassMutation } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import { add, roundToNearestHours } from "date-fns";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const createCompClass = createCompClassMutation(contestId);

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
  <sl-button
    size="small"
    type="submit"
    loading={$createCompClass.isPending}
    variant="primary"
    >Create
  </sl-button>
</CompClassForm>
