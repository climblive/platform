<script lang="ts">
  import { LabeledText } from "@climblive/lib/components";
  import type { RaffleWinner } from "@climblive/lib/models";
  import { getRaffleWinnersByContestQuery } from "@climblive/lib/queries";
  import { format } from "date-fns";

  interface Props {
    contestId: number;
  }

  let { contestId }: Props = $props();

  const raffleWinnersQuery = $derived(
    getRaffleWinnersByContestQuery(contestId),
  );
  const raffleWinners = $derived(raffleWinnersQuery.data);

  const sortedRaffleWinners = $derived.by(() => {
    const winners = [...(raffleWinners ?? [])];
    winners.sort((a, b) => {
      return b.timestamp.getTime() - a.timestamp.getTime();
    });
    return winners;
  });
</script>

{#if sortedRaffleWinners && sortedRaffleWinners.length > 0}
  <section>
    <h2>Raffle Winners</h2>
    <div class="winners">
      {#each sortedRaffleWinners as winner (winner.id)}
        <div class="winner">
          <LabeledText label={format(winner.timestamp, "yyyy-MM-dd HH:mm")}>
            {winner.contenderName}
          </LabeledText>
        </div>
      {/each}
    </div>
  </section>
{/if}

<style>
  section {
    padding: var(--wa-space-m);
    background-color: var(--wa-color-surface-default);
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-surface-border);
    border-radius: var(--wa-border-radius-m);
    font-size: var(--wa-font-size-s);
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
  }

  h2 {
    font-size: var(--wa-font-size-l);
    font-weight: var(--wa-font-weight-semibold);
    margin: 0;
  }

  .winners {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
  }

  .winner {
    padding: var(--wa-space-s);
    background-color: var(--wa-color-surface-subtle);
    border-radius: var(--wa-border-radius-s);
  }
</style>
