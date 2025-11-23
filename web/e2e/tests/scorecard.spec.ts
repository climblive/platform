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

  startedDbContainer = await new MariaDbContainer("mariadb:11.4")
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

  const schema = await readFile("../../backend/database/climblive.sql", "utf8")
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
    .withWaitStrategy(Wait.forLogMessage(/score engine started/))

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

  const codeInput = page.getByRole("textbox", { name: "Registration code *" })
  await codeInput.pressSequentially("abcd0002");

  await page.getByRole("button", { name: "Enter" }).click()

  await page.waitForURL('/ABCD0002/register');

  await page.getByRole("textbox", { name: "Full name *" }).pressSequentially("Dwight Schrute")
  const compClass = page.getByRole("combobox", { name: "Competition class *" })
  await compClass.click()
  await page.getByRole("option", { name: "Males", exact: true }).click()

  await page.getByRole("button", { name: "Register" }).click()

  await page.waitForURL('/ABCD0002');
});

test('registration code is saved in local storage for 12 hours', async ({ page }) => {
  await page.clock.install({ time: new Date() });
  await page.goto('/');

  const codeInput = page.getByRole("textbox", { name: "Registration code *" })
  await codeInput.pressSequentially("abcd0001");

  await page.getByRole("button", { name: "Enter" }).click()

  await page.waitForURL('/ABCD0001');
  await expect(page.getByText("Albert Einstein")).toBeVisible();

  await page.goto('/');
  await page.waitForURL('/');

  const region = await page.getByRole("region", { name: "Saved session ABCD0001" });

  await region.getByRole("button", { name: "Restore" }).click()

  await page.waitForURL('/ABCD0001');
  await expect(page.getByText("Albert Einstein")).toBeVisible();

  await page.clock.fastForward("12:00:00");

  await page.goto('/');
  await page.waitForURL('/');

  await expect(page.getByRole("button", { name: "Restore" })).not.toBeVisible();
});

test('the three most recently used registration codes can be restored', async ({ page }) => {
  await page.clock.install({ time: new Date() });

  for (const code of ["ABCD0001", "ABCD0002", "ABCD0003", "ABCD0004"]) {
    await page.goto(`/${code}`);
    await page.waitForURL(`/${code}`);

    await page.clock.fastForward("00:00:01");
  }

  await page.goto('/');
  await page.waitForURL('/');

  await expect(page.getByRole("region", { name: "Saved session ABCD0001" })).not.toBeVisible();

  await expect(page.getByRole("region", { name: "Saved session ABCD0002" })).toBeVisible();
  await expect(page.getByRole("region", { name: "Saved session ABCD0003" })).toBeVisible();
  await expect(page.getByRole("region", { name: "Saved session ABCD0004" })).toBeVisible();
});

test('deep link into scorecard', async ({ page }) => {
  await page.goto('/abcd0001');

  await expect(page.getByText("Albert Einstein")).toBeVisible();
});

test('garbage session value in local storage is thrown out', async ({ page }) => {
  await page.goto('/');

  await page.evaluate(() => localStorage.setItem('sessions', 'bad_data'))

  await page.goto('/');

  await expect(page.getByRole("textbox", { name: "Registration code *" })).toHaveValue("");
});

test('edit profile', async ({ page }) => {
  await page.goto('/ABCD0003');

  await expect(page.getByText("Michael Scott")).toBeVisible()
  await expect(page.getByText("Males", { exact: true })).toBeVisible()

  await page.getByRole("button", { name: "Edit" }).click({ force: true });

  await page.waitForURL('/ABCD0003/edit');

  const nameInput = page.getByRole("textbox", { name: "Full name *" })
  await nameInput.fill("")
  await nameInput.pressSequentially("Phyllis Lapin-Vance")

  const compClass = page.getByRole("combobox", { name: "Competition class *" })
  await compClass.click()
  await page.getByRole("option", { name: "Females", exact: true }).click()

  await page.getByRole("button", { name: "Save" }).click()

  await page.waitForURL('/ABCD0003');

  await expect(page.getByText("Phyllis Lapin-Vance")).toBeVisible()
  await expect(page.getByText("Females", { exact: true })).toBeVisible()
});

