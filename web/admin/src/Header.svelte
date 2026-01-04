<script lang="ts">
  import "@awesome.me/webawesome/dist/components/badge/badge.js";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { FullLogo } from "@climblive/lib/components";
  import { getPendingUnlockRequestsQuery } from "@climblive/lib/queries";
  import { getContext, onMount } from "svelte";
  import { navigate } from "svelte-routing";
  import type { Authenticator } from "./authenticator.svelte";

  let print = $state(false);

  const authenticator = getContext<Authenticator>("authenticator");
  const pendingRequestsQuery = $derived(getPendingUnlockRequestsQuery());
  const pendingRequests = $derived(pendingRequestsQuery.data);
  const pendingCount = $derived(pendingRequests?.length ?? 0);

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
      <wa-button
        onclick={() => navigate("./unlock-requests")}
        size="small"
        variant="neutral"
        appearance="filled-outlined"
        title="Unlock Requests"
        style="position: relative;"
      >
        <wa-icon name="lock-open"></wa-icon>
        {#if pendingCount > 0}
          <wa-badge
            variant="danger"
            size="small"
            attention="bounce"
            style="position: absolute; top: -8px; right: -8px;"
          >
            {pendingCount}
          </wa-badge>
        {/if}
      </wa-button>
      <wa-button
        onclick={() => navigate("./help")}
        size="small"
        variant="success"
        appearance="filled-outlined"
        ><wa-icon name="headset"></wa-icon></wa-button
      >
      <wa-button
        size="small"
        appearance="outlined"
        onclick={authenticator.logout}
      >
        Sign out<wa-icon slot="start" name="right-from-bracket"
        ></wa-icon></wa-button
      >
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

  @media print {
    header {
      display: none;
    }
  }
</style>
