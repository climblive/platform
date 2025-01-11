<script lang="ts">
  import { scorecardSessionSchema, type ScorecardSession } from "@/types";
  import { authenticateContender } from "@/utils/auth";
  import { PinInput } from "@climblive/lib/components";
  import "@shoelace-style/shoelace/dist/components/alert/alert.js";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import "@shoelace-style/shoelace/dist/components/icon/icon.js";
  import "@shoelace-style/shoelace/dist/components/input/input.js";
  import { useQueryClient } from "@tanstack/svelte-query";
  import { differenceInHours, format } from "date-fns";
  import { getContext, onMount } from "svelte";
  import { navigate } from "svelte-routing";
  import type { Writable } from "svelte/store";
  import { ZodError } from "zod";

  let loadingContender = $state(false);
  let loadingFailed = $state(false);
  let queryClient = useQueryClient();
  let registrationCode: string | undefined = $state();
  let form: HTMLFormElement | undefined = $state();
  let restoredSession: ScorecardSession | undefined = $state();

  onMount(() => {
    const data = localStorage.getItem("session");
    if (data) {
      try {
        const obj = JSON.parse(data);
        const sess = scorecardSessionSchema.parse(obj);

        if (differenceInHours(new Date(), sess.timestamp) < 12) {
          restoredSession = sess;
        }
      } catch {
        /* discard corrupt session data */
      }
    }
  });

  const session = getContext<Writable<ScorecardSession>>("scorecardSession");

  const handleCodeChange = (code: string) => {
    const autoSubmit = registrationCode === undefined && code.length === 8;
    registrationCode = code;
    if (autoSubmit) {
      form?.requestSubmit();
    }
  };

  const submitForm = async (event: SubmitEvent) => {
    event.preventDefault();

    if (!registrationCode) {
      return;
    }

    handleEnter(registrationCode);
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
    <h1>Welcome</h1>
    <p>Enter your unique registration code!</p>
  </header>
  <form bind:this={form} onsubmit={submitForm}>
    <PinInput
      length={8}
      defaultValue={registrationCode}
      onChange={handleCodeChange}
      disabled={loadingContender}
    />
    {#if loadingFailed}
      <sl-alert open variant="danger">
        <sl-icon slot="icon" name="exclamation-octagon"></sl-icon>
        The registration code is not valid.
      </sl-alert>
    {/if}
    <sl-button
      variant="primary"
      type="submit"
      size="small"
      disabled={registrationCode?.length !== 8}
      loading={loadingContender}
      >Enter
    </sl-button>
    {#if restoredSession}
      <div class="restoredSession">
        <h3>Saved session {restoredSession.registrationCode}</h3>
        <p class="timestamp">{format(restoredSession.timestamp, "PPpp")}</p>
        <sl-button
          onclick={() => {
            restoredSession && handleEnter(restoredSession.registrationCode);
          }}
          loading={loadingContender}
          size="small"
          >Restore
        </sl-button>
      </div>
    {/if}
  </form>
  <footer>by ClimbLive™</footer>
</main>

<style>
  main {
    width: 100%;
    height: 100vh;
    text-align: center;
    display: flex;
    flex-direction: column;
  }

  header {
    margin-top: 25%;
  }

  form {
    display: grid;
    grid-template-columns: min-content;
    gap: var(--sl-spacing-small);
    justify-content: center;
    margin-bottom: auto;
  }

  footer {
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
    }

    & .timestamp {
      font-size: var(--sl-font-size-x-small);
    }

    & sl-button {
      width: 100%;
    }
  }
</style>