test('withdraw from finals and reenter', async ({ page }) => {
  await page.goto('/ABCD0003/edit');

  await expect(page.getByRole("switch", { name: "Opt out of finals" })).not.toBeChecked();
  await page.getByRole("switch", { name: "Opt out of finals" }).check({ force: true });

  await page.getByRole("button", { name: "Save" }).click()

  await page.waitForURL('/ABCD0003');

  await page.getByRole("button", { name: "Edit" }).click({ force: true });

  await expect(page.getByRole("switch", { name: "Opt out of finals" })).toBeChecked();
  await page.getByRole("switch", { name: "Opt out of finals" }).uncheck({ force: true });

  await page.getByRole("button", { name: "Save" }).click()

  await page.waitForURL('/ABCD0003');

  await page.getByRole("button", { name: "Edit" }).click({ force: true });

  await expect(page.getByRole("switch", { name: "Opt out of finals" })).not.toBeChecked();
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
  await expect(page.getByText("1st")).toBeVisible()

  for (let p = 1; p <= 5; p++) {
    const problem = page.getByRole("region", { name: `Problem ${p}` });
    await expect(problem).toBeVisible();

    await problem.getByRole("button", { name: "Untick" }).click();

    await expect(problem.getByText(`+${p * 100}p`)).not.toBeVisible();
  }

  await expect(page.getByText("0p", { exact: true })).toBeVisible()
  await expect(page.getByText("1st")).toBeVisible()
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

test("info tab", async ({ page }) => {
  await page.goto('/ABCD0001');

  await page.getByRole("tab", { name: "Info" }).click();

  await expect(page.getByText("Name World Testing Championships")).toBeVisible();
  await expect(page.getByText("Description The world's number one competition for testing")).toBeVisible();
  await expect(page.getByText("Location On the web")).toBeVisible();
  await expect(page.getByText("Classes Males, Females")).toBeVisible();
  await expect(page.getByText("Number of problems 5")).toBeVisible();
  await expect(page.getByText("Qualifying problems 10 hardest")).toBeVisible();
  await expect(page.getByText("Number of finalists 7")).toBeVisible();

  await page.getByRole("button", { name: "Rules" }).click();

  await expect(page.getByText("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")).toBeVisible();
})

test.describe("contest states", () => {
  test.beforeEach(async ({ page }) => {
    await page.clock.setFixedTime(new Date('2023-11-01T00:00:00'));
  })

  test("before contest has started", async ({ page }) => {
    await page.goto('/ABCD0001');

    const timer = page.getByRole("timer", { name: "Time until start" })
    await expect(timer).toHaveText("2 months")

    await expect(page.getByRole("button", { name: "Edit" })).toBeEnabled();

    const problem = page.getByRole("region", { name: "Problem 1" });
    await expect(problem).toBeVisible();

    await expect(problem.getByRole("button", { name: "Tick" })).toBeDisabled();
  });

  test("while contest is running", async ({ page }) => {
    await page.goto('/ABCD0001');

    await page.clock.setFixedTime(new Date('2024-01-01T00:00:00'));

    const timer = page.getByRole("timer", { name: "Time remaining" })
    await expect(timer).toHaveText("almost 2 years")

    await expect(page.getByRole("button", { name: "Edit" })).toBeEnabled();

    const problem = page.getByRole("region", { name: "Problem 1" });
    await expect(problem).toBeVisible();

    await expect(problem.getByRole("button", { name: "Tick" })).toBeEnabled();
  });

  test("during grace period", async ({ page }) => {
    await page.goto('/ABCD0001');

    await page.clock.setFixedTime(new Date('2026-01-01T00:00:00'));

    const timer = page.getByRole("timer", { name: "Time remaining" })
    await expect(timer).toHaveText("00:00:00")

    await expect(page.getByRole("button", { name: "Edit" })).toBeEnabled();

    const problem = page.getByRole("region", { name: "Problem 1" });
    await expect(problem).toBeVisible();

    await expect(problem.getByRole("button", { name: "Tick" })).toBeEnabled();
  });

  test("after contest has ended", async ({ page }) => {
    await page.goto('/ABCD0001');

    await page.clock.setFixedTime(new Date('2026-01-01T00:05:00'));

    const timer = page.getByRole("timer", { name: "Time remaining" })
    await expect(timer).toHaveText("00:00:00")

    await expect(page.getByRole("button", { name: "Edit" })).toBeDisabled();

    const problem = page.getByRole("region", { name: "Problem 1" });
    await expect(problem).toBeVisible();

    await expect(problem.getByRole("button", { name: "Tick" })).toBeDisabled();
  })
})

test.describe("failsafe mode", () => {
  test('enter contest by entering registration code', async ({ page }) => {
    await page.goto('/failsafe');

    await expect(page.getByRole("heading", { name: "Welcome!" })).toBeVisible();

    const codeInput = page.getByRole("textbox", { name: "Registration code" })
    await codeInput.pressSequentially("abcd0005");

    await page.getByRole("button", { name: "Enter" }).click()

    await expect(page.getByRole("heading", { name: "Profile" })).toBeVisible();

    await page.getByRole("textbox", { name: "Name" }).pressSequentially("Andy Bernard")
    await page.getByRole("combobox", { name: "Competition class" }).selectOption({ label: "Females" });

    await page.getByRole("button", { name: "Register" }).click()

    await expect(page.getByRole("heading", { name: "Scorecard" })).toBeVisible();
  });

  test('deep link into scorecard', async ({ page }) => {
    await page.goto('/failsafe/abcd0005');

    await expect(page.getByRole("textbox", { name: "Name" })).toHaveValue("Andy Bernard");
    await expect(page.getByRole("combobox", { name: "Competition class" })).toHaveValue("2");
  });

  test("tick and untick all problems", async ({ page }) => {
    await page.goto('/failsafe/ABCD0005');

    for (let p = 1; p <= 5; p++) {
      const problem = page.getByRole("region", { name: `Problem ${p}` });
      await expect(problem).toBeVisible();

      await expect(problem.getByRole("button", { name: "Flash" })).toBeVisible();
      await problem.getByRole("button", { name: "Top" }).click();

      await expect(problem.getByRole("button", { name: "Unsend" })).toBeVisible();
    }

    for (let p = 1; p <= 5; p++) {
      const problem = page.getByRole("region", { name: `Problem ${p}` });
      await expect(problem).toBeVisible();

      await problem.getByRole("button", { name: "Unsend" }).click();

      await expect(problem.getByRole("button", { name: "Top" })).toBeVisible();
      await expect(problem.getByRole("button", { name: "Flash" })).toBeVisible();
    }
  })
})