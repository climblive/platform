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
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-brand-border-normal);
    border-radius: var(--wa-border-radius-m);
    margin-block: var(--wa-space-m);
    padding: var(--wa-space-m);
    display: flex;
    justify-content: space-between;
    align-items: center;

    & .left {
      flex-grow: 1;
      display: flex;
      flex-direction: column;
      gap: var(--wa-space-xs);
    }

    & h1 {
      font-size: var(--wa-font-size-l);
      margin: 0;
    }

    & .link {
      font-size: var(--wa-font-size-s);
    }

    & .code {
      font-size: var(--wa-font-size-xl);
      font-weight: var(--wa-font-weight-bold);
      font-family: monospace;
      letter-spacing: 0.2rem;
    }

    & img {
      height: 1.5rem;
      margin: auto var(--wa-space-m) var(--wa-space-2xs) auto;
    }
  }

  section:nth-of-type(7n) {
    break-after: page;
  }
</style>
