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

class CompatElement extends HTMLElement {
  constructor() {
    super();
  }
}

window.customElements.define("cl-compat", CompatElement);

const element: CompatElement = document.createElement("cl-compat");

if (element.attachInternals === undefined) {
  alert("ElementInternals is not supported");
} else {
  alert("ElementInternals supported");

  const internals = element.attachInternals();

  if (internals.states === undefined) {
    alert("CustomStateSet is not supported");
  } else {
    alert("CustomStateSet supported");
  }
}

const app = mount(App, {
  target: document.body,
});

export default app;
