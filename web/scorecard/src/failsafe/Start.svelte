<script lang="ts">
  import "@awesome.me/webawesome/dist/components/callout/callout.js";
  import "@awesome.me/webawesome/dist/components/icon/icon.js";
  import { ApiClient, ContenderCredentialsProvider } from "@climblive/lib";
  import type { Contender } from "@climblive/lib/models";
  import { getContestQuery } from "@climblive/lib/queries";
  import { add, formatDistance } from "date-fns";
  import { useQueryClient } from "@tanstack/svelte-query";
  import { onMount } from "svelte";
  import Scorecard from "./Scorecard.svelte";

  let code = $state<string>();
  let contender = $state<Contender>();
  const contestQuery = $derived(
    contender ? getContestQuery(contender.contestId) : undefined,
  );
  const contest = $derived(contestQuery?.data);
  const retentionDuration = $derived.by(() => {
    const base = new Date(0);
    return formatDistance(
      add(base, { minutes: (contest?.nameRetentionTime ?? 0) / 60000000000 }),
      base,
    );
  });

  let form: HTMLFormElement | undefined = $state();

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
      tryEnter(code);
    }
  };

  const tryEnter = async (code: string) => {
    contender = await authenticate(code);

    queryClient.setQueryData(
      ["contender", { id: contender.id }],
      () => contender,
    );

    history.replaceState({}, "", `/failsafe/${code}`);
  };

  const extractCodeFromPath = () => {
    const match = window.location.pathname.match(/\/failsafe\/([A-Z0-9]{8})/i);
    return match ? match[1] : null;
  };

  onMount(() => {
    const extractedCode = extractCodeFromPath();

    if (extractedCode) {
      code = extractedCode;

      tryEnter(extractedCode);
    }
  });
</script>

{#if contender}
  {#if !contender.entered && contest}<wa-callout variant="neutral" size="s"
      ><wa-icon slot="icon" name="circle-info"></wa-icon>Your name will be
      stored for {retentionDuration} after the competition ends, after which it will
      be removed and your results anonymized.</wa-callout
    >{/if}
  <Scorecard contestId={contender.contestId} contenderId={contender.id} />
{:else}
  <h2>Welcome!</h2>
  <form bind:this={form} onsubmit={handleEnter}>
    <div>
      <small>Input your 8 digit registration code.</small>
      <input
        placeholder="ABCD0001"
        name="code"
        type="text"
        value={code}
        minlength="8"
        maxlength="8"
        aria-label="Registration code"
      />
    </div>
    <button type="submit">Enter</button>
  </form>
{/if}

<style>
  form {
    display: flex;
    flex-direction: column;
    gap: var(--wa-space-m);
  }

  h2 {
    margin-top: var(--wa-space-m);
  }

  small {
    display: block;
    margin-block-end: var(--wa-space-xs);
  }

  input {
    width: 100%;
    display: block;
    text-transform: uppercase;
    letter-spacing: 0.25rem;
  }
</style>
