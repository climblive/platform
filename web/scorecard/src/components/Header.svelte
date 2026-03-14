<script lang="ts">
  import { type ContestState } from "@climblive/lib/types";
  import { navigate } from "svelte-routing";

  interface Props {
    registrationCode: string;
    contestName: string;
    compClassName: string | undefined;
    contenderName: string | undefined;
    contestState: ContestState;
  }

  let {
    registrationCode,
    contestName,
    compClassName,
    contenderName,
    contestState,
  }: Props = $props();
</script>

<header>
  <div class="info">
    <h1>{contestName}</h1>
    <p class="subtitle">
      {contenderName}
      {#if compClassName}
        <span class="separator">&middot;</span>
        <span class="comp-class">{compClassName}</span>
      {/if}
    </p>
  </div>
  <wa-button
    size="small"
    onclick={() => navigate(`/${registrationCode}/edit`)}
    disabled={contestState === "ENDED"}
    appearance="plain"
  >
    <wa-icon name="gear" label="Settings"></wa-icon>
  </wa-button>
</header>

<style>
  header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--wa-space-s);
    padding: var(--wa-space-s) 0;
  }

  .info {
    min-width: 0;
  }

  h1 {
    margin: 0;
    font-size: var(--wa-font-size-l);
    font-weight: var(--wa-font-weight-bold);
    line-height: var(--wa-line-height-condensed);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .subtitle {
    margin: 0;
    font-size: var(--wa-font-size-s);
    color: var(--wa-color-text-quiet);
    line-height: var(--wa-line-height-condensed);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .separator {
    margin-inline: var(--wa-space-3xs);
  }

  .comp-class {
    font-weight: var(--wa-font-weight-bold);
  }

  wa-button::part(base) {
    padding: 0;
  }

  wa-button {
    font-size: var(--wa-font-size-xl);
    flex-shrink: 0;
  }
</style>
