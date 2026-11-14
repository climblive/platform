<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { restoreCompetitionMutation } from "@climblive/lib/queries";
  import { toastUnexpectedError } from "@climblive/lib/utils";
  import type { Snippet } from "svelte";
  import { navigate } from "svelte-routing";

  type Props = {
    competitionId: number;
    children?: Snippet<[{ restoreCompetition: () => void }]>;
  };

  let { competitionId, children }: Props = $props();

  const restoreCompetition = $derived(
    restoreCompetitionMutation(competitionId),
  );

  const handleRestore = async () => {
    try {
      const restoredCompetition =
        await restoreCompetition.mutateAsync(undefined);
      navigate(`/admin/competitions/${restoredCompetition.id}`);
    } catch {
      toastUnexpectedError("Failed to restore competition.");
    }
  };
</script>

{#if children}
  {@render children({
    restoreCompetition: handleRestore,
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
      loading={restoreCompetition.isPending}
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
