<script lang="ts">
  import logoUrl from "@/static/logo.svg";
  import SvgQR from "@svelte-put/qr/svg/QR.svelte";

  type Props = {
    contestName: string;
    registrationCode: string;
  };

  let { contestName, registrationCode }: Props = $props();
</script>

<section>
  <div class="left">
    <h1>{contestName}</h1>
    <span class="link">Register at {location.protocol}//{location.host}</span>
    <span class="code">{registrationCode}</span>
  </div>
  <img src={logoUrl} alt="ClimbLive Logo" />
  <SvgQR
    data={`${location.protocol}//${location.host}/${registrationCode}`}
    width={80}
    margin={0}
    shape="circle"
    errorCorrectionLevel="Q"
  />
</section>

<style>
  section {
    width: 18cm;
    border: 1px solid var(--wa-color-primary-600);
    border-radius: var(--wa-border-radius-small);
    margin-block: var(--wa-spacing-medium);
    padding: var(--wa-spacing-medium);
    display: flex;
    justify-content: space-between;
    align-items: center;

    & .left {
      flex-grow: 1;
      display: flex;
      flex-direction: column;
      gap: var(--wa-spacing-x-small);
    }

    & h1 {
      font-size: var(--wa-font-size-large);
      margin: 0;
    }

    & .link {
      font-size: var(--wa-font-size-small);
    }

    & .code {
      font-size: var(--wa-font-size-x-large);
      font-weight: bold;
      font-family: monospace;
      letter-spacing: 0.2rem;
    }

    & img {
      height: 1.5rem;
      margin: auto var(--wa-spacing-medium) var(--wa-spacing-2x-small) auto;
    }
  }

  section:nth-of-type(7n) {
    break-after: page;
  }
</style>
