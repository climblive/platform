<script lang="ts">
  import { type ScorecardSession } from "@/types";
  import { authenticateContender, readStoredSessions } from "@/utils/auth";
  import { serialize } from "@shoelace-style/shoelace";
  import "@shoelace-style/shoelace/dist/components/alert/alert.js";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import "@shoelace-style/shoelace/dist/components/icon/icon.js";
  import "@shoelace-style/shoelace/dist/components/input/input.js";
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
    <sl-input
      required
      placeholder="ABCD1234"
      label="Registration code"
      help-text="Input your 8 digit registration code"
      name="code"
      type="text"
      minlength="8"
      maxlength="8"
    >
      <sl-icon name="key" slot="prefix"></sl-icon>
    </sl-input>
    {#if loadingFailed}
      <sl-alert open variant="danger">
        <sl-icon slot="icon" name="exclamation-octagon"></sl-icon>
        The registration code is not valid.
      </sl-alert>
    {/if}
    <sl-button variant="primary" type="submit" loading={loadingContender}>
      <sl-icon slot="prefix" name="box-arrow-in-right"></sl-icon>
      Enter
    </sl-button>
  </form>
  {#each restoredSessions as restoredSession (restoredSession.registrationCode)}
    <sl-divider style="--color: var(--sl-color-primary-600);"></sl-divider>
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
      <sl-button
        onclick={() => {
          if (restoredSession) {
            handleEnter(restoredSession.registrationCode);
          }
        }}
        loading={loadingContender}
        size="small"
        >Restore
      </sl-button>
    </section>
  {/each}
  <footer>by ClimbLive™</footer>
</main>

<style>
  main {
    display: flex;
    flex-direction: column;
    padding-inline: var(--sl-spacing-large);
  }

  header {
    margin-top: 25%;
  }

  form {
    display: flex;
    flex-direction: column;
    text-align: left;
    gap: var(--sl-spacing-small);

    & sl-input::part(input) {
      text-transform: uppercase;
      font-family: monospace;
      white-space: pre;

      width: 100%;
    }
  }

  footer {
    margin-top: auto;
    text-align: center;
    font-weight: var(--sl-font-weight-semibold);
    line-height: 4rem;
    font-size: var(--sl-font-size-x-small);
    color: var(--sl-color-primary-900);
  }

  .restoredSession {
    background-color: var(--sl-color-primary-600);
    border-radius: var(--sl-border-radius-medium);
    padding: var(--sl-spacing-small);
    text-align: left;
    color: white;

    & h3 {
      margin: 0;
      font-weight: normal;
    }

    & .timestamp {
      font-size: var(--sl-font-size-x-small);
    }

    & sl-button {
      width: 100%;
    }

    & .code {
      text-transform: uppercase;
      font-weight: bold;
    }
  }
</style>
