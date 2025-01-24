<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";
  import { modals } from "svelte-modals";
  import NewProjectModal from "./NewProjectModal.svelte";
  import { getApiClient, handleApiError } from "$lib";
  import { goto } from "$app/navigation";
  import Spacer from "$lib/components/Spacer.svelte";

  const { data } = $props();
  const apiClient = getApiClient();
</script>

<Button
  onclick={async () => {
    const modalData = await modals.open(NewProjectModal, {});
    if (modalData) {
      const res = await apiClient.createProject({
        name: modalData.name,
      });
      if (!res.success) {
        handleApiError(res.error);
        return;
      }

      goto(`/projects/${res.data.id}`);
    }
  }}
>
  New Project
</Button>

<Spacer size="md" />

<div class="flex flex-col gap-2">
  {#each data.projects as project}
    <div class="rounded bg-primary px-2 py-2 text-primary-foreground">
      <a href="/projects/{project.id}">{project.name}</a>
    </div>
  {/each}
</div>
