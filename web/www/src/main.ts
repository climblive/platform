import "@awesome.me/webawesome/dist/components/badge/badge.js";
import "@awesome.me/webawesome/dist/components/button/button.js";
import "@awesome.me/webawesome/dist/components/card/card.js";
import "@awesome.me/webawesome/dist/components/carousel/carousel.js";
import "@awesome.me/webawesome/dist/components/carousel-item/carousel-item.js";
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

const compactCarouselLayout = window.matchMedia("(max-width: 768px)");
const updateCarouselSlides = () => {
  document.querySelectorAll("wa-carousel").forEach((carousel) => {
    carousel.setAttribute(
      "slides-per-page",
      compactCarouselLayout.matches ? "1" : "3",
    );
    carousel.toggleAttribute("navigation", !compactCarouselLayout.matches);
  });
};

compactCarouselLayout.addEventListener("change", updateCarouselSlides);
updateCarouselSlides();
