<script lang="ts">
  import RegistrationForm from "@/forms/RegistrationForm.svelte";
  import type { ScorecardSession } from "@/types";
  import type { RegistrationFormData } from "@climblive/lib/models";
  import {
    getContenderQuery,
    updateContenderMutation,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
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
        onError: () => toastError("Failed to save registration data."),
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
    <div class="controls">
      <sl-button
        size="small"
        type="button"
        variant="text"
        on:click={gotoScorecard}
        >Cancel
      </sl-button>
      <sl-button
        size="small"
        type="submit"
        loading={$updateContender.isPending}
        disabled={false}
        variant="primary"
        >Save
      </sl-button>
    </div>
  </RegistrationForm>
{/if}

<style>
  .controls {
    display: flex;
    justify-content: end;
    gap: var(--sl-spacing-small);
  }
</style>
