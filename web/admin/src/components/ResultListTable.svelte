<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import "@awesome.me/webawesome/dist/components/qr-code/qr-code.js";
  import { Table, type ColumnDefinition } from "@climblive/lib/components";
  import type { ScoreboardEntry } from "@climblive/lib/models";
  import { ordinalSuperscript } from "@climblive/lib/utils";
  import type { Readable } from "svelte/store";

  interface Props {
    scoreboard: Readable<Map<number, ScoreboardEntry[]>>;
    compClassId: number;
  }

  const { scoreboard, compClassId }: Props = $props();

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
  ];
</script>

{#snippet renderName({ publicName }: ScoreboardEntry)}
  {publicName}
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

<Table
  {columns}
  data={$scoreboard.get(compClassId) ?? []}
  getId={({ contenderId }) => contenderId}
></Table>
