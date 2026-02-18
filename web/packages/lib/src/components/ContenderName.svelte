<script lang="ts">
  import "@awesome.me/webawesome/dist/components/tooltip/tooltip.js";
  import { format } from "date-fns";
  import { maskScrubbedName } from "../utils";

  interface Props {
    id: number;
    name: string | undefined;
    scrubbedAt: Date | undefined;
  }

  const { id, name, scrubbedAt }: Props = $props();

  const displayName = $derived(
    scrubbedAt ? maskScrubbedName(id) : (name ?? ""),
  );

  const tooltipId = $props.id();
</script>

{#if scrubbedAt}
  <wa-tooltip for={tooltipId}
    >Anonymized on {format(scrubbedAt, "yyyy-MM-dd HH:mm")}</wa-tooltip
  >
{/if}

<span id={tooltipId}>{displayName}</span>
