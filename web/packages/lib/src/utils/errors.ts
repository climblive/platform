import "@awesome.me/webawesome/dist/components/icon/icon.js";
import WaToast from "@awesome.me/webawesome/dist/components/toast/toast.js";

export const toastError = (message: string, duration = 5000) => {
  const toast = document.querySelector("wa-toast");

  if (!toast || !(toast instanceof WaToast)) {
    return;
  }

  toast.create(
    `
        <wa-icon name="circle-exclamation" slot="icon"></wa-icon>
        <strong>An unexpected error occurred</strong><br />
        ${message}
    `,
    {
      allowHtml: true,
      variant: "brand",
      duration,
    },
  );
};
