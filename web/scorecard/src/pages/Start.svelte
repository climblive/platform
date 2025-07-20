<script lang="ts">
  import logoUrl from "@/static/logo.svg";
  import { type ScorecardSession } from "@/types";
  import { authenticateContender, readStoredSessions } from "@/utils/auth";
  import { serialize } from "@awesome.me/webawesome";
  import "@awesome.me/webawesome/dist/components/button/button.js";
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/divider/divider.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import "@awesome.me/webawesome/dist/components/input/input.js";
  import { useQueryClient } from "@tanstack/svelte-query";
  import { format } from "date-fns";
  import { getContext, onMount } from "svelte";
  import { navigate } from "svelte-routing";
  import type { Writable } from "svelte/store";
  import * as z from "zod";
  import { ZodError } from "zod";

  const enterFormSchema = z.object({
    code: z.string().length(8),
  });

  let loadingContender = $state(false);
  let loadingFailed = $state(false);
  let queryClient = useQueryClient();
  let form: HTMLFormElement | undefined = $state();
  let restoredSessions: ScorecardSession[] = $state([]);

  onMount(() => {
    restoredSessions = readStoredSessions();
  });

  const session = getContext<Writable<ScorecardSession>>("scorecardSession");

  const handleSubmit = async (event: SubmitEvent) => {
    event.preventDefault();

    if (!form) {
      return;
    }

    const { data, success } = enterFormSchema.safeParse(serialize(form));

    if (success) {
      handleEnter(data.code.toUpperCase());
    }
  };

  const handleEnter = async (registrationCode: string) => {
    try {
      loadingFailed = false;
      loadingContender = true;

      const contender = await authenticateContender(
        registrationCode,
        queryClient,
        session,
      );

      if (contender.entered) {
        navigate(`/${registrationCode}`);
      } else {
        navigate(`/${registrationCode}/register`);
      }
    } catch (e) {
      if (e instanceof ZodError) {
        // eslint-disable-next-line no-console
        console.error(e);
      }

      loadingFailed = true;
    } finally {
      loadingContender = false;
    }
  };
</script>

<main>
  <header>
    <h1>Welcome!</h1>
  </header>
  <form bind:this={form} onsubmit={handleSubmit}>
    <wa-input
      required
      placeholder="ABCD1234"
      label="Registration code"
      hint="Input your 8 digit registration code"
      name="code"
      type="text"
      minlength="8"
      maxlength="8"
    >
      <wa-icon name="key" slot="start"></wa-icon>
    </wa-input>
    {#if loadingFailed}
      <wa-callout open variant="danger">
        <wa-icon slot="icon" name="exclamation-octagon"></wa-icon>
        The registration code is not valid.
      </wa-callout>
    {/if}
    <wa-button variant="brand" type="submit" loading={loadingContender}>
      <wa-icon slot="start" name="arrow-right-to-bracket"></wa-icon>
      Enter
    </wa-button>
  </form>

  {#if restoredSessions.length > 0}
    <wa-divider></wa-divider>
  {/if}

  {#each restoredSessions as restoredSession (restoredSession.registrationCode)}
    <section
      class="restoredSession"
      aria-label="Saved session {restoredSession.registrationCode}"
    >
      <h3>
        Saved session <span class="code"
          >{restoredSession.registrationCode}</span
        >
      </h3>
      <p class="timestamp">{format(restoredSession.timestamp, "pp")}</p>
      <wa-button
        onclick={() => {
          if (restoredSession) {
            handleEnter(restoredSession.registrationCode);
          }
        }}
        loading={loadingContender}
        size="small"
        appearance="outlined filled"
        >Restore
        <wa-icon slot="start" name="arrow-right-to-bracket"></wa-icon>
      </wa-button>
    </section>
  {/each}
  <footer>
    <img src={logoUrl} alt="ClimbLive" />
  </footer>
</main>

<style>
  main {
    display: flex;
    flex-direction: column;
    padding-inline: var(--wa-space-l);
    min-height: 100vh;
  }

  header {
    margin-top: 25%;
  }

  form {
    display: flex;
    flex-direction: column;
    text-align: left;
    gap: var(--wa-space-s);

    & wa-input::part(input) {
      text-transform: uppercase;
      font-family: monospace;
      white-space: pre;

      width: 100%;
    }
  }

  footer {
    margin-top: auto;
    text-align: center;
    padding-block: var(--wa-space-m);

    & img {
      height: var(--wa-font-size-l);
    }
  }

  .restoredSession {
    background-color: var(--wa-color-surface-default);
    border: var(--wa-border-width-s) var(--wa-border-style)
      var(--wa-color-surface-border);
    border-radius: var(--wa-border-radius-m);
    padding: var(--wa-space-s);
    text-align: left;

    & h3 {
      margin: 0;
      font-weight: normal;
    }

    & .timestamp {
      font-size: var(--wa-font-size-xs);
    }

    & wa-button {
      width: 100%;
    }

    & .code {
      text-transform: uppercase;
      font-weight: bold;
    }
  }

  .restoredSession:not(:last-of-type) {
    margin-bottom: var(--wa-space-s);
  }
</style>
