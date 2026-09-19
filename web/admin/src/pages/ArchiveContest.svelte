<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import {
    archiveContestMutation,
    getContestQuery,
  } from "@climblive/lib/queries";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import type { Snippet } from "svelte";
  import { navigate } from "svelte-routing";

  type Props = {
    contestId: number;
    children?: Snippet<[{ archiveContest: () => void }]>;
    organizerId: number;
  };

  let dialog: WaDialog | undefined = $state();

  let { contestId, organizerId, children }: Props = $props();

  const contestQuery = $derived(getContestQuery(contestId));
  const contest = $derived(contestQuery.data);

  const archiveContest = $derived(archiveContestMutation(contestId));

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
    archiveContest.mutate(undefined, {
      onSuccess: () => {
        handleCancel();
        navigate(`/admin/organizers/${organizerId}/contests`);
      },
      onError: () => {
        toastUnexpectedError("Failed to archive competition.");
      },
    });
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

<wa-dialog bind:this={dialog} label="Archive competition">
  This will hide the competition <strong>{contest?.name}</strong> for you and
  stop any running score engines.<br /><br />
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
