<script lang="ts">
  import { type WaSelectEvent } from "@awesome.me/webawesome";
  import "@awesome.me/webawesome/dist/components/badge/badge.js";
  import "@awesome.me/webawesome/dist/components/button-group/button-group.js";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/dropdown/dropdown.js";
  import WaDropdownItem from "@awesome.me/webawesome/dist/components/dropdown-item/dropdown-item.js";
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

  const healthQuery = $derived(self?.admin ? getHealthQuery() : undefined);
  const health = $derived(healthQuery?.data);

  const issues = $derived(
    health?.filter(({ healthy }) => !healthy).length ?? 0,
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
          variant="neutral"
          appearance="outlined"><wa-icon name="headset"></wa-icon></wa-button
        >
        {#if self?.admin}
          <wa-button
            onclick={() => navigate("./health")}
            size="s"
            appearance="outlined"
            variant={issues > 0 ? "danger" : "success"}
          >
            <wa-icon name="heart-pulse"></wa-icon>
          </wa-button>
        {/if}
        <wa-button-group label="Account actions">
          <wa-button
            size="s"
            appearance="outlined"
            aria-label="Open profile"
            onclick={() => navigate("/admin/profile")}
          >
            <wa-icon name="circle-user"></wa-icon>
          </wa-button>
          <wa-dropdown
            onwa-select={(event: WaSelectEvent) => {
              const { item } = event.detail;

              if (item instanceof WaDropdownItem && item.value === "sign-out") {
                authenticator.logout();
              }
            }}
          >
            <wa-button
              slot="trigger"
              size="s"
              appearance="outlined"
              aria-label="Open account actions"
            >
              <wa-icon name="chevron-down"></wa-icon>
            </wa-button>
            <wa-dropdown-item value="sign-out" variant="danger">
              <wa-icon slot="icon" name="right-from-bracket"></wa-icon>
              Sign out
            </wa-dropdown-item>
          </wa-dropdown>
        </wa-button-group>
      </div>
    </div>
  </header>
{/if}

<style>
  header {
    background-color: var(--wa-color-surface-lowered);
  }

  header > div {
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
  }

  @media print {
    header {
      display: none;
    }
  }
</style>
