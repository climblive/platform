<script lang="ts">
  import logoUrl from "@/assets/logo.svg";
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
  <SvgQR
    data={`${location.protocol}//${location.host}/${registrationCode}`}
    width={80}
    margin={0}
    shape="circle"
    errorCorrectionLevel="Q"
    logo={logoUrl}
    anchorOuterFill="black"
    moduleFill="black"
    anchorInnerFill="#ea384c"
  />
</section>

<style>
  section {
    width: 18cm;
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-brand-border-normal);
    border-radius: var(--wa-border-radius-s);
    margin-block: var(--wa-space-m);
    padding: var(--wa-space-s);
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
  }

  section:nth-of-type(7n) {
    break-after: page;
  }
</style>
