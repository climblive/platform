<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { FullLogo } from "@climblive/lib/components";
  import { getHealthQuery, getSelfQuery } from "@climblive/lib/queries";
  import { getContext, onMount } from "svelte";
  import { navigate } from "svelte-routing";
  import type { Authenticator } from "./authenticator.svelte";

  let print = $state(false);

  const authenticator = getContext<Authenticator>("authenticator");

  const selfQuery = $derived(getSelfQuery());
  const self = $derived(selfQuery.data);

  const healthQuery = $derived(
    self?.admin ? getHealthQuery() : { data: undefined },
  );
  const health = $derived(healthQuery.data);

  const hasIssues = $derived(
    health !== undefined &&
      (!health.scoreEngineManager.healthy ||
        !health.scoreKeeper.healthy ||
        !health.scrubber.healthy),
  );

  onMount(() => {
    const urlParams = new URLSearchParams(window.location.search);
    if (urlParams.get("print") !== null) {
      print = true;
    }
  });
</script>

{#if !print}
  <header>
    <div>
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div class="logo" onclick={() => navigate("./")}>
        <FullLogo />
      </div>
      <div class="right-actions">
        <wa-button
          onclick={() => navigate("./help")}
          size="s"
          variant="success"
          appearance="filled-outlined"
          ><wa-icon name="headset"></wa-icon></wa-button
        >
        {#if self?.admin}
          <div class="health-btn">
            <wa-button
              onclick={() => navigate("./health")}
              size="s"
              appearance="outlined"
            >
              <wa-icon name="heart-pulse"></wa-icon>
            </wa-button>
            {#if hasIssues}
              <span class="issue-badge"></span>
            {/if}
          </div>
        {/if}
        <wa-button size="s" appearance="outlined" onclick={authenticator.logout}>
          Sign out<wa-icon slot="start" name="right-from-bracket"
          ></wa-icon></wa-button
        >
      </div>
    </div>
  </header>
{/if}

<style>
  header {
    background-color: var(--wa-color-surface-lowered);
  }

  div {
    margin: 0 auto;
    max-width: 1024px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding-inline-end: var(--wa-space-m);
    height: 3.5rem;
    gap: var(--wa-space-xs);
  }

  .logo {
    text-align: left;
    height: var(--wa-font-size-xl);
    color: var(--wa-color-text-normal);
    padding-left: var(--wa-space-xs);
    flex-shrink: 0;
    margin-inline-start: var(--wa-space-xs);
    cursor: pointer;
  }

  .right-actions {
    display: flex;
    align-items: center;
    gap: var(--wa-space-xs);
    max-width: unset;
    margin: unset;
    height: unset;
    padding: unset;
    justify-content: flex-end;
  }

  .health-btn {
    position: relative;
    display: inline-flex;
  }

  .issue-badge {
    position: absolute;
    top: 2px;
    right: 2px;
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background-color: var(--wa-color-danger);
    animation: pulse 1.5s ease-in-out infinite;
  }

  @keyframes pulse {
    0%,
    100% {
      opacity: 1;
      transform: scale(1);
    }
    50% {
      opacity: 0.5;
      transform: scale(1.4);
    }
  }

  @media print {
    header {
      display: none;
    }
  }
</style>
