<script lang="ts">
  type Props = {
    code: string;
  };

  const { code }: Props = $props();

  const segments = $derived(code.match(/.{1,4}/g));
</script>

<span>{segments?.join("-")}</span>
