<script lang="ts">
  import {
    getCompClassesQuery,
    getContenderQuery,
    patchContenderMutation,
  } from "@climblive/lib/queries";

  type Props = {
    contestId: number;
    contenderId: number;
  };

  const { contestId, contenderId }: Props = $props();

  let form: HTMLFormElement | undefined = $state();

  const contenderQuery = $derived(getContenderQuery(contenderId));
  const compClassesQuery = $derived(getCompClassesQuery(contestId));
  const patchContender = $derived(patchContenderMutation(contenderId));

  let contender = $derived($contenderQuery.data);
  let compClasses = $derived($compClassesQuery.data);
  let selectedCompClass = $derived(
    compClasses?.find(({ id }) => id === contender?.compClassId),
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
      $patchContender.mutate({
        ...contender,
        name,
        publicName: name,
        compClassId: Number(compClassId),
      });
    }
  };
</script>

{#if compClasses && contender}
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
      aria-label="Competition class"
    >
      {#each compClasses as compClass (compClass.id)}
        <option value={compClass.id}>{compClass.name}</option>
      {/each}
    </select>
    <button type="submit" disabled={$patchContender.isPending}
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
