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

  const pinInput = await page.getByRole("textbox").first()
  pinInput.pressSequentially("abcd0002");

  await page.getByRole("textbox", { name: "Full name *" }).pressSequentially("Dwight Schrute")
  await page.getByRole("textbox", {
    name: "Club name"
  }).pressSequentially("Scranton Climbing Club")
  const compClass = page.getByRole("combobox", { name: "Competition class *" })
  await compClass.click()
  page.getByRole("option", { name: "Males", exact: true }).click()

  await page.getByRole("button", { name: "Register" }).click()

  await page.getByText("0p")
});
