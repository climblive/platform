<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { restoreContestMutation } from "@climblive/lib/queries";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import type { Snippet } from "svelte";
  import { navigate } from "svelte-routing";

  type Props = {
    contestId: number;
    children?: Snippet<[{ restoreContest: () => void }]>;
  };

  let { contestId, children }: Props = $props();

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

{#if children}
  {@render children({
    restoreContest: handleRestore,
  })}
{:else}
  <section>
    <wa-callout variant="danger" size="s">
      <wa-icon slot="icon" name="box-archive"></wa-icon>
      <p>
        <strong>This competition has been archived</strong><br />
        You can restore this competition at any time to make it active again.
      </p>
    </wa-callout>

    <wa-button
      onclick={handleRestore}
      loading={restoreContest.isPending}
      appearance="filled-outlined"
      variant="success">Restore</wa-button
    >
  </section>
{/if}

<style>
  section {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
  }
</style>
