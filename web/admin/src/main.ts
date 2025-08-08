import App from "@/App.svelte";
import "@/main.css";
import {
  prefersDarkColorScheme,
  updateTheme,
  watchColorSchemeChanges,
} from "@climblive/lib/utils";
import { mount } from "svelte";

watchColorSchemeChanges((prefersDarkColorScheme) =>
  updateTheme(prefersDarkColorScheme),
);
updateTheme(prefersDarkColorScheme());

const app = mount(App, {
  target: document.body,
});

export default app;
