import "@awesome.me/webawesome/dist/components/icon/icon.js";
import WaToast from "@awesome.me/webawesome/dist/components/toast/toast.js";

export const toastUnexpectedError = (message: string) => {
  toast({
    title: "An unexpected error occurred",
    message,
    icon: "circle-exclamation",
    variant: "danger",
    duration: 5000,
  });
};

type ToastOptions = {
  title: string;
  message: string;
  icon: string;
  variant?: "neutral" | "brand" | "success" | "danger" | "warning";
  duration?: number;
};

export const toast = (options: ToastOptions) => {
  const toast = document.querySelector("wa-toast");

  if (!toast || !(toast instanceof WaToast)) {
    return;
  }

  toast.create(
    `
        <wa-icon name="${options.icon}" slot="icon"></wa-icon>
        <strong>${options.title}</strong><br />
        ${options.message}
    `,
    {
      allowHtml: true,
      variant: options.variant ?? "neutral",
      duration: options.duration ?? 5000,
      size: "s",
    },
  );
};
