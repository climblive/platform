import App from "@/App.svelte";
import "@/main.css";
import { Fallback } from "@climblive/lib/components";
import {
  checkCompat,
  prefersDarkColorScheme,
  updateTheme,
  watchColorSchemeChanges,
} from "@climblive/lib/utils";
import { mount } from "svelte";

const urlParams = new URLSearchParams(window.location.search);
if (urlParams.get("print") === null) {
  watchColorSchemeChanges((prefersDarkColorScheme) =>
    updateTheme(prefersDarkColorScheme),
  );

  updateTheme(prefersDarkColorScheme());
}

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
    },
  });
}
