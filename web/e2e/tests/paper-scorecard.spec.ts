import { expect, test, type Page } from "@playwright/test";
import type {
  Category,
  Competition,
  Problem,
} from "../../packages/lib/src/models/generated";

const competition: Competition = {
  id: 1,
  ownership: { organizerId: 1 },
  name: "Emergency test competition",
  country: "SE",
  qualifyingProblems: 10,
  finalists: 6,
  usePoints: true,
  pooledPoints: true,
  maxAttempts: 99,
  pointDeduction: 1,
  gracePeriod: 0,
  nameRetentionTime: 0,
  created: new Date("2026-11-22T10:00:00Z"),
  registeredCompetitors: 100,
};

const categories: Category[] = [
  {
    id: 1,
    competitionId: 1,
    name: "Women",
    timeBegin: new Date("2026-11-22T10:00:00Z"),
    timeEnd: new Date("2026-11-22T14:00:00Z"),
  },
  {
    id: 2,
    competitionId: 1,
    name: "Men",
    timeBegin: new Date("2026-11-22T10:00:00Z"),
    timeEnd: new Date("2026-11-22T14:00:00Z"),
  },
];

const makeProblems = (count: number): Problem[] =>
  Array.from({ length: count }, (_, index) => ({
    id: index + 1,
    competitionId: 1,
    number: index + 1,
    holdColorPrimary: "red",
    zone1Enabled: index % 3 !== 0,
    zone2Enabled: index % 3 === 2,
    pointsTop: 100,
  }));

const mockCompetition = async (
  page: Page,
  problems = makeProblems(60),
  categoryList = categories,
  competitionData = competition,
) => {
  await page.addInitScript(() => {
    const payload = btoa(JSON.stringify({ exp: Date.now() / 1000 + 86_400 }));
    localStorage.setItem("access_token", `test.${payload}.test`);
    localStorage.setItem("refresh_token", "test");
  });

  await page.route("**/api/**", async (route) => {
    const pathname = new URL(route.request().url()).pathname;
    const responses: Record<string, unknown> = {
      "/api/users/self": {
        id: 1,
        username: "organizer",
        admin: false,
        organizers: [],
      },
      "/api/competitions/1": competitionData,
      "/api/competitions/1/categories": categoryList,
      "/api/competitions/1/problems": problems,
    };

    await route.fulfill({
      status: pathname in responses ? 200 : 404,
      json: responses[pathname] ?? {},
    });
  });
};

const openScorecard = async (page: Page) => {
  await page.goto("/admin/competitions/1/scorecard/print");
  await expect(
    page.getByRole("button", { name: "Print", exact: true }),
  ).toBeVisible();
};

test("offers paper scorecards from the competition's advanced section", async ({
  page,
}) => {
  await mockCompetition(page);
  await page.goto("/admin/competitions/1");
  const link = page.getByRole("link", { name: "Print paper scorecard" });
  await expect(link).toHaveAttribute(
    "href",
    "/admin/competitions/1/scorecard/print",
  );
  await expect(link).toHaveAttribute("target", "_blank");
});

