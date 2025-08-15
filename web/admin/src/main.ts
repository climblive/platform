import App from "@/App.svelte";
import "@/main.css";
import {
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

const app = mount(App, {
  target: document.body,
});

export default app;
