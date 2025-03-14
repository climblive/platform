<script lang="ts">
  import CreateContestForm from "@/forms/CreateContestForm.svelte";
  import type { Contest, ContestTemplate } from "@climblive/lib/models";
  import { createContestMutation } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import "@shoelace-style/shoelace/dist/components/input/input.js";
  import { navigate } from "svelte-routing";

  interface Props {
    organizerId: number;
  }

  let { organizerId }: Props = $props();

  const createContest = createContestMutation(organizerId);

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

<CreateContestForm
  submit={handleSubmit}
  data={{
    name: "Test",
    finalists: 7,
    qualifyingProblems: 10,
    gracePeriod: 15,
  }}
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
      loading={$createContest.isPending}
      variant="primary"
      >Create
    </sl-button>
  </div>
</CreateContestForm>
