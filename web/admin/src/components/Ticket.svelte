<script lang="ts">
  import QRCode from "qrcode-svg";

  type Props = {
    registrationCode: string;
  };

  let { registrationCode }: Props = $props();

  const svg = $derived(
    new QRCode({
      content: `${location.protocol}//${location.host}/${registrationCode}`,
      padding: 0,
      width: 64,
      height: 64,
      color: "#000000",
      background: "#ffffff",
      ecl: "M",
    }).svg(),
  );
</script>

<section>
  {@html svg}
  {registrationCode}
</section>

<style>
  section {
    border: 1px solid var(--sl-color-primary-600);
    border-radius: var(--sl-border-radius-medium);
    margin-block: var(--sl-spacing-medium);
    padding: var(--sl-spacing-medium);
  }

  section:nth-of-type(7n) {
    break-after: page;
  }
</style>
