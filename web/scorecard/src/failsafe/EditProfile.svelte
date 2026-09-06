<script lang="ts">
  import {
    getCompClassesQuery,
    getContenderQuery,
    getContestQuery,
    patchContenderMutation,
    scrubContenderMutation,
  } from "@climblive/lib/queries";
  import { maskScrubbedName, SyncedTime } from "@climblive/lib/utils";
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
  const time = new SyncedTime(60_000);
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
      add(base, {
        minutes: (contest?.nameRetentionTime ?? 0) / 60_000_000_000,
      }),
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
      patchContender.mutate(
        {
          name,
          compClassId: Number(compClassId),
        },
        {
          onError: () => window.alert("Failed to save registration data."),
        },
      );
    }
  };

  const handleScrub = () => {
    if (
      !window.confirm(
        "Your name will be permanently removed and your results will be anonymized. This action cannot be undone.\n\nBe aware that without a name, you will lose your chance at finals and you cannot take part in any prize raffles.",
      )
    ) {
      return;
    }

    scrubContender.mutate(undefined, {
      onSuccess: () => window.location.reload(),
      onError: () => window.alert("Failed to remove your name."),
    });
  };
</script>

{#if compClasses && contender}
  <form onsubmit={handleSubmit} bind:this={form}>
    <div class="name-row">
      <input
        required
        placeholder="Name"
        name="name"
        type="text"
        value={contender.scrubbedAt !== undefined
          ? maskScrubbedName(contender.id)
          : contender.name}
        aria-label="Name"
        disabled={contender.scrubbedAt !== undefined}
      />
      {#if contender.entered && !contender.scrubbedAt}
        <button type="button" onclick={() => (showInfo = !showInfo)}>
          Info
        </button>
      {/if}
    </div>
    {#if !contender.entered && contest}
      <p class="info">
        Your name will be stored for {registrationRetentionDuration} after the competition
        ends, after which it will be removed and your results anonymized.
      </p>
    {:else if showInfo}
      <p class="info">
        {#if retentionDuration}
          Your name will be kept stored for {retentionDuration} from now, after which
          it will be removed and your results anonymized.
        {:else}
          Your name will be removed and your results anonymized shortly.
        {/if}
      </p>
    {/if}
    <select
      name="compClassId"
      required
      value={selectedCompClass?.id}
      aria-label="Category"
      disabled={contender.scrubbedAt !== undefined}
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
        disabled={contender.scrubbedAt !== undefined ||
          patchContender.isPending ||
          scrubContender.isPending}
        >{contender.entered ? "Update" : "Register"}</button
      >
    </div>
  </form>
{/if}

<style>
  .info {
    padding: var(--wa-space-s);
    border: 1px solid;
    border-radius: 0.25rem;
    margin: 0;
  }

  form {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
  }

  .name-row,
  .actions {
    display: flex;
    align-items: center;
    gap: var(--wa-space-xs);
  }

  .name-row input {
    flex-grow: 1;
  }

  .actions {
    justify-content: space-between;
  }
</style>
