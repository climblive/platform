<script lang="ts">
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import {
    getCompClassesQuery,
    getContenderQuery,
    patchContenderMutation,
    scrubContenderMutation,
  } from "@climblive/lib/queries";
  import { SyncedTime, toastUnexpectedError } from "@climblive/lib/utils";
  import { formatDistance, isBefore } from "date-fns";
  import { onMount } from "svelte";

  type Props = {
    contestId: number;
    contenderId: number;
  };

  const { contestId, contenderId }: Props = $props();

  let form: HTMLFormElement | undefined = $state();

  const contenderQuery = $derived(getContenderQuery(contenderId));
  const compClassesQuery = $derived(getCompClassesQuery(contestId));
  const patchContender = $derived(patchContenderMutation(contenderId));
  const scrubContender = $derived(scrubContenderMutation(contenderId));
  const time = new SyncedTime(60000);
  onMount(() => {
    time.start();
    return () => time.stop();
  });

  let contender = $derived(contenderQuery.data);
  let compClasses = $derived(compClassesQuery.data);
  let selectedCompClass = $derived(
    compClasses?.find(({ id }) => id === contender?.compClassId),
  );

  const retentionDuration = $derived.by(() =>
    contender?.scrubBefore && !isBefore(contender.scrubBefore, time.current)
      ? formatDistance(contender.scrubBefore, time.current)
      : undefined,
  );

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
    )
      return;
    scrubContender.mutate(undefined, {
      onSuccess: () => window.location.reload(),
      onError: () => toastUnexpectedError("Failed to remove your name."),
    });
  };
</script>

{#if compClasses && contender}
  {#if !contender.scrubbedAt}<wa-callout variant="neutral" size="s"
      ><wa-icon slot="icon" name="circle-info"
      ></wa-icon>{#if retentionDuration}Your name will be kept stored for {retentionDuration}
        from now, after which it will be removed and your results anonymized.{:else}Your
        name will be removed and your results anonymized shortly.{/if}</wa-callout
    >{/if}
  <form onsubmit={handleSubmit} bind:this={form}>
    <input
      required
      placeholder="Name"
      name="name"
      type="text"
      value={contender.name}
      aria-label="Name"
    />
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
  </form>
{/if}

<style>
  form {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
  }
</style>
