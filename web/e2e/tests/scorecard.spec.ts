import { expect, test } from '@playwright/test';
import {
  MariaDbContainer,
  StartedMariaDbContainer
} from "@testcontainers/mariadb";
import { readFile } from 'fs/promises';
import { Connection, createConnection } from "mariadb";
import { GenericContainer, Network, StartedTestContainer, Wait } from "testcontainers";

let dbConnection: Connection | undefined;
let startedDbContainer: StartedMariaDbContainer | undefined;
let startedApiContainer: StartedTestContainer | undefined;
let startedWebContainer: StartedTestContainer | undefined;

test.describe.configure({ mode: 'serial' });

test.beforeAll(async () => {
  const network = await new Network().start();

  startedDbContainer = await new MariaDbContainer()
    .withUsername("climblive")
    .withUserPassword("secretpassword")
    .withDatabase("climblive")
    .withExposedPorts(3306)
    .withNetwork(network)
    .withNetworkAliases("e2e")
    .start();

  dbConnection = await createConnection({
    host: "localhost",
    port: startedDbContainer.getMappedPort(3306),
    user: "climblive",
    password: "secretpassword",
    database: "climblive",
    multipleStatements: true,
  });

  const schema = await readFile("../../backend/database/scoreboard.sql", "utf8")
  const samples = await readFile("./samples.sql", "utf8")

  dbConnection.query(schema)
  dbConnection.query(samples)

  const apiContainer = new GenericContainer("climblive-api:latest")
    .withEnvironment({
      "DB_USERNAME": "climblive",
      "DB_PASSWORD": "secretpassword",
      "DB_HOST": "e2e",
      "DB_PORT": "3306",
      "DB_DATABASE": "climblive",
    })
    .withNetwork(network)
    .withExposedPorts({ container: 8090, host: 8090 })
    .withWaitStrategy(Wait.forLogMessage(/score engine hydration complete/))

  const webContainer = new GenericContainer("climblive-web:latest")
    .withNetwork(network)
    .withExposedPorts({ container: 80, host: 8080 })
    .withWaitStrategy(Wait.forListeningPorts())

  const startedContainers = await Promise.all([apiContainer.start(), webContainer.start()])

  startedApiContainer = startedContainers[0];
  startedWebContainer = startedContainers[1]
})

test.afterAll(async () => {
  await startedWebContainer?.stop()
  await startedApiContainer?.stop();
  await dbConnection?.end();
  await startedDbContainer?.stop()
})

test('enter contest by entering registration code', async ({ page }) => {
  await page.goto('/');

  await expect(page).toHaveTitle(/ClimbLive/);

  const pinInput = page.getByRole("textbox", { name: "Pin character 1 out of 8" })
  await pinInput.pressSequentially("abcd0002");

  await page.waitForURL('/ABCD0002/register');

  await page.getByRole("textbox", { name: "Full name *" }).pressSequentially("Dwight Schrute")
  await page.getByRole("textbox", {
    name: "Club name"
  }).pressSequentially("Scranton Climbing Club")
  const compClass = page.getByRole("combobox", { name: "Competition class *" })
  await compClass.click()
  await page.getByRole("option", { name: "Males", exact: true }).click()

  await page.getByRole("button", { name: "Register" }).click()

  await page.waitForURL('/ABCD0002');
});

test('registration code is saved in local storage', async ({ page }) => {
  await page.goto('/');

  const pinInput = page.getByRole("textbox", { name: "Pin character 1 out of 8" })
  await pinInput.pressSequentially("abcd0001");

  await page.waitForURL('/ABCD0001');
  await expect(page.getByText("Albert Einstein")).toBeVisible();

  await page.goto('/');
  await page.waitForURL('/');

  await page.getByRole("button", { name: "Enter" }).click()

  await page.waitForURL('/ABCD0001');
  await expect(page.getByText("Albert Einstein")).toBeVisible();
});

test('deep link into scorecard', async ({ page }) => {
  await page.goto('/abcd0001');

  await expect(page.getByText("Albert Einstein")).toBeVisible();
});

test('garbage session value in local storage is thrown out', async ({ page }) => {
  await page.goto('/');

  await page.evaluate(() => localStorage.setItem('session', 'bad_data'))

  await page.goto('/');

  await expect(page.getByRole("textbox", { name: "Pin character 1 out of 8" })).toHaveValue("");
});

test('edit profile', async ({ page }) => {
  await page.goto('/ABCD0003');

  await expect(page.getByText("Michael Scott")).toBeVisible()
  await expect(page.getByText("Scranton Climbing Club")).toBeVisible()
  await expect(page.getByText("Males")).toBeVisible()

  await page.getByRole("button", { name: "Edit" }).click();

  await page.waitForURL('/ABCD0003/edit');

  const nameInput = page.getByRole("textbox", { name: "Full name *" })
  await nameInput.fill("")
  await nameInput.pressSequentially("Phyllis Lapin-Vance")

  const clubNameInput = page.getByRole("textbox", {
    name: "Club name"
  })
  await clubNameInput.fill("")
  await clubNameInput.pressSequentially("Dunder Mifflin Climbing Club")

  const compClass = page.getByRole("combobox", { name: "Competition class *" })
  await compClass.click()
  await page.getByRole("option", { name: "Females", exact: true }).click()

  await page.getByRole("button", { name: "Save" }).click()

  await page.waitForURL('/ABCD0003');

  await expect(page.getByText("Phyllis Lapin-Vance")).toBeVisible()
  await expect(page.getByText("Dunder Mifflin Climbing Club")).toBeVisible()
  await expect(page.getByText("Females")).toBeVisible()
});

