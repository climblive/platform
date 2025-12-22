<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { patchContestMutation } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";

  type Props = {
    contestId: number;
    organizerId: number;
  };

  let dialog: WaDialog | undefined = $state();

  let { contestId, organizerId }: Props = $props();

  const archiveContest = $derived(patchContestMutation(contestId));

  const handleArchive = async () => {
    if (dialog) {
      dialog.open = true;
    }
  };

  const handleCancel = () => {
    if (dialog) {
      dialog.open = false;
    }
  };

  const confirmArchivation = () => {
    archiveContest.mutate(
      { archived: true },
      {
        onSuccess: () => {
          handleCancel();
          navigate(`/admin/organizers/${organizerId}`);
        },
        onError: () => {
          toastError("Failed to archive contest.");
        },
      },
    );
  };
</script>

<div class="actions">
  <wa-button onclick={handleArchive} appearance="outlined" variant="danger"
    >Archive
    <wa-icon name="box-archive" slot="start"></wa-icon>
  </wa-button>
</div>

<wa-dialog bind:this={dialog} label="Archive contest">
  This will hide the contest for you and stop any running score engines.
  Archived contests may be permanently deleted in the future.
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
