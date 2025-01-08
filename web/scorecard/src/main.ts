import App from "@/App.svelte";
import "@/main.css";
import 'animate.css';
import { mount } from "svelte";

const app = mount(App, {
  target: document.body,
});

export default app;
