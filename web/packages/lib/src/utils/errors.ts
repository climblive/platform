import "@awesome.me/webawesome/dist/components/callout/callout.js";
import "@awesome.me/webawesome/dist/components/icon/icon.js";

export const toastError = (message: string, duration = 5000) => {
  const alert = Object.assign(document.createElement("wa-callout"), {
    variant: "danger",
    closable: true,
    duration: duration,
    innerHTML: `
        <wa-icon name="circle-exclamation" slot="icon"></wa-icon>
        <strong>An unexpected error occurred</strong><br />
        ${message}
      `,
  });

  document.getElementById("callout-stack")?.append(alert);

  setTimeout(() => {
    alert.remove();
  }, duration);
};
