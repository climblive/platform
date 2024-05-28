import "@shoelace-style/shoelace/dist/components/alert/alert.js";
import "@shoelace-style/shoelace/dist/components/icon/icon.js";

export const toastError = (message: string, duration = 5000) => {
  const alert = Object.assign(document.createElement("sl-alert"), {
    variant: "danger",
    closable: true,
    duration: duration,
    innerHTML: `
        <sl-icon name="exclamation-octagon" slot="icon"></sl-icon>
        <strong>An unexpected error occurred</strong><br />
        ${message}
      `,
  });

  document.body.append(alert);
  return alert.toast();
};
