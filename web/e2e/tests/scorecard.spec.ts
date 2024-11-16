import { expect, test } from '@playwright/test';
import {
  MariaDbContainer
} from "@testcontainers/mariadb";
import { readFile } from "fs/promises";
import { createConnection } from "mariadb";
import { GenericContainer } from "testcontainers";

test.beforeAll(async () => {
  const dbContainer = await new MariaDbContainer()
    .withUsername("climblive")
    .withUserPassword("secretpassword")
    .withDatabase("climblive")
    .withExposedPorts(3306)
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
  const samples = await readFile("../../backend/database/samples.sql", "utf8")

  dbConnection.query(schema)
  dbConnection.query(samples)

  const apiContainer = await GenericContainer
    .fromDockerfile("../../backend")
    .build();

  apiContainer.withExposedPorts({ container: 8090, host: 8090 }).start()
})

test('enter contest by entering code', async ({ page }) => {
  await page.goto('/');

  await expect(page).toHaveTitle(/ClimbLive/);

  const pinInput = await page.getByRole("textbox").first()
  pinInput.pressSequentially("abcd0001");

  await page.getByRole("button", { name: "Enter" }).click();
});
