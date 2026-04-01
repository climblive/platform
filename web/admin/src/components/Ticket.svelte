<script lang="ts">
  import QrCode from "./QrCode.svelte";

  type Props = {
    contestName: string;
    registrationCode: string;
    ticketNumber: number;
  };

  let { contestName, registrationCode, ticketNumber }: Props = $props();
</script>

<section>
  <div class="left">
    <h1>{contestName}</h1>
    <span class="link">
      Scan the QR code or register at {location.protocol}//{location.host}
      using your <em>secret code:</em></span
    >
    <span class="code">{registrationCode}</span>
  </div>
  <QrCode {registrationCode} width={80}></QrCode>
  <span class="ticket-number">#{ticketNumber}</span>
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
    gap: var(--wa-space-m);

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

    :global(& > svg) {
      flex-shrink: 0;
    }
  }

  .ticket-number {
    writing-mode: vertical-rl;
    font-size: var(--wa-font-size-xs);
    color: var(--wa-color-text-quiet);
    font-family: monospace;
    flex-shrink: 0;
  }
</style>
