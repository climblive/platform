<script lang="ts">
  import "@awesome.me/webawesome/dist/components/tooltip/tooltip.js";
  import { maskScrubbedName } from "../utils";

  interface Props {
    id: number;
    name?: string;
    scrubbedAt?: Date;
  }

  const { id, name, scrubbedAt }: Props = $props();

  const isScrubbed = $derived(scrubbedAt !== undefined || name === "");
  const displayName = $derived(
    isScrubbed ? maskScrubbedName(id) : (name ?? ""),
  );

  const tooltipId = $props.id();
</script>

{#if isScrubbed}
  <wa-tooltip for={tooltipId}>Name has been scrubbed</wa-tooltip>
  <span id={tooltipId}>{displayName}</span>
{:else}
  {displayName}
{/if}
