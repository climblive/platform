<script lang="ts">
  import RegistrationForm from "@/forms/RegistrationForm.svelte";
  import type { ScorecardSession } from "@/types";
  import type { RegistrationFormData } from "@climblive/lib/models";
  import {
    getContenderQuery,
    updateContenderMutation,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import "@shoelace-style/shoelace/dist/components/alert/alert.js";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import "@shoelace-style/shoelace/dist/components/icon/icon.js";
  import { getContext } from "svelte";
  import { navigate } from "svelte-routing";
  import type { Readable } from "svelte/store";
  import Loading from "./Loading.svelte";

  const session = getContext<Readable<ScorecardSession>>("scorecardSession");

  const contenderQuery = getContenderQuery($session.contenderId);
  const updateContender = updateContenderMutation($session.contenderId);

  $: contender = $contenderQuery.data;

  const gotoScorecard = () => {
    navigate(`/${contender?.registrationCode}`);
  };

  const handleSubmit = ({
    detail: form,
  }: CustomEvent<RegistrationFormData>) => {
    if (!contender || $updateContender.isPending) {
      return;
    }

    $updateContender.mutate(
      {
        ...contender,
        ...form,
        publicName: form.name,
      },
      {
        onSuccess: gotoScorecard,
        onError: () => toastError("Registration was not successful."),
      },
    );
  };
</script>

{#if !contender}
  <Loading />
{:else}
  <RegistrationForm
    on:submit={handleSubmit}
    data={{
      name: contender.name,
      clubName: contender.clubName,
      compClassId: contender.compClassId,
    }}
  >
    <sl-button
      size="small"
      type="submit"
      loading={$updateContender.isPending}
      variant="primary"
      >Register
    </sl-button>
  </RegistrationForm>
{/if}
