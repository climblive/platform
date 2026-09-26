<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { archiveContestMutation } from "@climblive/lib/queries";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import type { Snippet } from "svelte";
  import { navigate } from "svelte-routing";

  type Props = {
    contestId: number;
    contestName: string;
    open?: boolean;
    onClose?: () => void;
    children?: Snippet<[{ archiveContest: () => void }]>;
    organizerId: number;
  };

  let {
    contestId,
    contestName,
    organizerId,
    children,
    open: initialOpen = false,
    onClose,
  }: Props = $props();

  let open = $derived(initialOpen);

  const archiveContest = $derived(archiveContestMutation(contestId));

  const handleArchive = () => {
    open = true;
  };

  const handleCancel = () => {
    open = false;
  };

  const confirmArchivation = async () => {
    const destination = `/admin/organizers/${organizerId}/contests`;

    try {
      await archiveContest.mutateAsync(undefined);
      handleCancel();
      navigate(destination);
    } catch {
      toastUnexpectedError("Failed to archive competition.");
    }
  };
</script>

{#if children}
  {@render children({ archiveContest: handleArchive })}
{:else}
  <div class="actions">
    <wa-button onclick={handleArchive} appearance="outlined" variant="danger"
      >Archive
      <wa-icon name="box-archive" slot="start"></wa-icon>
    </wa-button>
  </div>
{/if}

<wa-dialog
  label="Archive competition"
  {open}
  onwa-after-hide={() => {
    open = false;
    onClose?.();
  }}
>
  This will hide the competition <strong>{contestName}</strong> for you and stop
  any running score engines.<br /><br />
  Archived competitions may be permanently deleted in the future.
  <wa-button slot="footer" appearance="plain" onclick={handleCancel}>
    Cancel</wa-button
  >
  <wa-button
    slot="footer"
    variant="danger"
    onclick={confirmArchivation}
    loading={archiveContest.isPending}
  >
    Archive
    <wa-icon slot="start" name="box-archive"></wa-icon>
  </wa-button>
</wa-dialog>

<style>
  wa-dialog {
    white-space: normal;
  }
</style>
