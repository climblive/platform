<script lang="ts">
  import ProblemForm, { formSchema } from "@/forms/ProblemForm.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import type { ProblemTemplate } from "@climblive/lib/models";
  import { createProblemMutation } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const createProblem = $derived(createProblemMutation(contestId));

  const handleSubmit = async (tmpl: ProblemTemplate) => {
    $createProblem.mutate(tmpl, {
      onSuccess: () => navigate(`/admin/contests/${contestId}#problems`),
      onError: () => toastError("Failed to create problem."),
    });
  };
</script>

<ProblemForm
  submit={handleSubmit}
  data={{
    number: 1,
    holdColorPrimary: "#000000",
    pointsTop: 100,
    pointsZone: 0,
    flashBonus: 0,
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
      loading={$createProblem.isPending}
      variant="brand"
      >Create
    </wa-button>
  </div>
</ProblemForm>
