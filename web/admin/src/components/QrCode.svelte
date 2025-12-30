<script lang="ts">
  import logoUrl from "@/assets/logo.svg";
  import SvgQR from "@svelte-put/qr/svg/QR.svelte";

  type Props = {
    registrationCode: string;
    width: number;
    fill?: string;
  };

  const { registrationCode, width, fill = "black" }: Props = $props();
</script>

<SvgQR
  data={`${location.protocol}//${location.host}/${registrationCode}`}
  {width}
  margin={0}
  shape="circle"
  errorCorrectionLevel="Q"
  logo={logoUrl}
  anchorOuterFill={fill}
  moduleFill={fill}
  anchorInnerFill="#ea384c"
/>
