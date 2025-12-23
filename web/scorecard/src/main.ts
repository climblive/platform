import "@/main.css";
import { Fallback } from "@climblive/lib/components";
import {
  checkCompat,
  prefersDarkColorScheme,
  updateTheme,
  watchColorSchemeChanges,
} from "@climblive/lib/utils";
import * as Sentry from "@sentry/svelte";
import { mount } from "svelte";
import App from "./App.svelte";
import FailsafeApp from "./FailsafeApp.svelte";
import TryFailsafe from "./TryFailsafe.svelte";

if (import.meta.env.PROD) {
  Sentry.init({
    dsn: "https://019099d850441f60cea5d465e217f768@o4509937603641344.ingest.de.sentry.io/4509937616093264",
    sendDefaultPii: false,
    environment: import.meta.env.VITE_SENTRY_ENVIRONMENT,
    integrations: [Sentry.captureConsoleIntegration({ levels: ["error"] })],
  });
}

if (location.pathname.startsWith("/failsafe")) {
  mount(FailsafeApp, {
    target: document.body,
  });
} else {
  watchColorSchemeChanges((prefersDarkColorScheme) =>
    updateTheme(prefersDarkColorScheme),
  );
  updateTheme(prefersDarkColorScheme());

  const [compatible, missingFeatures] = checkCompat();

  const ignoreCompat = sessionStorage.getItem("compat") === "ignore";

  if (!compatible && !ignoreCompat) {
    console.warn(
      "Incompatible browser detected. Missing features:",
      missingFeatures,
    );
  }

  if (missingFeatures.includes("ElementInternals")) {
    await import("element-internals-polyfill");
  }

  if (compatible || ignoreCompat) {
    mount(App, {
      target: document.body,
    });
  } else {
    mount(Fallback, {
      target: document.body,
      props: {
        missingFeatures,
        alternative: TryFailsafe,
      },
    });
  }
}
