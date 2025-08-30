<script lang="ts">
  import ResultListTable from "@/components/ResultListTable.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { ScoreboardProvider } from "@climblive/lib/components";
  import { getApiUrl } from "@climblive/lib/utils";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();
</script>

<section>
  <div class="controls">
    <a href={`${getApiUrl()}/contests/${contestId}/results`}>
      <wa-button appearance="outlined"
        >Download results
        <wa-icon name="file-excel" slot="start"></wa-icon>
      </wa-button>
    </a>

    <a href={`/scoreboard/${contestId}`} target="_blank">
      <wa-button appearance="outlined">
        <wa-icon slot="start" name="arrow-up-right-from-square"></wa-icon>
        Open public scoreboard
      </wa-button>
    </a>
  </div>

  <ScoreboardProvider {contestId}>
    {#snippet children({ scoreboard, loading })}
      <ResultListTable {contestId} {scoreboard} {loading}></ResultListTable>
    {/snippet}
  </ScoreboardProvider>
</section>

<style>
  section {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
  }

  .controls {
    display: flex;
    gap: var(--wa-space-xs);
    flex-wrap: wrap;
  }
</style>
