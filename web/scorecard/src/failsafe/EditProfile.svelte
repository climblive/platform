<script lang="ts">
  import {
    getCompClassesQuery,
    getContenderQuery,
    getContestQuery,
    patchContenderMutation,
    scrubContenderMutation,
  } from "@climblive/lib/queries";
  import { SyncedTime, toastUnexpectedError } from "@climblive/lib/utils";
  import { add, formatDistance, isBefore } from "date-fns";
  import { onMount } from "svelte";

  type Props = {
    contestId: number;
    contenderId: number;
  };

  const { contestId, contenderId }: Props = $props();

  let form: HTMLFormElement | undefined = $state();
  let showInfo = $state(false);

  const contenderQuery = $derived(getContenderQuery(contenderId));
  const contestQuery = $derived(getContestQuery(contestId));
  const compClassesQuery = $derived(getCompClassesQuery(contestId));
  const patchContender = $derived(patchContenderMutation(contenderId));
  const scrubContender = $derived(scrubContenderMutation(contenderId));
  const time = new SyncedTime(60000);
  onMount(() => {
    time.start();
    return () => time.stop();
  });

  let contender = $derived(contenderQuery.data);
  let contest = $derived(contestQuery.data);
  let compClasses = $derived(compClassesQuery.data);
  let selectedCompClass = $derived(
    compClasses?.find(({ id }) => id === contender?.compClassId),
  );

  const retentionDuration = $derived.by(() =>
    contender?.scrubBefore && !isBefore(contender.scrubBefore, time.current)
      ? formatDistance(contender.scrubBefore, time.current)
      : undefined,
  );
  const registrationRetentionDuration = $derived.by(() => {
    const base = new Date(0);
    return formatDistance(
      add(base, { minutes: (contest?.nameRetentionTime ?? 0) / 60000000000 }),
      base,
    );
  });

  const handleSubmit = (event: SubmitEvent) => {
    event.preventDefault();

    if (!form || !contender) {
      return;
    }

    const formData = new FormData(form);
    const name = formData.get("name")?.toString().trim();
    const compClassId = formData.get("compClassId")?.toString().trim();

    if (name && compClassId) {
      patchContender.mutate({
        ...contender,
        name,
        compClassId: Number(compClassId),
      });
    }
  };

  const handleScrub = () => {
    if (
      !confirm(
        "Your name will be permanently removed and your results will be anonymized. This action cannot be undone.\n\nBe aware that without a name, you will lose your chance at finals and you cannot take part in any prize raffles.",
      )
    ) {
      return;
    }
    scrubContender.mutate(undefined, {
      onSuccess: () => window.location.reload(),
      onError: () => toastUnexpectedError("Failed to remove your name."),
    });
  };
</script>

{#if compClasses && contender}
  <form onsubmit={handleSubmit} bind:this={form}>
    <div class="name-field">
      <div class="name-row">
        <input
          required
          placeholder="Name"
          name="name"
          type="text"
          value={contender.name}
          aria-label="Name"
        />
        {#if contender.entered && !contender.scrubbedAt}
          <button
            class="info-button"
            type="button"
            aria-expanded={showInfo}
            onclick={() => (showInfo = !showInfo)}
          >
            Info
          </button>
        {/if}
      </div>
      {#if !contender.entered && contest}
        <p class="info" role="note">
          Your name will be stored for {registrationRetentionDuration} after the competition
          ends, after which it will be removed and your results anonymized.
        </p>
      {:else if showInfo && !contender.scrubbedAt}
        <p class="info" role="note">
          {#if retentionDuration}
            Your name will be kept stored for {retentionDuration} from now, after
            which it will be removed and your results anonymized.
          {:else}
            Your name will be removed and your results anonymized shortly.
          {/if}
        </p>
      {/if}
    </div>
    <select
      name="compClassId"
      required
      value={selectedCompClass?.id}
      aria-label="Category"
    >
      {#each compClasses as compClass (compClass.id)}
        <option value={compClass.id}>{compClass.name}</option>
      {/each}
    </select>
    <div class="actions">
      {#if contender.name}<button
          type="button"
          disabled={scrubContender.isPending}
          onclick={handleScrub}>Remove my name</button
        >{/if}
      <button
        type="submit"
        disabled={patchContender.isPending || scrubContender.isPending}
        >{contender.entered ? "Update" : "Register"}</button
      >
    </div>
  </form>
{/if}

<style>
  .info-button {
    padding: 0.15rem 0.35rem;
    font-size: 0.75rem;
    white-space: nowrap;
  }

  .info {
    padding: var(--wa-space-s);
    border: 1px solid;
    border-radius: 0.25rem;
  }

  form {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
  }

  .name-field {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-xs);
  }

  .name-row,
  .actions {
    display: flex;
    gap: var(--wa-space-xs);
  }

  .name-row input {
    flex: 1;
    min-width: 0;
  }

  .actions {
    flex-wrap: nowrap;
  }
</style>
