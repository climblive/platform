import "@/main.css";
import {
  prefersDarkColorScheme,
  updateTheme,
  watchColorSchemeChanges,
} from "@climblive/lib/utils";
import { mount } from "svelte";
import App from "./App.svelte";
import Fallback from "./Fallback.svelte";

watchColorSchemeChanges((prefersDarkColorScheme) =>
  updateTheme(prefersDarkColorScheme),
);
updateTheme(prefersDarkColorScheme());

class CompatElement extends HTMLElement {
  constructor() {
    super();
  }
}

window.customElements.define("cl-compat", CompatElement);

const checkCompat = () => {
  const element: CompatElement = document.createElement("cl-compat");

  if (element.attachInternals === undefined) {
    return false;
  }

  const internals = element.attachInternals();

  if (internals.states === undefined) {
    return false;
  }

  return true;
};

const appComponent = checkCompat() ? App : Fallback;

const app = mount(appComponent, {
  target: document.body,
});

export default app;
