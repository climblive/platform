<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import type WaInput from "@awesome.me/webawesome/dist/components/input/input.js";
  import "@awesome.me/webawesome/dist/components/option/option.js";
  import "@awesome.me/webawesome/dist/components/select/select.js";
  import type WaSelect from "@awesome.me/webawesome/dist/components/select/select.js";
  import "@awesome.me/webawesome/dist/components/switch/switch.js";
  import type WaSwitch from "@awesome.me/webawesome/dist/components/switch/switch.js";
  import { Table, type ColumnDefinition } from "@climblive/lib/components";
  import { checked, value } from "@climblive/lib/forms";
  import type { ScoreboardEntry } from "@climblive/lib/models";
  import { getCompClassesQuery } from "@climblive/lib/queries";
  import { ordinalSuperscript } from "@climblive/lib/utils";
  import { Link } from "svelte-routing";
  import type { Readable } from "svelte/store";
  import Loader from "./Loader.svelte";

  interface Props {
    contestId: number;
    scoreboard: Readable<Map<number, ScoreboardEntry[]>>;
    loading: boolean;
  }

  const { contestId, scoreboard, loading }: Props = $props();

  let tableData = $state<ScoreboardEntry[]>([]);

  let compClassSelector: WaSelect | undefined = $state();
  let quickFilter: WaInput | undefined = $state();

  const compClassesQuery = $derived(getCompClassesQuery(contestId));

  const compClasses = $derived($compClassesQuery.data);

  let filterText = $state<string>();
  let selectedCompClassId: number | undefined = $state();
  let liveSwitch: WaSwitch | undefined = $state();
  let live = $state(true);

  $effect(() => {
    if (
      compClasses &&
      compClasses.length > 0 &&
      selectedCompClassId === undefined
    ) {
      selectedCompClassId = compClasses[0].id;
    }
  });

  const toggleLive = () => {
    if (!liveSwitch) {
      return;
    }

    live = Boolean(liveSwitch.checked);
  };

  $effect(() => {
    if (!live) {
      return;
    }

    if (selectedCompClassId === undefined) {
      return;
    }

    let scores = [...($scoreboard.get(selectedCompClassId) ?? [])];

    if (filterText) {
      const search = filterText.toLowerCase();

      scores = scores.filter(({ publicName }) =>
        publicName.toLowerCase().includes(search),
      );
    }

    scores.sort(
      (a: ScoreboardEntry, b: ScoreboardEntry) =>
        (a.score?.rankOrder ?? Infinity) - (b.score?.rankOrder ?? Infinity),
    );

    tableData = scores;
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
      align: "left",
    },
  ];
</script>

{#snippet renderName({
  contenderId,
  publicName,
  disqualified,
}: ScoreboardEntry)}
  <Link to={`./contenders/${contenderId}`}>
    {#if disqualified}
      <del>{publicName}</del>
    {:else}
      {publicName}
    {/if}
  </Link>
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

{#if compClasses}
  <div class="controls">
    <wa-input
      bind:this={quickFilter}
      size="small"
      label="Quick filter"
      placeholder="Search by name..."
      oninput={() => {
        filterText = quickFilter?.value ?? "";
      }}
      with-clear
    ></wa-input>
    <wa-select
      bind:this={compClassSelector}
      size="small"
      label="Class filter"
      {@attach value(selectedCompClassId)}
      onchange={() => {
        selectedCompClassId = Number(compClassSelector?.value);
      }}
    >
      {#each compClasses as compClass (compClass.id)}
        <wa-option value={compClass.id}>{compClass.name}</wa-option>
      {/each}
    </wa-select>
  </div>
{/if}

<wa-switch bind:this={liveSwitch} {@attach checked(live)} onchange={toggleLive}
  >Live</wa-switch
>

{#if loading}
  <Loader />
{:else}
  <Table {columns} data={tableData} getId={({ contenderId }) => contenderId}
  ></Table>
{/if}

<style>
  wa-switch {
    margin-left: auto;
  }

  .controls {
    display: flex;
    gap: var(--wa-space-m);
    justify-content: space-evenly;

    & * {
      width: 100%;
    }
  }
</style>
