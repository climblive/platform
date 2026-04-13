<script lang="ts">
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import type WaDialog from "@awesome.me/webawesome/dist/components/dialog/dialog.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/number-input/number-input.js";
  import type WaNumberInput from "@awesome.me/webawesome/dist/components/number-input/number-input.js";
  import { value } from "@climblive/lib/forms";
  import type {
    Contender,
    CreateContendersArguments,
  } from "@climblive/lib/models";
  import { createContendersMutation } from "@climblive/lib/queries";
  import { toastError } from "@climblive/lib/utils";

  const maxTickets = 500;

  interface Props {
    contestId: number;
    remainingCodes: number | undefined;
    onCreated?: (contenders: Contender[]) => void;
  }

  const { contestId, remainingCodes, onCreated }: Props = $props();

  let dialog: WaDialog | undefined = $state();
  let printDialog: WaDialog | undefined = $state();
  let numberInput: WaNumberInput | undefined = $state();
  let printUrl: string | undefined = $state(undefined);

  const createContenders = createContendersMutation(contestId);

  export const open = () => {
    if (dialog) {
      dialog.open = true;
    }
  };

  const closeDialog = () => {
    if (dialog) {
      dialog.open = false;
    }
  };

  const closePrintDialog = () => {
    if (printDialog) {
      printDialog.open = false;
    }
  };

  const handlePrintNewTickets = () => {
    closePrintDialog();

    if (printUrl) {
      window.open(printUrl, "_blank");
    }
  };

  const handleCreate = () => {
    if (numberInput) {
      const args: CreateContendersArguments = {
        number: Number(numberInput.value),
      };

      createContenders.mutate(args, {
        onSuccess: (newContenders) => {
          closeDialog();

          if (newContenders.length > 0) {
            const ids = newContenders.map(({ id }) => id);
            const from = Math.min(...ids);
            const to = Math.max(...ids);
            printUrl = `/admin/contests/${contestId}/tickets/print?from=${from}&to=${to}`;

            onCreated?.(newContenders);

            if (printDialog) {
              printDialog.open = true;
            }
          }
        },
        onError: () => toastError("Failed to create tickets."),
      });
    }
  };
</script>

<wa-dialog bind:this={dialog} label="Create tickets">
  <div class="dialog-content">
    <wa-callout variant="neutral">
      <wa-icon slot="icon" name="circle-exclamation"></wa-icon>
      You have {remainingCodes} tickets remaining out of your maximum allotted
      {maxTickets}.
    </wa-callout>

    <wa-number-input
      bind:this={numberInput}
      name="number"
      {@attach value(Math.min(100, remainingCodes ?? 0))}
      min="1"
      max={remainingCodes}
      label="Number of tickets to create"
    ></wa-number-input>
  </div>

  <wa-button slot="footer" appearance="plain" onclick={closeDialog}
    >Cancel</wa-button
  >
  <wa-button
    slot="footer"
    size="small"
    variant="neutral"
    appearance="accent"
    loading={createContenders.isPending}
    onclick={handleCreate}
  >
    Create
  </wa-button>
</wa-dialog>

<wa-dialog bind:this={printDialog} label="Print tickets">
  Do you want to print the tickets that you just created?
  <wa-button slot="footer" appearance="plain" onclick={closePrintDialog}
    >Later</wa-button
  >
  <wa-button slot="footer" variant="neutral" onclick={handlePrintNewTickets}>
    <wa-icon slot="start" name="print"></wa-icon>
    Print now
  </wa-button>
</wa-dialog>

<style>
  .dialog-content {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-s);
  }
</style>
