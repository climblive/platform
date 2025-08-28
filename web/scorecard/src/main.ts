import "@/main.css";
import {
  prefersDarkColorScheme,
  updateTheme,
  watchColorSchemeChanges
} from "@climblive/lib/utils";
import { mount } from "svelte";
import App from "./App.svelte";
import Fallback from "./Fallback.svelte";

watchColorSchemeChanges((prefersDarkColorScheme) =>
  updateTheme(prefersDarkColorScheme),
);
updateTheme(prefersDarkColorScheme());

const [compatible, missingFeatures] = [false, ["ElementInternals", "CustomElementRegistry", "CustomStateSet"]]// checkCompat();

if (compatible) {
  mount(App, {
    target: document.body,
  })
} else {
  mount(Fallback, {
    target: document.body,
    props: {
      missingFeatures,
      app: App
    },
  });
}
