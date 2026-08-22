<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { restoreContestMutation } from "@climblive/lib/queries";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import { navigate } from "svelte-routing";

  type Props = {
    contestId: number;
  };

  let { contestId }: Props = $props();

  const restoreContest = $derived(restoreContestMutation(contestId));

  const handleRestore = () => {
    restoreContest.mutate(undefined, {
      onSuccess: () => {
        navigate(`/admin/contests/${contestId}`);
      },
      onError: () => {
        toastUnexpectedError("Failed to restore competition.");
      },
    });
  };
</script>

<section>
  <wa-callout variant="danger" size="s">
    <wa-icon slot="icon" name="box-archive"></wa-icon>
    <p>
      <strong>This competition has been archived.</strong><br />
      You can restore this competition at any time to make it active again.
    </p>
  </wa-callout>

  <wa-button
    onclick={handleRestore}
    appearance="filled-outlined"
    variant="success">Restore</wa-button
  >
</section>

<style>
  section {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
  }
</style>