test('cancel edit profile', async ({ page }) => {
  await page.goto('/ABCD0001/edit');

  await page.getByRole("button", { name: "Cancel" }).click();

  await page.waitForURL('/ABCD0001');

  await expect(page.getByText("Albert Einstein")).toBeVisible()
});

test("tick and untick all problems", async ({ page }) => {
  await page.goto('/ABCD0003');

  for (let p = 1; p <= 5; p++) {
    const problem = page.getByRole("region", { name: `Problem ${p}` });
    await expect(problem).toBeVisible();

    await problem.getByRole("button", { name: "Tick" }).click();
    await problem.getByRole("button", { name: "Top" }).click();

    await expect(problem.getByText(`+${p * 100}p`)).toBeVisible();
  }

  await expect(page.getByText("1500p")).toBeVisible()
  await expect(page.getByText("1st place")).toBeVisible()

  for (let p = 1; p <= 5; p++) {
    const problem = page.getByRole("region", { name: `Problem ${p}` });
    await expect(problem).toBeVisible();

    await problem.getByRole("button", { name: "Untick" }).click();

    await expect(problem.getByText(`+${p * 100}p`)).not.toBeVisible();
  }

  await expect(page.getByText("0p", { exact: true })).toBeVisible()
  await expect(page.getByText("1st place")).toBeVisible()
})

test("flash a problem", async ({ page }) => {
  await page.goto('/ABCD0003');

  const problem = page.getByRole("region", { name: "Problem 1" });
  await expect(problem).toBeVisible();

  await problem.getByRole("button", { name: "Tick" }).click();
  await problem.getByRole("button", { name: "Flash" }).click();

  await expect(problem.getByText("+110p")).toBeVisible();

  await problem.getByRole("button", { name: "Untick" }).click();

  await expect(problem.getByText("+110p")).not.toBeVisible();
})

test("tick buttons are disabled before contest has started", async ({ page }) => {
  await page.clock.setFixedTime(new Date('2023-11-01T00:00:00'));

  await page.goto('/ABCD0001');

  const timer = page.getByRole("timer", { name: "Time until start" })
  await expect(timer).toHaveText("2 months")

  await expect(page.getByRole("button", { name: "Edit" })).toBeEnabled();

  const problem = page.getByRole("region", { name: "Problem 1" });
  await expect(problem).toBeVisible();

  await expect(problem.getByRole("button", { name: "Tick" })).toBeDisabled();
})

test.describe("screenshots", () => {
  test("registration code input", async ({ page }) => {
    await page.goto('/')

    await expect(page.getByRole("heading", { name: "Welcome" })).toBeVisible()
    await expect(page).toHaveScreenshot({ maxDiffPixels: 100 })
  });

  test("registration", async ({ page }) => {
    await page.goto('/abcd0004/register')

    await expect(page.getByRole("textbox", { name: "Full name *" })).toBeVisible();
    await expect(page).toHaveScreenshot({ maxDiffPixels: 100 })
  });

  test("edit profile", async ({ page }) => {
    await page.goto('/abcd0001/edit')

    await expect(page.getByRole("textbox", { name: "Full name *" })).toBeVisible();
    await expect(page).toHaveScreenshot({ maxDiffPixels: 100 })
  });

  test.describe("scorecard", () => {
    test("before contest has started", async ({ page }) => {
      await page.clock.setFixedTime(new Date('2023-11-01T00:00:00'));

      await page.goto('/abcd0001')

      await expect(page.getByText("Albert Einstein")).toBeVisible();
      await expect(page).toHaveScreenshot({ maxDiffPixels: 100 })
    });

    test("while contest is running", async ({ page }) => {
      await page.clock.setFixedTime(new Date('2024-01-01T00:00:00'));

      await page.goto('/abcd0001')

      await expect(page.getByText("Albert Einstein")).toBeVisible();
      await expect(page).toHaveScreenshot({ maxDiffPixels: 100 })
    });

    test("during grace period", async ({ page }) => {
      await page.clock.setFixedTime(new Date('2025-01-01T00:00:00'));

      await page.goto('/abcd0001')

      await expect(page.getByText("Albert Einstein")).toBeVisible();
      await expect(page).toHaveScreenshot({ maxDiffPixels: 100 })
    });

    test("after contest has ended", async ({ page }) => {
      await page.clock.setFixedTime(new Date('2025-01-01T00:05:00'));

      await page.goto('/abcd0001')

      await expect(page.getByText("Albert Einstein")).toBeVisible();
      await expect(page).toHaveScreenshot({ maxDiffPixels: 100 })
    });
  })
});
