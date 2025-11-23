<script lang="ts">
  import RegistrationForm from "@/forms/RegistrationForm.svelte";
  import type { ScorecardSession } from "@/types";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
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

  const contender = $derived(contenderQuery.data);
  const contest = $derived(contestQuery.data);

  const gotoScorecard = () => {
    navigate(`/${contender?.registrationCode}`, { replace: true });
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
        onError: () => toastError("Registration was not successful."),
      },
    );
  };
</script>

{#if !contender || !contest}
  <Loading />
{:else}
  <h1>{contest.name}</h1>
  <RegistrationForm
    submit={handleSubmit}
    data={{
      name: contender.name,
      clubName: contender.clubName,
      compClassId: contender.compClassId,
      withdrawnFromFinals: contender.withdrawnFromFinals,
    }}
  >
    <wa-button
      size="small"
      type="submit"
      loading={patchContender.isPending}
      variant="neutral"
      appearance="accent"
      >Register
    </wa-button>
  </RegistrationForm>
{/if}

<style>
  h1 {
    font-size: var(--wa-font-size-l);
    padding: var(--wa-space-m);
    padding-block-end: 0;
  }
</style>
