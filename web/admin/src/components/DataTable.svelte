<script lang="ts" generics="TData, TValue">
  import {
    type ColumnDef,
    getCoreRowModel,
    type VisibilityState,
  } from "@tanstack/table-core";
  import { onMount } from "svelte";
  import {
    createSvelteTable,
    FlexRender,
  } from "../$lib/components/ui/data-table/index.js";
  import * as Table from "../$lib/components/ui/table/index.js";

  type DataTableProps<TData, TValue> = {
    columns: ColumnDef<TData, TValue>[];
    data: TData[];
  };

  let { data, columns }: DataTableProps<TData, TValue> = $props();

  let columnVisibility = $state<VisibilityState>({});

  const table = createSvelteTable({
    get data() {
      return data;
    },
    columns,
    getCoreRowModel: getCoreRowModel(),
    initialState: {
      columnVisibility: {
        status: true,
      },
    },
    onColumnVisibilityChange: (updater) => {
      if (typeof updater === "function") {
        columnVisibility = updater(columnVisibility);
      } else {
        columnVisibility = updater;
      }
    },
    state: {
      get columnVisibility() {
        return columnVisibility;
      },
    },
  });

  function handleResize() {
    const mobile = window.innerWidth < 768;

    for (const column of table.getAllColumns()) {
      if (mobile && column.columnDef.enableHiding) {
        column.toggleVisibility(false);
      } else {
        column.toggleVisibility(true);
      }
    }
  }

  onMount(() => {
    handleResize();
  });
</script>

<svelte:window onresize={handleResize} />

<Table.Root>
  <Table.Header>
    {#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
      <Table.Row>
        {#each headerGroup.headers as header (header.id)}
          <Table.Head>
            {#if !header.isPlaceholder}
              <FlexRender
                content={header.column.columnDef.header}
                context={header.getContext()}
              />
            {/if}
          </Table.Head>
        {/each}
      </Table.Row>
    {/each}
  </Table.Header>
  <Table.Body>
    {#each table.getRowModel().rows as row (row.id)}
      <Table.Row>
        {#each row.getVisibleCells() as cell (cell.id)}
          <Table.Cell>
            <FlexRender
              content={cell.column.columnDef.cell}
              context={cell.getContext()}
            />
          </Table.Cell>
        {/each}
      </Table.Row>
    {/each}
  </Table.Body>
</Table.Root>
