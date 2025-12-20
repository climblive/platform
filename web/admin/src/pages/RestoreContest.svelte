<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { patchContestMutation } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";

  type Props = {
    contestId: number;
  };

  let { contestId }: Props = $props();

  const patchContest = $derived(patchContestMutation(contestId));

  const handleRestore = () => {
    patchContest.mutate(
      { archived: false },
      {
        onSuccess: () => {
          navigate(`/admin/contests/${contestId}`);
        },
        onError: () => {
          toastError("Failed to restore contest.");
        },
      },
    );
  };
</script>

<wa-callout variant="danger" size="small">
  <wa-icon slot="icon" name="box-archive"></wa-icon>
  <p>
    <strong>This contest has been archived.</strong><br />
    You can restore this contest at any time to make it active again.
  </p>
  <wa-button
    onclick={handleRestore}
    appearance="filled-outlined"
    variant="success">Restore</wa-button
  >
</wa-callout>

<style>
  wa-callout::part(message) {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--wa-space-m);
  }
</style>
