<script lang="ts">
  import type { CreateContendersArguments } from "@climblive/lib/models";
  import {
    createContendersMutation,
    getContendersByContestQuery,
  } from "@climblive/lib/queries";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import "@shoelace-style/shoelace/dist/components/icon/icon.js";
  import type SlRange from "@shoelace-style/shoelace/dist/components/range/range.js";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  let range = $state<SlRange | undefined>();
  const contendersQuery = getContendersByContestQuery(contestId);
  const createContenders = createContendersMutation(contestId);

  let contenders = $derived($contendersQuery.data);

  let remainingCodes = $derived(
    contenders === undefined ? undefined : 500 - contenders.length,
  );

  const increments = [1, 10, 100];

  const addContenders = async (args: CreateContendersArguments) => {
    $createContenders.mutate(args, {
      onSuccess: () => {
        if (range) {
          range.value = 0;
        }
      },
    });
  };
</script>

{#if contenders}
  <p>
    You have {remainingCodes} codes remaining out of your maximum 500.
  </p>
  <section>
    {#each increments as increment (increment)}
      <sl-button
        size="small"
        type="button"
        variant="primary"
        loading={$createContenders.isPending}
        disabled={!remainingCodes || remainingCodes < increment}
        onclick={() => addContenders({ number: increment })}
      >
        <sl-icon slot="prefix" name="plus-lg"></sl-icon>
        Add {increment} code{#if increment != 1}s{/if}
      </sl-button>
    {/each}
  </section>

  <ul>
    {#each contenders as contender (contender.id)}
      <li>{contender.registrationCode}</li>
    {/each}
  </ul>
{/if}

<style>
  section {
    display: flex;
    gap: var(--sl-spacing-x-small);
  }
</style>
