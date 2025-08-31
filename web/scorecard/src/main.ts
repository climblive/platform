import "@/main.css";
import { Fallback } from "@climblive/lib/components";
import {
  checkCompat,
  prefersDarkColorScheme,
  updateTheme,
  watchColorSchemeChanges,
} from "@climblive/lib/utils";
import { mount } from "svelte";
import App from "./App.svelte";
import NativeStyles from "./NativeStyles.svelte";

import * as Sentry from "@sentry/svelte";

Sentry.init({
  dsn: "https://019099d850441f60cea5d465e217f768@o4509937603641344.ingest.de.sentry.io/4509937616093264",
  sendDefaultPii: false,
});

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
