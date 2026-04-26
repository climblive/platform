<script lang="ts">
  import logoUrl from "@/assets/logo.svg";
  import { ContenderName } from "@climblive/lib/components";
  import { Link, navigate } from "svelte-routing";

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
  <div class="identity">
    <Link to="/">
      <img class="logo" src={logoUrl} alt="ClimbLive logo" />
    </Link>
    <div class="info">
      <h1>
        <ContenderName
          id={contenderId}
          name={contenderName}
          scrubbedAt={contenderScrubbedAt}
          withTooltip
        />
      </h1>
      <p class="subtitle">
        <span>{compClassName}</span><span class="separator">–</span
        >{contestName}
      </p>
    </div>
  </div>
  <wa-button
    size="medium"
    onclick={() => navigate(`/${registrationCode}/edit`)}
    appearance="plain"
  >
    <wa-icon name="gear" label="Edit"></wa-icon>
  </wa-button>
</header>

<style>
  header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--wa-space-s);
    padding-block-start: var(--wa-space-l);
  }

  .identity {
    display: flex;
    align-items: center;
    gap: var(--wa-space-xs);
    min-width: 0;
  }

  .info {
    min-width: 0;
  }

  .logo {
    width: calc(var(--wa-font-size-l) * 2);
    height: calc(var(--wa-font-size-l) * 2);
    flex-shrink: 0;
  }

  h1 {
    margin: 0;
    margin-block-start: -0.25rem;
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

  wa-button {
    flex-shrink: 0;
  }
</style>
