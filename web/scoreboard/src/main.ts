import App from "@/App.svelte";
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
import NativeStyles from "./NativeStyles.svelte";

if (import.meta.env.PROD) {
  Sentry.init({
    dsn: "https://019099d850441f60cea5d465e217f768@o4509937603641344.ingest.de.sentry.io/4509937616093264",
    sendDefaultPii: false,
    environment: import.meta.env.VITE_SENTRY_ENVIRONMENT,
  });
}

watchColorSchemeChanges((prefersDarkColorScheme) =>
  updateTheme(prefersDarkColorScheme),
);
updateTheme(prefersDarkColorScheme());

const [compatible, missingFeatures] = checkCompat();

if (compatible) {
  mount(App, {
    target: document.body,
  });
} else {
  mount(Fallback, {
    target: document.body,
    props: {
      missingFeatures,
      app: App,
      styles: NativeStyles,
    },
  });
}
