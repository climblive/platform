import { expect, test } from '@playwright/test';
import {
  MariaDbContainer
} from "@testcontainers/mariadb";
import { readFile } from "fs/promises";
import { createConnection } from "mariadb";
import { GenericContainer, Network, Wait } from "testcontainers";

test.beforeAll(async () => {
  const network = await new Network().start();

  const dbContainer = await new MariaDbContainer()
    .withUsername("climblive")
    .withUserPassword("secretpassword")
    .withDatabase("climblive")
    .withExposedPorts({ container: 3306, host: 3306 })
    .withNetwork(network)
    .withNetworkAliases("e2e")
    .start();

  var dbConnection = await createConnection({
    host: "localhost",
    port: dbContainer.getMappedPort(3306),
    user: "climblive",
    password: "secretpassword",
    database: "climblive",
    multipleStatements: true,
  });

  const schema = await readFile("../../backend/database/scoreboard.sql", "utf8")
  const samples = await readFile("./samples.sql", "utf8")

  dbConnection.query(schema)
  dbConnection.query(samples)

  const apiContainer = await GenericContainer
    .fromDockerfile("../../backend")
    .build();

  await apiContainer
    .withEnvironment({
      "DB_USERNAME": "climblive",
      "DB_PASSWORD": "secretpassword",
      "DB_HOST": "e2e",
      "DB_DATABASE": "climblive",
    })
    .withNetwork(network)
    .withExposedPorts({ container: 8090, host: 8090 })
    .withWaitStrategy(Wait.forHttp("/contests/1", 8090))
    .start()
})

test('enter contest by entering code', async ({ page }) => {
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
