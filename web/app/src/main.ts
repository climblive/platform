import "normalize.css/normalize.css";
import "@shoelace-style/shoelace/dist/themes/light.css";
import "@climblive/shared/theme.css";
import "./main.css";
import App from "./App.svelte";

const app = new App({
  target: document.body,
});

export default app;
