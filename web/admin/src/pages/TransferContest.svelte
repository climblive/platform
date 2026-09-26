<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import type WaSelect from "@awesome.me/webawesome/dist/components/select/select.js";
  import { value } from "@climblive/lib/forms";
  import {
    getSelfQuery,
    transferContestMutation,
  } from "@climblive/lib/queries";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import type { Snippet } from "svelte";
  import { navigate } from "svelte-routing";

  type Props = {
    contestId: number;
    contestName: string;
    open?: boolean;
    onClose?: () => void;
    children?: Snippet<[{ transferContest: () => void; disabled: boolean }]>;
    organizerId: number;
  };

  let selectedOrganizerId: number | undefined = $state();

  const {
    contestId,
    contestName,
    organizerId,
    children,
    open: initialOpen = false,
    onClose,
  }: Props = $props();

  let open = $derived(initialOpen);

  const selfQuery = $derived(getSelfQuery());
  const transferContest = $derived(transferContestMutation(contestId));

  const organizers = $derived(selfQuery.data?.organizers ?? []);
  const otherOrganizers = $derived(
    organizers.filter(({ id }) => id !== organizerId),
  );

  const handleTransfer = () => {
    open = true;
  };

  const handleCancel = () => {
    open = false;
  };

  const confirmTransfer = async () => {
    if (selectedOrganizerId === undefined) {
      return;
    }

    try {
      const transferredContest =
        await transferContest.mutateAsync(selectedOrganizerId);
      handleCancel();
      navigate(
        `/admin/organizers/${transferredContest.ownership.organizerId}/contests`,
      );
    } catch {
      toastUnexpectedError("Failed to transfer competition.");
    }
  };

  const handleSelect = (event: Event) => {
    const select = event.target as WaSelect;
    selectedOrganizerId = Number(select.value);
  };
</script>

{#if children}
  {@render children({
    transferContest: handleTransfer,
    disabled: otherOrganizers.length === 0,
  })}
{:else}
  <wa-button
    onclick={handleTransfer}
    appearance="outlined"
    disabled={otherOrganizers.length === 0}
  >
    Transfer
    <wa-icon name="arrow-right" slot="start"></wa-icon>
  </wa-button>
{/if}

<wa-dialog
  label="Transfer competition"
  {open}
  onwa-after-hide={() => {
    open = false;
    onClose?.();
  }}
>
  <wa-select
    label="Select new organizer"
    onchange={handleSelect}
    {@attach value(selectedOrganizerId)}
    hint="Select one of the other organizers you belong to."
  >
    {#each otherOrganizers as organizer (organizer.id)}
      <wa-option value={organizer.id}>{organizer.name}</wa-option>
    {/each}
  </wa-select>

  {#if selectedOrganizerId}
    {@const currentOrganizer = organizers.find(
      (organizer) => organizer.id === organizerId,
    )}
    {@const newOrganizer = organizers.find(
      (organizer) => organizer.id === selectedOrganizerId,
    )}

    {#if currentOrganizer && newOrganizer}
      <wa-callout variant="warning">
        <wa-icon slot="icon" name="triangle-exclamation"></wa-icon>
        This will transfer all competition data of
        <strong>{contestName}</strong>
        from the current organizer
        <strong>
          {currentOrganizer.name}
        </strong>
        to
        <strong>
          {newOrganizer.name}
        </strong>.
      </wa-callout>
    {/if}
  {/if}

  <wa-button slot="footer" appearance="plain" onclick={handleCancel}>
    Cancel
  </wa-button>
  <wa-button
    slot="footer"
    variant="warning"
    onclick={confirmTransfer}
    loading={transferContest.isPending}
    disabled={selectedOrganizerId === undefined}
  >
    Transfer
    <wa-icon slot="start" name="arrow-right"></wa-icon>
  </wa-button>
</wa-dialog>

<style>
  wa-dialog {
    white-space: normal;
  }

  wa-dialog::part(body) {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
  }
</style>
