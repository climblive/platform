<script lang="ts">
  import ContestForm from "@/forms/ContestForm.svelte";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import type { Contest, ContestTemplate } from "@climblive/lib/models";
  import { createContestMutation } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";

  interface Props {
    organizerId: number;
  }

  let { organizerId }: Props = $props();

  const createContest = $derived(createContestMutation(organizerId));

  const handleSubmit = (form: ContestTemplate) => {
    if ($createContest.isPending) {
      return;
    }

    $createContest.mutate(
      {
        ...form,
      },
      {
        onSuccess: (contest: Contest) => navigate(`contests/${contest.id}`),
        onError: () => toastError("Failed to create contest."),
      },
    );
  };
</script>

<ContestForm
  submit={handleSubmit}
  data={{
    name: "Test",
    finalists: 7,
    qualifyingProblems: 10,
    gracePeriod: 15,
    rules: "",
  }}
>
  <div class="controls">
    <wa-button size="small" appearance="plain" onclick={history.back()}
      >Cancel</wa-button
    >
    <wa-button
      size="small"
      type="submit"
      loading={$createContest.isPending}
      variant="brand"
      appearance="accent"
      >Create
    </wa-button>
  </div>
</ContestForm>
