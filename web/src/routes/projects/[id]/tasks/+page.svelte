<script lang="ts">
  const { data } = $props();
</script>

{#each data.tasks as task}
  <p>{task.name}</p>
{/each}
