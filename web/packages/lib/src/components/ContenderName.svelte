<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/popover/popover.js";
  import { format } from "date-fns";
  import { maskScrubbedName } from "../utils";

  interface Props {
    id: number;
    name: string | undefined;
    scrubbedAt: Date | undefined;
    withTooltip?: boolean;
  }

  const { id, name, scrubbedAt, withTooltip = false }: Props = $props();

  const tooltipId = $props.id();
</script>

{#if scrubbedAt && withTooltip}
  <wa-popover for={tooltipId}>
    The contender name was removed and anonymized on {format(
      scrubbedAt,
      "yyyy-MM-dd HH:mm",
    )}.

    <wa-button data-popover="close" variant="primary" size="small"
      >Got it!</wa-button
    >
  </wa-popover>
{/if}

<span id={tooltipId}>
  {#if scrubbedAt}
    <pre>{maskScrubbedName(id)}</pre>
  {:else}
    {name ?? ""}
  {/if}
</span>

<style>
  wa-popover {
    --max-width: 300px;

    font-weight: var(--wa-font-weight-normal);
  }

  wa-popover::part(body) {
    margin-inline-start: var(--wa-space-l);
  }

  wa-button {
    margin-block-start: var(--wa-space-m);
  }
</style>
