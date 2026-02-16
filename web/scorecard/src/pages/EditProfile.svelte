<script lang="ts">
  import RegistrationForm from "@/forms/RegistrationForm.svelte";
  import type { ScorecardSession } from "@/types";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import type { ContenderPatch } from "@climblive/lib/models";
  import {
    getContenderQuery,
    getContestQuery,
    patchContenderMutation,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { getContext } from "svelte";
  import { navigate } from "svelte-routing";
  import type { Readable } from "svelte/store";
  import Loading from "./Loading.svelte";

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");

  const contenderQuery = $derived(getContenderQuery($session.contenderId));
  const contestQuery = $derived(getContestQuery($session.contestId));
  const patchContender = $derived(patchContenderMutation($session.contenderId));

  let contender = $derived(contenderQuery.data);
  let contest = $derived(contestQuery.data);

  const gotoScorecard = () => {
    navigate(`/${contender?.registrationCode}`);
  };

  const handleSubmit = (form: ContenderPatch) => {
    if (!contender || patchContender.isPending) {
      return;
    }

    patchContender.mutate(
      {
        ...form,
      },
      {
        onSuccess: gotoScorecard,
        onError: () => toastError("Failed to save registration data."),
      },
    );
  };
</script>

{#if !contender || !contest}
  <Loading />
{:else}
  <RegistrationForm
    submit={handleSubmit}
    nameRetentionTime={contest.nameRetentionTime}
    data={{
      name: contender.name,
      compClassId: contender.compClassId,
      withdrawnFromFinals: contender.withdrawnFromFinals,
    }}
  >
    <div class="controls">
      <wa-button
        size="small"
        type="button"
        onclick={gotoScorecard}
        appearance="plain"
        >Cancel
      </wa-button>
      <wa-button
        size="small"
        type="submit"
        loading={patchContender.isPending}
        disabled={false}
        variant="neutral"
        appearance="accent"
        >Save
      </wa-button>
    </div>
  </RegistrationForm>
{/if}

<style>
  .controls {
    display: flex;
    justify-content: end;
    gap: var(--wa-space-xs);
  }
</style>
