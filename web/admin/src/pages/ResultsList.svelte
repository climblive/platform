<script lang="ts">
  import ResultListTable from "@/components/ResultListTable.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/qr-code/qr-code.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import type WaSelect from "@awesome.me/webawesome/dist/components/select/select.js";
  import { ScoreboardProvider } from "@climblive/lib/components";
  import { value } from "@climblive/lib/forms";
  import {
    getCompClassesQuery,
    getContendersByContestQuery,
  } from "@climblive/lib/queries";
  import { getApiUrl } from "@climblive/lib/utils";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  let compClassSelector: WaSelect | undefined = $state();

  const contendersQuery = $derived(getContendersByContestQuery(contestId));
  const compClassesQuery = $derived(getCompClassesQuery(contestId));

  let contenders = $derived(
    new Map(
      $contendersQuery.data?.map((contender) => [contender.id, contender]) ??
        [],
    ),
  );
  let compClasses = $derived($compClassesQuery.data);

  let selectedCompClassId: number | undefined = $state();

  $effect(() => {
    if (
      compClasses &&
      compClasses.length > 0 &&
      selectedCompClassId === undefined
    ) {
      selectedCompClassId = compClasses[0].id;
    }
  });
</script>

<a href={`${getApiUrl()}/contests/${contestId}/results`}>
  <wa-button appearance="outlined"
    >Download results
    <wa-icon name="download" slot="start"></wa-icon>
  </wa-button>
</a>

{#if compClasses && compClasses.length > 1}
  <wa-select
    bind:this={compClassSelector}
    size="small"
    name="compClassId"
    label="Competition class"
    {@attach value(selectedCompClassId)}
    onchange={() => {
      selectedCompClassId = Number(compClassSelector?.value);
    }}
  >
    {#each compClasses as compClass (compClass.id)}
      <wa-option value={compClass.id}>{compClass.name}</wa-option>
    {/each}
  </wa-select>
{/if}

<ScoreboardProvider {contestId}>
  {#snippet children({ scoreboard })}
    {#if selectedCompClassId}
      <ResultListTable
        {scoreboard}
        {contenders}
        compClassId={selectedCompClassId}
      ></ResultListTable>
    {/if}
  {/snippet}
</ScoreboardProvider>

<style>
  wa-select {
    margin-top: var(--wa-space-m);
  }
</style>
