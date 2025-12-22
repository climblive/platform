<script lang="ts">
  import ResultListTable from "@/components/ResultListTable.svelte";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { ApiClient } from "@climblive/lib";
  import { ScoreboardProvider } from "@climblive/lib/components";
  import { toastError } from "@climblive/lib/utils";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const handleDownloadResults = async () => {
    try {
      const blob = await ApiClient.getInstance().downloadResults(contestId);

      const a = document.createElement("a");
      document.body.appendChild(a);
      a.setAttribute("style", "display: none");

      const url = window.URL.createObjectURL(blob);
      a.href = url;
      a.download = `contest_${contestId}_results.xlsx`;
      a.click();

      window.URL.revokeObjectURL(url);
    } catch {
      toastError("Failed to download results.");
    }
  };
</script>

<section>
  <div class="controls">
    <wa-button appearance="outlined" onclick={handleDownloadResults}
      >Download results
      <wa-icon name="file-excel" slot="start"></wa-icon>
    </wa-button>

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
