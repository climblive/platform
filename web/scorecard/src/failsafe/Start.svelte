<script lang="ts">
  import { ApiClient, ContenderCredentialsProvider } from "@climblive/lib";
  import type { Contender } from "@climblive/lib/models";
  import { useQueryClient } from "@tanstack/svelte-query";
  import { onMount } from "svelte";
  import EditProfile from "./EditProfile.svelte";
  import Scorecard from "./Scorecard.svelte";

  let code = $state<string>();
  let contender = $state<Contender>();

  let form: HTMLFormElement | undefined;

  const queryClient = useQueryClient();

  const authenticate = async (code: string) => {
    const contender = await ApiClient.getInstance().findContender(code);

    const provider = new ContenderCredentialsProvider(code);
    ApiClient.getInstance().setCredentialsProvider(provider);

    return contender;
  };

  const handleEnter = async (event: SubmitEvent) => {
    event.preventDefault();

    if (!form) {
      return;
    }

    const formData = new FormData(form);
    const code = formData.get("code")?.toString().trim();

    if (code && code.length === 8) {
      try {
        contender = await authenticate(code);

        queryClient.setQueryData(
          ["contender", { id: contender.id }],
          () => contender,
        );

        history.replaceState({}, "", `/failsafe/${code}`);
      } catch {}
    }
  };

  const extractCodeFromPath = () => {
    const match = window.location.pathname.match(/\/failsafe\/([A-Z0-9]{8})/i);
    return match ? match[1] : null;
  };

  onMount(() => {
    const extractedCode = extractCodeFromPath();

    if (extractedCode) {
      code = extractedCode;
    }
  });
</script>

{#if contender}
  <h2>Profile</h2>
  <EditProfile contestId={contender.contestId} contenderId={contender.id} />

  {#if contender?.entered}
    <h2>Scorecard</h2>
    <Scorecard contestId={contender.contestId} contenderId={contender.id}
    ></Scorecard>
  {/if}
{:else}
  <form bind:this={form} onsubmit={handleEnter}>
    <input
      placeholder="Registration code"
      name="code"
      type="text"
      value={code}
      minlength="8"
      maxlength="8"
    />
    <button type="submit">Enter</button>
  </form>
{/if}

<style>
  form {
    display: flex;
    gap: var(--wa-space-m);
    margin: 2rem auto;
  }

  h2:not(:first-of-type) {
    margin-top: var(--wa-space-m);
  }

  input {
    text-transform: uppercase;
    letter-spacing: 0.25rem;

    &::placeholder {
      text-transform: none;
      letter-spacing: normal;
    }
  }
</style>
