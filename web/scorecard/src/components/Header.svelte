<script lang="ts">
  import { ContenderName, FullLogo } from "@climblive/lib/components";
  import { navigate } from "svelte-routing";

  interface Props {
    registrationCode: string;
    contestName: string;
    compClassName: string | undefined;
    contenderId: number;
    contenderName: string | undefined;
    contenderScrubbedAt: Date | undefined;
  }

  const {
    registrationCode,
    contestName,
    compClassName,
    contenderId,
    contenderName,
    contenderScrubbedAt,
  }: Props = $props();
</script>

<header>
  <div class="top">
    <div class="logo">
      <FullLogo />
    </div>

    <wa-button
      size="small"
      onclick={() => navigate(`/${registrationCode}/edit`)}
      appearance="plain"
    >
      <wa-icon name="gear" label="Edit"></wa-icon>
    </wa-button>
  </div>
  <div>
    <h1>
      <ContenderName
        id={contenderId}
        name={contenderName}
        scrubbedAt={contenderScrubbedAt}
        withTooltip
      />
    </h1>
    <p class="subtitle">
      {contestName}<span class="separator"> / </span>{compClassName}
    </p>
  </div>
</header>

<style>
  header {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-xs);
    padding-block-start: var(--wa-space-l);
  }

  .logo {
    height: var(--wa-font-size-l);
  }

  .top {
    display: flex;
    align-items: center;
    justify-content: space-between;
    overflow-y: visible;
    height: var(--wa-font-size-l);
  }

  h1 {
    margin: 0;
    font-size: var(--wa-font-size-xl);
    font-weight: var(--wa-font-weight-bold);
    line-height: var(--wa-line-height-condensed);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .subtitle {
    margin: 0;
    margin-block-start: var(--wa-space-2xs);
    font-size: var(--wa-font-size-m);
    color: var(--wa-color-text-quiet);
    line-height: var(--wa-line-height-condensed);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .separator {
    margin-inline: var(--wa-space-2xs);
  }
  wa-button {
    flex-shrink: 0;
  }
</style>
