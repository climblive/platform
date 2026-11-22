<script lang="ts">
  import { FullLogo } from "@climblive/lib/components";
  import type { Category, Problem } from "@climblive/lib/models";
  import { localizedPointsValue } from "@climblive/lib/utils";
  import { onMount, tick } from "svelte";

  interface Props {
    competitionName: string;
    categories: Category[];
    problems: Problem[];
  }

  const { competitionName, categories, problems }: Props = $props();

  const hasZone1 = (problem: Problem) =>
    problem.zone1Enabled && (problem.pointsZone1 ?? 0) > 0;
  const hasZone2 = (problem: Problem) =>
    problem.zone2Enabled && (problem.pointsZone2 ?? 0) > 0;
  const hasFlash = (problem: Problem) => (problem.flashBonus ?? 0) > 0;

  let rowsPerColumn = $state(24);
  let tableContainers = $state<HTMLDivElement[]>([]);

  const sortedProblems = $derived(
    problems.toSorted((a, b) => a.number - b.number),
  );

  const pages = $derived.by(() => {
    const capacity = rowsPerColumn * 3;
    const pages: Problem[][][] = [];

    for (let start = 0; start < sortedProblems.length; start += capacity) {
      const page = sortedProblems.slice(start, start + capacity);
      const columns: Problem[][] = [[], [], []];

      for (let index = 0; index < page.length; index++) {
        columns[index % 3].push(page[index]);
      }

      pages.push(columns);
    }

    return pages;
  });

  const updateRowsPerColumn = () => {
    const container = tableContainers[0];
    const heading = container?.querySelector("thead");
    const row = container?.querySelector("tbody tr");

    if (!container || !heading || !row) {
      return;
    }

    const availableHeight =
      container.clientHeight - heading.getBoundingClientRect().height - 1;
    rowsPerColumn = Math.max(
      1,
      Math.floor(availableHeight / row.getBoundingClientRect().height),
    );
  };

  export const preparePrint = async () => {
    updateRowsPerColumn();
    await tick();
  };

  onMount(() => {
    const observer = new ResizeObserver(updateRowsPerColumn);
    observer.observe(tableContainers[0]);
    updateRowsPerColumn();

    return () => observer.disconnect();
  });
</script>

{#snippet achievement(points: number)}
  {@const value = localizedPointsValue(points)}
  <span class="achievement">
    <span class="checkbox" aria-hidden="true"></span>
    <span style:font-size={`${Math.min(6.5, 48 / value.length)}pt`}
      >{value}</span
    >
  </span>
{/snippet}

{#each pages as columns, pageIndex (pageIndex)}
  <section class="sheet" aria-label={`Scorecard page ${pageIndex + 1}`}>
    <header>
      <div class="logo" role="img" aria-label="ClimbLive"><FullLogo /></div>
      <h1>{competitionName}</h1>
      <div class="name">Name <span class="writing-line"></span></div>
      <div class="categories">
        <strong>Category</strong>
        {#each categories as category (category.id)}
          <span class="category">
            <span class="checkbox" aria-hidden="true"></span>{category.name}
          </span>
        {:else}
          <span class="writing-line"></span>
        {/each}
      </div>
      <p>
        Tick only your highest achievement for each problem. Flash = top on the
        first attempt. Leave unclimbed problems blank.
      </p>
    </header>

    <div class="problem-tables" bind:this={tableContainers[pageIndex]}>
      {#each columns as column, columnIndex (columnIndex)}
        {#if column.length > 0}
          {@const zone1 = column.some(hasZone1)}
          {@const zone2 = column.some(hasZone2)}
          {@const flash = column.some(hasFlash)}
          <table
            aria-label={`Problems ${column[0].number}–${column[column.length - 1].number}`}
          >
            <thead>
              <tr>
                <th scope="col">No.</th>
                {#if zone1}<th scope="col">Z1</th>{/if}
                {#if zone2}<th scope="col">Z2</th>{/if}
                <th scope="col">Top</th>
                {#if flash}<th scope="col">Flash</th>{/if}
              </tr>
            </thead>
            <tbody>
              {#each column as problem (problem.id)}
                {@const zone1Available = hasZone1(problem)}
                {@const zone2Available = hasZone2(problem)}
                <tr>
                  <th scope="row">{problem.number}</th>
                  {#if zone1}
                    <td
                      aria-label={zone1Available ? "Zone 1" : "Not available"}
                    >
                      {#if zone1Available}
                        {@render achievement(problem.pointsZone1 ?? 0)}
                      {:else}—{/if}
                    </td>
                  {/if}
                  {#if zone2}
                    <td
                      aria-label={zone2Available ? "Zone 2" : "Not available"}
                    >
                      {#if zone2Available}
                        {@render achievement(problem.pointsZone2 ?? 0)}
                      {:else}—{/if}
                    </td>
                  {/if}
                  <td aria-label="Top">
                    {#if problem.pointsTop > 0}
                      {@render achievement(problem.pointsTop)}
                    {:else}—{/if}
                  </td>
                  {#if flash}
                    <td
                      aria-label={hasFlash(problem) ? "Flash" : "Not available"}
                    >
                      {#if hasFlash(problem)}
                        {@render achievement(
                          problem.pointsTop + (problem.flashBonus ?? 0),
                        )}
                      {:else}—{/if}
                    </td>
                  {/if}
                </tr>
              {/each}
            </tbody>
          </table>
        {/if}
      {/each}
    </div>

    <footer>
      <span>Page {pageIndex + 1} of {pages.length}</span>
    </footer>
  </section>
{/each}

<style>
  .sheet {
    width: 190mm;
    height: 276mm;
    margin: 0 auto var(--wa-space-l);
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: 4mm;
    background: white;
    color: black;
    font-family: Arial, sans-serif;
    font-size: 10pt;
    line-height: 1.25;
  }

  header {
    display: flex;
    flex-direction: column;
    gap: 4mm;
  }

  .logo {
    height: 8mm;
    flex-shrink: 0;
    print-color-adjust: exact;
  }

  h1 {
    margin: 0;
    font-size: 16pt;
    overflow-wrap: anywhere;
  }

  .name {
    display: flex;
    gap: 3mm;
    font-weight: bold;
  }

  .writing-line {
    display: inline-block;
    flex: 1;
    min-width: 30mm;
    min-height: 6mm;
    border-bottom: 0.25mm solid black;
  }

  .categories {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 3mm 5mm;
  }

  .category {
    display: inline-flex;
    align-items: center;
    gap: 2mm;
    max-width: 100%;
    overflow-wrap: anywhere;
  }

  .checkbox {
    display: inline-block;
    width: 3.5mm;
    height: 3.5mm;
    flex-shrink: 0;
    border: 0.25mm solid black;
    vertical-align: middle;
  }

  .achievement {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 0.5mm;
    font-size: 6.5pt;
    line-height: 1;
    white-space: nowrap;
  }

  p {
    margin: 0;
    font-size: 9pt;
  }

  .problem-tables {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    align-items: start;
    gap: 4mm;
    flex: 1;
    min-height: 0;
  }

  table {
    width: 100%;
    table-layout: fixed;
    border-collapse: collapse;
    font-size: 9pt;
  }

  th,
  td {
    height: 8mm;
    padding: 0 1mm;
    border: 0.25mm solid black;
    text-align: center;
    vertical-align: middle;
  }

  footer {
    display: flex;
    justify-content: flex-end;
    font-size: 8pt;
  }

  @media print {
    .sheet {
      margin: 0;
      break-after: page;

      &:last-child {
        break-after: auto;
      }
    }
  }
</style>