test("prints a blank 60-problem scorecard in three columns on one A4 page", async ({
  page,
}) => {
  const problems = makeProblems(60).reverse();
  problems[0].number = 100;
  const requests: string[] = [];
  page.on("request", (request) => {
    if (request.url().includes("/api/")) {
      requests.push(new URL(request.url()).pathname);
    }
  });
  await mockCompetition(page, problems);
  await openScorecard(page);

  const sheet = page.getByRole("region", {
    name: "Scorecard page 1",
    exact: true,
  });
  await expect(page.locator(".sheet")).toHaveCount(1);
  await expect(sheet.getByRole("img", { name: "ClimbLive" })).toBeVisible();
  await expect(
    sheet.getByRole("heading", { name: competition.name }),
  ).toBeVisible();
  await expect(sheet.locator(".name")).toHaveText("Name");
  await expect(sheet.locator(".category")).toHaveText(["Women", "Men"]);
  await expect(sheet.getByRole("table")).toHaveCount(3);
  await expect(sheet.locator("tbody tr")).toHaveCount(60);
  await expect(sheet.locator("tbody th")).toHaveText([
    ...Array.from({ length: 59 }, (_, index) => String(index + 1)),
    "100",
  ]);
  await expect(sheet).toContainText("Tick only your highest achievement");
  await expect(sheet).toContainText("Page 1 of 1");
  expect(requests).not.toContain("/api/competitions/1/competitors");
  expect(requests.some((path) => /ticks|points|results/.test(path))).toBe(
    false,
  );

  await page.evaluate(() => {
    window.print = () => {
      document.documentElement.dataset.printed = "true";
    };
  });
  await page.getByRole("button", { name: "Print", exact: true }).click();
  await expect(page.locator("html")).toHaveAttribute("data-printed", "true");

  await page.emulateMedia({ media: "print" });
  await expect(page.locator(".print-prompt")).toBeHidden();
  await expect(page.locator("header").first()).toBeHidden();
  await expect(sheet).toBeVisible();
});

test("only provides checkboxes for enabled zones", async ({ page }) => {
  await mockCompetition(page, makeProblems(9));
  await openScorecard(page);

  const rows = page.locator("tbody tr");
  await expect(rows.nth(0).locator(".checkbox")).toHaveCount(2);
  await expect(
    rows.nth(0).getByRole("cell", { name: "Not available", exact: true }),
  ).toHaveCount(2);
  await expect(rows.nth(1).locator(".checkbox")).toHaveCount(3);
  await expect(rows.nth(2).locator(".checkbox")).toHaveCount(4);
  await expect(
    rows.getByRole("cell", { name: "Top", exact: true }),
  ).toHaveCount(9);
  await expect(
    rows.getByRole("cell", { name: "Flash", exact: true }),
  ).toHaveCount(9);
});

test("omits zone columns when no problems use zones", async ({ page }) => {
  const problems = makeProblems(60).map((problem) => ({
    ...problem,
    zone1Enabled: false,
    zone2Enabled: false,
  }));
  await mockCompetition(page, problems);
  await openScorecard(page);

  await expect(page.getByRole("columnheader", { name: /^Z[12]$/ })).toHaveCount(
    0,
  );
  await expect(page.locator("tbody .checkbox")).toHaveCount(120);
});

test("overflows onto numbered pages with repeated identification and no lost problems", async ({
  page,
}) => {
  await mockCompetition(page, makeProblems(100));
  await openScorecard(page);

  await expect(page.locator(".sheet")).toHaveCount(2);
  await expect(page.locator("tbody th")).toHaveText(
    Array.from({ length: 100 }, (_, index) => String(index + 1)),
  );

  for (let index = 0; index < 2; index++) {
    const sheet = page.locator(".sheet").nth(index);
    await expect(sheet.getByRole("img", { name: "ClimbLive" })).toBeVisible();
    await expect(sheet.locator(".name")).toHaveText("Name");
    await expect(sheet.locator(".category")).toHaveText(["Women", "Men"]);
    await expect(sheet).toContainText(`Page ${index + 1} of 2`);
  }

  await page.emulateMedia({ media: "print" });
  const fits = await page.locator(".sheet").evaluateAll((sheets) =>
    sheets.every((sheet) => {
      const footer = sheet.querySelector("footer")!.getBoundingClientRect();
      return [...sheet.querySelectorAll("table")].every(
        (table) => table.getBoundingClientRect().bottom <= footer.top,
      );
    }),
  );
  expect(fits).toBe(true);
  await expect(page.locator(".sheet").last()).toHaveCSS("break-after", "auto");
});

