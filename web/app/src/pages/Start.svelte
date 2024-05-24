<script lang="ts">
  import "@shoelace-style/shoelace/dist/components/alert/alert.js";
  import "@shoelace-style/shoelace/dist/components/button/button.js";
  import "@shoelace-style/shoelace/dist/components/icon/icon.js";
  import "@shoelace-style/shoelace/dist/components/input/input.js";
  import { useQueryClient } from "@tanstack/svelte-query";
  import { differenceInHours } from "date-fns";
  import { getContext, onMount } from "svelte";
  import { navigate } from "svelte-routing";
  import type { Writable } from "svelte/store";
  import PinInput from "@climblive/shared/components/PinInput.svelte";
  import { scorecardSessionSchema, type ScorecardSession } from "../types";
  import { authenticateContender } from "../utils/auth";

  let loadingContender = false;
  let loadingFailed = false;
  let queryClient = useQueryClient();
  let registrationCode: string | undefined;

  onMount(() => {
    const data = localStorage.getItem("session");
    if (data) {
      try {
        const obj = JSON.parse(data);
        const sess = scorecardSessionSchema.parse(obj);

        if (differenceInHours(new Date(), sess.timestamp) < 12) {
          registrationCode = sess.registrationCode;
        }
      } catch (_) {}
    }
  });

  const session = getContext<Writable<ScorecardSession>>("scorecardSession");

  const handleCodeChange = (code: string) => {
    const autoSubmit = registrationCode === undefined && code.length === 8;
    registrationCode = code;
    autoSubmit && submitForm();
  };

  const submitForm = async () => {
    if (!registrationCode) {
      return;
    }

    try {
      loadingFailed = false;
      loadingContender = true;

      const contender = await authenticateContender(
        registrationCode,
        queryClient,
        session
      );

      if (contender.entered) {
        navigate(`/${registrationCode}`);
      } else {
        navigate(`/${registrationCode}/register`);
      }
    } catch (e) {
      loadingFailed = true;
    } finally {
      loadingContender = false;
    }
  };
</script>

<main>
  <header>
    <h1>Welcome</h1>
    <p>Enter your unique registration code:</p>
  </header>
  <form on:submit|preventDefault={submitForm}>
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
      loading={loadingContender}>Enter</sl-button
    >
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
</style>
