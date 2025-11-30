<script lang="ts">
  import Loader from "@/components/Loader.svelte";
  import ContestForm, {
      formSchema,
      minuteInNanoseconds,
  } from "@/forms/ContestForm.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import type { Contest, ContestPatch } from "@climblive/lib/models";
  import {
      getContestQuery,
      patchContestMutation,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const contestQuery = $derived(getContestQuery(contestId));
  const patchContest = $derived(patchContestMutation(contestId));

  const contest = $derived(contestQuery.data);

  const handleSubmit = async (tmpl: ContestPatch) => {
    patchContest.mutate(
      {
        ...tmpl,
        gracePeriod:
          tmpl.gracePeriod !== undefined
            ? tmpl.gracePeriod * minuteInNanoseconds
            : undefined,
      },
      {
        onSuccess: (contest: Contest) =>
          navigate(`/admin/contests/${contest.id}`),
        onError: () => toastError("Failed to save contest."),
      },
    );
  };
</script>

{#if contest === undefined}
  <Loader />
{:else}
  <ContestForm submit={handleSubmit} data={contest} schema={formSchema}>
    <div class="controls">
      <wa-button
        size="small"
        type="button"
        appearance="plain"
        onclick={history.back()}>Cancel</wa-button
      >
      <wa-button
        size="small"
        type="submit"
        loading={patchContest.isPending}
        variant="neutral"
        >Save
      </wa-button>
    </div>
  </ContestForm>
{/if}

<style>
  .controls {
    display: flex;
    justify-content: end;
    gap: var(--wa-space-xs);
  }
</style>
