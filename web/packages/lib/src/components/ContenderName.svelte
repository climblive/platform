<script lang="ts">
  import "@awesome.me/webawesome/dist/components/popover/popover.js";
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
  <wa-popover for={tooltipId}>
    The name of the contender was removed and anonymized on {format(
      scrubbedAt,
      "yyyy-MM-dd HH:mm",
    )}.

    <wa-button data-popover="close" variant="primary" size="small"
      >Got it!</wa-button
    >
  </wa-popover>
{/if}

<span id={tooltipId}>{displayName}</span>

<style>
  wa-popover {
    --max-width: 300px;
  }

  wa-popover::part(body) {
    margin-inline-start: var(--wa-space-l);
  }

  wa-button {
    margin-block-start: var(--wa-space-m);
  }
</style>
