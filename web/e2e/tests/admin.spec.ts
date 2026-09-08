import { expect, test } from "../fixtures";

test("create a competition, register a contender, enter results and draw a raffle winner", async ({
  page,
}) => {
  const username = process.env.ADMIN_USERNAME;
  const password = process.env.ADMIN_PASSWORD;

  if (!username || !password) {
    throw new Error(
      "Set ADMIN_USERNAME and ADMIN_PASSWORD to run the admin test.",
    );
  }

  let contestUrl: string;
  let contenderUrl: string;

  await test.step("Sign in through Cognito", async () => {
    await page.goto("/admin");
    await page.getByRole("button", { name: "Sign in", exact: true }).click();
    await page.waitForURL(
      "https://clmb.auth.eu-west-1.amazoncognito.com/login?**",
    );
    await page.locator('input[name="username"]:visible').fill(username);
    await page.locator('input[name="password"]:visible').fill(password);
    await page.getByRole("button", { name: "Sign in", exact: true }).click();
    await page.waitForURL(/\/admin\/?$/);
    await expect(
      page.getByRole("button", { name: "Create new competition" }),
    ).toBeVisible();
  });

  await test.step("Create the competition", async () => {
    await page.getByRole("button", { name: "Create new competition" }).click();

    await page
      .getByRole("textbox", { name: "Name *", exact: true })
      .fill("World Testing Championships");
    await page
      .getByRole("textbox", { name: "Description" })
      .fill("A friendly bouldering competition");
    await page
      .getByRole("textbox", { name: "Location" })
      .fill("Test climbing gym");
    await page.getByRole("combobox", { name: "Country" }).click();
    await page.getByRole("option", { name: "Sweden", exact: true }).click();
    await page.getByRole("spinbutton", { name: "Grace period" }).fill("5");
    await page.getByRole("radio", { name: "30 days", exact: true }).check();
    await page.getByRole("button", { name: "Add general info" }).click();
    await page
      .getByLabel("General info", { exact: true })
      .fill("Check in at reception before climbing.");
    await page.getByRole("button", { name: "Create", exact: true }).click();

    await page.waitForURL(/\/admin\/contests\/\d+$/);
    contestUrl = page.url();

    await expect(
      page.getByRole("heading", {
        name: "World Testing Championships",
        exact: true,
      }),
    ).toBeVisible();
    await expect(
      page.getByText("A friendly bouldering competition", { exact: true }),
    ).toBeVisible();
    await expect(page.getByText("Test climbing gym, Sweden")).toBeVisible();
  });

  await test.step("Set rules and categories", async () => {
    const finalists = page.locator("form").filter({
      has: page.getByRole("heading", { name: "Finalists", exact: true }),
    });
    await finalists.getByRole("checkbox").check({ force: true });
    await expect(finalists.getByText("Saved", { exact: true })).toBeVisible();

    const start = new Date(Date.now() - 60 * 60 * 1_000)
      .toISOString()
      .slice(0, 16);
    const end = new Date(Date.now() + 3 * 60 * 60 * 1_000)
      .toISOString()
      .slice(0, 16);

    for (const category of ["Open", "Recreational"]) {
      await page.getByRole("button", { name: "Create category" }).click();
      await page
        .getByRole("textbox", { name: "Name *", exact: true })
        .fill(category);
      await page.getByRole("textbox", { name: "Description" }).fill("All ages");
      await page.getByLabel("Start time", { exact: true }).fill(start);
      await page.getByLabel("End time", { exact: true }).fill(end);
      await page.getByRole("button", { name: "Create", exact: true }).click();
      await expect(
        page.getByRole("cell", { name: category, exact: true }),
      ).toBeVisible();
    }

    await page.reload();
    await expect(finalists.getByRole("checkbox")).toBeChecked();
    await expect(
      finalists.getByRole("spinbutton", { name: "Finalists" }),
    ).toHaveValue("7");
  });

  await test.step("Create problems and tickets", async () => {
    for (const number of [1, 2]) {
      await page
        .getByRole("button", { name: "Create problem", exact: true })
        .click();
      await page
        .getByRole("spinbutton", { name: "Number" })
        .fill(String(number));
      await page.getByRole("button", { name: "Primary hold color" }).click();
      await page
        .getByRole("menuitem", { name: "Select color #dc3146", exact: true })
        .click();
      await page
        .getByRole("spinbutton", { name: "Points top" })
        .fill(String(number * 100));
      await page.getByRole("spinbutton", { name: "Flash bonus" }).fill("10");
      await page
        .getByRole("switch", { name: "Enable zone Z1" })
        .check({ force: true });
      await page.getByRole("spinbutton", { name: "Points Z1" }).fill("10");
      await page
        .getByRole("switch", { name: "Enable zone Z2" })
        .check({ force: true });
      await page.getByRole("spinbutton", { name: "Points Z2" }).fill("20");
      await page.getByRole("button", { name: "Create", exact: true }).click();
      await expect(
        page.getByRole("cell", { name: `#${number}`, exact: true }),
      ).toBeVisible();
    }

    await page
      .getByRole("button", { name: "Create tickets", exact: true })
      .click();
    await page
      .getByRole("spinbutton", { name: "Number of tickets to create" })
      .fill("1");
    await page.getByRole("button", { name: "Create", exact: true }).click();
    await page.getByRole("button", { name: "Later" }).click();
    await page.getByRole("link", { name: "View and print tickets" }).click();

    const ticket = page
      .getByRole("row")
      .filter({ has: page.getByText("Unused", { exact: true }) });
    await expect(ticket).toHaveCount(1);
    await ticket.getByRole("link").click();
    await page.waitForURL(/\/admin\/contenders\/\d+$/);
    contenderUrl = page.url();
  });

  await test.step("Register the contender and enter results", async () => {
    const popup = page.waitForEvent("popup");
    await page.getByRole("link", { name: "Open scorecard" }).click();
    const scorecard = await popup;

    await scorecard
      .getByRole("textbox", { name: "Name *" })
      .fill("Dwight Schrute");
    await scorecard.getByRole("combobox", { name: "Category *" }).click();
    await scorecard
      .getByRole("option", { name: "Open All ages", exact: true })
      .click();
    await scorecard
      .getByRole("button", { name: "Register", exact: true })
      .click();

    const firstProblem = scorecard.getByRole("region", {
      name: "Problem 1",
      exact: true,
    });
    await firstProblem
      .getByRole("button", { name: "Tick", exact: true })
      .click();
    await firstProblem
      .getByRole("button", { name: "Flash", exact: true })
      .click();
    await expect(firstProblem.getByText("+110p", { exact: true })).toBeVisible({
      timeout: 15_000,
    });

    const secondProblem = scorecard.getByRole("region", {
      name: "Problem 2",
      exact: true,
    });
    await secondProblem
      .getByRole("button", { name: "Tick", exact: true })
      .click();
    await secondProblem
      .getByRole("button", { name: "Top", exact: true })
      .click();
    await expect(
      secondProblem.getByText("+200p", { exact: true }),
    ).toBeVisible();
    await expect(scorecard.getByText("310p", { exact: true })).toBeVisible();
    await scorecard.close();
  });

  await test.step("Verify the contender and results pages", async () => {
    await page.goto(contenderUrl);
    await expect(
      page.getByRole("heading", { name: "Dwight Schrute", exact: true }),
    ).toBeVisible();
    await expect(
      page.getByText("Category Open", { exact: true }),
    ).toBeVisible();
    await expect(
      page.getByText("Placement 1st", { exact: true }),
    ).toBeVisible();
    await expect(page.getByText("Score 310", { exact: true })).toBeVisible();

    for (const number of [1, 2]) {
      const result = page.getByRole("row").filter({
        has: page.getByRole("cell", { name: `#${number}`, exact: true }),
      });
      await expect(
        result.getByRole("cell", { name: "T", exact: true }),
      ).toBeVisible();
    }

    await page.goto(contestUrl);
    await page
      .getByRole("button", { name: "View results", exact: true })
      .click();
    await expect(
      page.getByRole("heading", { name: "Results", exact: true }),
    ).toBeVisible();
    const result = page.getByRole("row").filter({
      has: page.getByRole("link", { name: "Dwight Schrute", exact: true }),
    });
    await expect(
      result.getByRole("cell", { name: "310 pts", exact: true }),
    ).toBeVisible();
    await expect(
      result.getByRole("cell", { name: "1st", exact: true }),
    ).toBeVisible();
  });

  await test.step("Finish with a raffle", async () => {
    await page.goto(contestUrl);
    await page
      .getByRole("button", { name: "Start new raffle", exact: true })
      .click();
    await page
      .getByRole("button", { name: "Draw winner at random", exact: true })
      .click();
    await expect(
      page.getByRole("cell", { name: "Dwight Schrute", exact: true }),
    ).toBeVisible();
    await expect(
      page.getByText("All eligible winners have been drawn.", { exact: true }),
    ).toBeVisible();
  });
});
