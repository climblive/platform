<script lang="ts">
  import logoUrl from "@/static/logo.svg";
  import SvgQR from "@svelte-put/qr/svg/QR.svelte";

  type Props = {
    registrationCode: string;
  };

  let { registrationCode }: Props = $props();
</script>

<section>
  <img src={logoUrl} alt="ClimbLive Logo" height="32" />
  {registrationCode}
  <SvgQR
    data={`${location.protocol}//${location.host}/${registrationCode}`}
    width={64}
  />
</section>

<style>
  section {
    border: 1px solid var(--sl-color-primary-600);
    border-radius: var(--sl-border-radius-medium);
    margin-block: var(--sl-spacing-medium);
    padding: var(--sl-spacing-medium);
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  section:nth-of-type(7n) {
    break-after: page;
  }
</style>
