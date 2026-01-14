<script lang="ts">
  import ContestForm, {
    formSchema,
    minuteInNanoseconds,
  } from "@/forms/ContestForm.svelte";
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

  const handleSubmit = (
    form: Omit<ContestTemplate, "qualifyingProblems" | "finalists">,
  ) => {
    if (createContest.isPending) {
      return;
    }

    createContest.mutate(
      {
        ...form,
        gracePeriod: form.gracePeriod * minuteInNanoseconds,
        qualifyingProblems: 0,
        finalists: 0,
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
    name: "",
    country: "SE",
    gracePeriod: 15 * minuteInNanoseconds,
    info: "",
  }}
  schema={formSchema}
>
  <div class="controls">
    <wa-button
      size="small"
      appearance="plain"
      onclick={() => navigate(`./organizers/${organizerId}/contests`)}
      >Cancel</wa-button
    >
    <wa-button
      size="small"
      type="submit"
      loading={createContest.isPending}
      variant="neutral"
      appearance="accent"
      >Create
    </wa-button>
  </div>
</ContestForm>
