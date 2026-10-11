import { test as base } from "@playwright/test";
import {
  MariaDbContainer,
  StartedMariaDbContainer,
} from "@testcontainers/mariadb";
import { readFile } from "fs/promises";
import { Connection, createConnection } from "mariadb";
import path from "path";
import {
  GenericContainer,
  Network,
  StartedTestContainer,
  Wait,
} from "testcontainers";

export { expect } from "@playwright/test";

export const test = base.extend<{}, { app: void }>({
  app: [
    async ({}, use) => {
      let dbConnection: Connection | undefined;
      let startedDbContainer: StartedMariaDbContainer | undefined;
      let startedAppContainer: StartedTestContainer | undefined;
      const network = await new Network().start();

      try {
        startedDbContainer = await new MariaDbContainer("mariadb:11.4")
          .withUsername("climblive")
          .withUserPassword("secretpassword")
          .withDatabase("climblive")
          .withExposedPorts(3306)
          .withNetwork(network)
          .withNetworkAliases("e2e")
          .start();

        dbConnection = await createConnection({
          host: startedDbContainer.getHost(),
          port: startedDbContainer.getMappedPort(3306),
          user: "climblive",
          password: "secretpassword",
          database: "climblive",
          multipleStatements: true,
        });

        const schema = await readFile(
          "../../backend/database/climblive.sql",
          "utf8",
        );
        const samples = await readFile("./samples.sql", "utf8");

        await dbConnection.query(schema);
        await dbConnection.query(samples);

        const appContainer = new GenericContainer("climblive-api:latest")
          .withEnvironment({
            DB_USERNAME: "climblive",
            DB_PASSWORD: "secretpassword",
            DB_HOST: "e2e",
            DB_PORT: "3306",
            DB_DATABASE: "climblive",
            RUN_AS_USER: "climblive",
            TLS_APP_CERT_FILE: "/certs/cert.pem",
            TLS_APP_KEY_FILE: "/certs/key.pem",
            TLS_WWW_CERT_FILE: "/certs/cert.pem",
            TLS_WWW_KEY_FILE: "/certs/key.pem",
          })
          .withNetwork(network)
          .withCopyDirectoriesToContainer([
            {
              source: path.resolve(__dirname, ".local/certs"),
              target: "/certs",
            },
          ])
          .withExposedPorts({ container: 443, host: 8443 })
          .withWaitStrategy(Wait.forLogMessage(/score engine started/));

        startedAppContainer = await appContainer.start();
        await use();
      } finally {
        await startedAppContainer?.stop();
        await dbConnection?.end();
        await startedDbContainer?.stop();
        await network.stop();
      }
    },
    { scope: "worker", auto: true },
  ],
});
