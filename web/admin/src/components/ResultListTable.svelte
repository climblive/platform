<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import "@awesome.me/webawesome/dist/components/qr-code/qr-code.js";
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

  const columns: ColumnDefinition<ScoreboardEntry>[] = [
    {
      label: "Name",
      mobile: true,
      render: renderName,
      width: "3fr",
    },
    {
      label: "Score",
      mobile: true,
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
      mobile: false,
      render: renderFinalist,
      width: "max-content",
      align: "right",
    },
  ];

  $effect(() => console.log($scoreboard));
</script>

{#snippet renderName({
  contenderId,
  publicName,
  disqualified,
}: ScoreboardEntry)}
  {@const contender = contenders.get(contenderId)}
  {#if contender}
    <a href={`/${contender.registrationCode}`} target="_blank">
      {#if disqualified}
        <strike>{publicName}</strike>
      {:else}
        {publicName}
      {/if}
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

{#snippet renderFinalist({ score, withdrawnFromFinals }: ScoreboardEntry)}
  <wa-icon name={score?.finalist ? "medal" : "minus"}></wa-icon>
{/snippet}

<Table
  {columns}
  data={$scoreboard.get(compClassId) ?? []}
  getId={({ contenderId }) => contenderId}
></Table>
