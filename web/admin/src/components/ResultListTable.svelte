<script lang="ts">
  import { Table, type ColumnDefinition } from "@climblive/lib/components";
  import type { Contender, ScoreboardEntry } from "@climblive/lib/models";
  import { ordinalSuperscript } from "@climblive/lib/utils";
  import type { Readable } from "svelte/store";

  interface Props {
    scoreboard: Readable<Map<number, ScoreboardEntry[]>>;
    contenders: Map<number, Contender>;
    compClassId: number;
  }

  const { scoreboard, contenders, compClassId }: Props = $props();

  const data = $derived.by(() => {
    const scores = [...($scoreboard.get(compClassId) ?? [])];
    scores.sort(
      (a: ScoreboardEntry, b: ScoreboardEntry) =>
        (a.score?.rankOrder ?? Infinity) - (b.score?.rankOrder ?? Infinity),
    );
    return scores;
  });

  const columns: ColumnDefinition<ScoreboardEntry>[] = [
    {
      label: "Name",
      mobile: true,
      render: renderName,
      width: "3fr",
    },
    {
      label: "Score",
      mobile: false,
      render: renderScore,
      width: "max-content",
      align: "right",
    },
    {
      label: "Placement",
      mobile: true,
      render: renderPlacement,
      width: "max-content",
      align: "right",
    },
    {
      label: "Finalist",
      mobile: true,
      render: renderFinalist,
      width: "max-content",
      align: "right",
    },
  ];
</script>

{#snippet renderName({
  contenderId,
  publicName,
  disqualified,
}: ScoreboardEntry)}
  {@const contender = contenders.get(contenderId)}
  {#if contender}
    <a href={`/${contender.registrationCode}`} target="_blank">
      <wa-button appearance="plain" variant="brand" size="small">
        <wa-icon slot="start" name="arrow-up-right-from-square"></wa-icon>
        {#if disqualified}
          <strike>{publicName}</strike>
        {:else}
          {publicName}
        {/if}
      </wa-button>
    </a>
  {/if}
{/snippet}

{#snippet renderScore({ score }: ScoreboardEntry)}
  {#if score}
    {score.score} pts
  {/if}
{/snippet}

{#snippet renderPlacement({ score }: ScoreboardEntry)}
  {#if score}
    {score.placement}<sup>{ordinalSuperscript(score.placement)}</sup>
  {/if}
{/snippet}

{#snippet renderFinalist({ score }: ScoreboardEntry)}
  <wa-icon name={score?.finalist ? "medal" : "minus"}></wa-icon>
{/snippet}

{#if data && data.length > 0}
  <Table {columns} {data} getId={({ contenderId }) => contenderId}></Table>
{/if}