for (const [problemCount, pageCount] of [
  [60, 1],
  [100, 2],
]) {
  test(`prints ${problemCount} problems on exactly ${pageCount} physical A4 pages`, async ({
    page,
    browserName,
  }) => {
    test.skip(browserName !== "chromium", "PDF printing requires Chromium");
    await mockCompetition(page, makeProblems(problemCount));
    await openScorecard(page);
    await expect(page.locator(".sheet")).toHaveCount(pageCount);

    const pdf = await page.pdf({ preferCSSPageSize: true });
    const contents = pdf.toString("latin1");
    expect(contents.match(/\/Type \/Page\b/g)).toHaveLength(pageCount);
    const mediaBox = contents.match(/\/MediaBox \[0 0 ([\d.]+) ([\d.]+)\]/);
    expect(mediaBox).not.toBeNull();
    expect(Number(mediaBox![1])).toBeCloseTo(595, 0);
    expect(Number(mediaBox![2])).toBeCloseTo(842, 0);
    await test.info().attach("paper-scorecard.pdf", {
      body: pdf,
      contentType: "application/pdf",
    });
  });
}

test("allows space for long competition names and many categories", async ({
  page,
}) => {
  const categoryList = Array.from({ length: 16 }, (_, index) => ({
    ...categories[0],
    id: index + 1,
    name: `Competition category ${index + 1}`,
  }));
  await mockCompetition(page, makeProblems(100), categoryList, {
    ...competition,
    name: "A very long competition name that needs multiple lines on the paper scorecard for this emergency fallback competition",
  });
  await openScorecard(page);
  await expect(page.locator("tbody tr")).toHaveCount(100);
  await page.emulateMedia({ media: "print" });

  await expect
    .poll(() =>
      page.locator(".sheet").evaluateAll((sheets) =>
        sheets.every((sheet) => {
          const footer = sheet.querySelector("footer")!.getBoundingClientRect();
          return [...sheet.querySelectorAll("table")].every(
            (table) => table.getBoundingClientRect().bottom <= footer.top,
          );
        }),
      ),
    )
    .toBe(true);
});

test("provides a handwritten category field when no categories exist", async ({
  page,
}) => {
  await mockCompetition(page, makeProblems(1), []);
  await openScorecard(page);
  await expect(page.locator(".categories .writing-line")).toBeVisible();
  await expect(page.locator("tbody tr")).toHaveCount(1);
});

test("does not offer printing before problems finish loading", async ({
  page,
}) => {
  await mockCompetition(page);
  let release: () => void = () => {};
  const ready = new Promise<void>((resolve) => {
    release = resolve;
  });
  await page.route("**/api/competitions/1/problems", async (route) => {
    await ready;
    await route.fulfill({ json: makeProblems(60) });
  });
  await page.goto("/admin/competitions/1/scorecard/print");
  await expect(
    page.getByRole("button", { name: "Print", exact: true }),
  ).toHaveCount(0);
  release();
  await expect(
    page.getByRole("button", { name: "Print", exact: true }),
  ).toBeVisible();
});

test("does not print a competition without problems", async ({ page }) => {
  await mockCompetition(page, []);
  await page.goto("/admin/competitions/1/scorecard/print");
  await expect(
    page.getByRole("heading", { name: "No problems to print" }),
  ).toBeVisible();
  await expect(
    page.getByRole("button", { name: "Print", exact: true }),
  ).toHaveCount(0);
  await expect(page.locator(".sheet")).toHaveCount(0);
});

test("shows a retry action instead of an incomplete scorecard on a loading error", async ({
  page,
}) => {
  await mockCompetition(page);
  await page.route("**/api/competitions/1/categories", (route) =>
    route.fulfill({ status: 503, json: {} }),
  );
  await page.goto("/admin/competitions/1/scorecard/print");
  await expect(
    page.getByRole("heading", { name: "Unable to load scorecard" }),
  ).toBeVisible({ timeout: 20_000 });
  await expect(
    page.getByRole("button", { name: "Print", exact: true }),
  ).toHaveCount(0);
  await page.unroute("**/api/competitions/1/categories");
  await page.getByRole("button", { name: "Try again" }).click();
  await expect(
    page.getByRole("button", { name: "Print", exact: true }),
  ).toBeVisible();
});
