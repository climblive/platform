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

  const apiContainerBuilder = GenericContainer
    .fromDockerfile("../../backend")

  const webContainerBuilder = GenericContainer
    .fromDockerfile("..")

  let [apiContainer, webContainer] = await Promise.all([apiContainerBuilder.build(), webContainerBuilder.build()]);

  apiContainer = apiContainer
    .withEnvironment({
      "DB_USERNAME": "climblive",
      "DB_PASSWORD": "secretpassword",
      "DB_HOST": "e2e",
      "DB_PORT": "3306",
      "DB_DATABASE": "climblive",
    })
    .withNetwork(network)
    .withExposedPorts({ container: 8090, host: 8090 })
    .withWaitStrategy(Wait.forHttp("/contests/1", 8090))

  webContainer = webContainer
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

  await expect(page.getByText("2nd place")).toBeVisible()
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
