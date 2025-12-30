<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import type WaSelect from "@awesome.me/webawesome/dist/components/select/select.js";
  import {
    getSelfQuery,
    transferContestMutation,
  } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";

  type Props = {
    contestId: number;
    organizerId: number;
  };

  let dialog: WaDialog | undefined = $state();
  let selectedOrganizerID: number | undefined = $state();

  const { contestId, organizerId }: Props = $props();

  const selfQuery = $derived(getSelfQuery());
  const transferContest = $derived(transferContestMutation(contestId));

  const organizers = $derived(selfQuery.data?.organizers ?? []);
  const otherOrganizers = $derived(
    organizers.filter(({ id }) => id !== organizerId),
  );

  const handleTransfer = async () => {
    if (dialog) {
      dialog.open = true;
      selectedOrganizerID = undefined;
    }
  };

  const handleCancel = () => {
    if (dialog) {
      dialog.open = false;
    }
  };

  const confirmTransfer = () => {
    if (selectedOrganizerID === undefined) {
      return;
    }

    transferContest.mutate(selectedOrganizerID, {
      onSuccess: () => {
        handleCancel();

        navigate(`./organizers/${selectedOrganizerID}/contests`);
      },
      onError: (error) => {
        toastError("Failed to transfer contest.");
      },
    });
  };

  const handleSelect = (event: Event) => {
    const select = event.target as WaSelect;
    selectedOrganizerID = Number(select.value);
  };
</script>

<div class="actions">
  <wa-button
    onclick={handleTransfer}
    appearance="outlined"
    disabled={otherOrganizers.length === 0}
  >
    Transfer
    <wa-icon name="arrow-right" slot="start"></wa-icon>
  </wa-button>
</div>

<wa-dialog bind:this={dialog} label="Transfer contest">
  Move this contest to one of the other organizers you belong to.
  <wa-select label="Select new organizer" onchange={handleSelect}>
    {#each otherOrganizers as organizer (organizer.id)}
      <wa-option value={organizer.id}>{organizer.name}</wa-option>
    {/each}
  </wa-select>

  <wa-button slot="footer" appearance="plain" onclick={handleCancel}>
    Cancel
  </wa-button>
  <wa-button
    slot="footer"
    variant="warning"
    onclick={confirmTransfer}
    loading={transferContest.isPending}
    disabled={selectedOrganizerID === undefined || otherOrganizers.length === 0}
  >
    Transfer
    <wa-icon slot="start" name="arrow-right"></wa-icon>
  </wa-button>
</wa-dialog>

<style>
  wa-select {
    margin-block-start: var(--wa-space-m);
  }
</style>
