import "@awesome.me/webawesome/dist/components/button/button.js";
import {
  prefersDarkColorScheme,
  updateTheme,
  watchColorSchemeChanges,
} from "@climblive/lib/utils";
import "../styles.css";

document.getElementById("current-year")!.textContent = new Date()
  .getFullYear()
  .toString();

watchColorSchemeChanges((prefersDarkColorScheme) =>
  updateTheme(prefersDarkColorScheme),
);
updateTheme(prefersDarkColorScheme());
