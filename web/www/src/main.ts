import "@awesome.me/webawesome/dist/components/badge/badge.js";
import "@awesome.me/webawesome/dist/components/button/button.js";
import "@awesome.me/webawesome/dist/components/card/card.js";
import "@awesome.me/webawesome/dist/components/carousel/carousel.js";
import "@awesome.me/webawesome/dist/components/carousel-item/carousel-item.js";
import "@awesome.me/webawesome/dist/components/number-ticker/number-ticker.js";
import type WaNumberTicker from "@awesome.me/webawesome/dist/components/number-ticker/number-ticker.js";
import {
  prefersDarkColorScheme,
  updateTheme,
  watchColorSchemeChanges,
} from "@climblive/lib/utils";
import type { GlobalStats } from "@climblive/lib/models";
import { getApiBaseUrl } from "./api";
import "../styles.css";

const currentYear = document.getElementById("current-year");
if (currentYear) {
  currentYear.textContent = new Date().getFullYear().toString();
}

watchColorSchemeChanges((prefersDarkColorScheme) =>
  updateTheme(prefersDarkColorScheme),
);
updateTheme(prefersDarkColorScheme());

const loadGlobalStats = async () => {
  const section = document.getElementById("global-stats");
  const ticks = document.querySelector<WaNumberTicker>("#stats-ticks");
  const registrations = document.querySelector<WaNumberTicker>(
    "#stats-registrations",
  );
  const competitions = document.querySelector<WaNumberTicker>(
    "#stats-competitions",
  );
  if (!section || !ticks || !registrations || !competitions) {
    return;
  }

  const response = await fetch(`${getApiBaseUrl()}/stats/global`);
  if (!response.ok) {
    return;
  }

  const stats: GlobalStats = await response.json();
  ticks.value = stats.ticks;
  registrations.value = stats.competitors;
  competitions.value = stats.competitions;
  section.hidden = false;
};

void loadGlobalStats().catch(() => undefined);

const compactCarouselLayout = window.matchMedia("(max-width: 768px)");

const updateCardSizes = () => {
  document.querySelectorAll("wa-carousel").forEach((carousel) => {
    const cards = carousel.querySelectorAll<HTMLElement>(
      "wa-carousel-item > wa-card",
    );

    if (cards.length === 0) {
      return;
    }

    cards.forEach((card) => {
      card.style.height = "auto";
    });

    const height = Math.ceil(
      Math.max(
        ...Array.from(cards, (card) => card.getBoundingClientRect().height),
      ),
    );

    cards.forEach((card) => {
      card.style.height = `${height}px`;
    });
  });
};

let cardSizingFrame = 0;

const updateCarouselLayout = () => {
  document.querySelectorAll("wa-carousel").forEach((carousel) => {
    carousel.setAttribute(
      "slides-per-page",
      compactCarouselLayout.matches ? "1" : "3",
    );
    carousel.toggleAttribute("navigation", !compactCarouselLayout.matches);
  });

  cancelAnimationFrame(cardSizingFrame);
  cardSizingFrame = requestAnimationFrame(updateCardSizes);
};

window.addEventListener("resize", updateCarouselLayout);
document.fonts.ready.then(updateCarouselLayout);
updateCarouselLayout();
