import "@awesome.me/webawesome/dist/components/alert/alert.js";
import "@awesome.me/webawesome/dist/components/icon/icon.js";

export const toastError = (message: string, duration = 5000) => {
  const alert = Object.assign(document.createElement("wa-alert"), {
    variant: "danger",
    closable: true,
    duration: duration,
    innerHTML: `
        <wa-icon name="exclamation-octagon" slot="icon"></wa-icon>
        <strong>An unexpected error occurred</strong><br />
        ${message}
      `,
  });

  document.body.append(alert);
  return alert.toast();
};
