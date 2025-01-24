<script lang="ts">
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Dialog, Input, Label } from "@nanoteck137/nano-ui";
  import type { ModalProps } from "svelte-modals";

  export type Props = {};

  export type Result = {
    title: string;
    tags: string;
  };

  const { isOpen, close }: Props & ModalProps<Result | null> = $props();

  let result = $state<Result>({
    title: "",
    tags: "",
  });
</script>

<Dialog.Root
  controlledOpen
  open={isOpen}
  onOpenChange={(v) => {
    if (!v) {
      close(null);
    }
  }}
>
  <Dialog.Content>
    <form
      class="flex flex-col gap-4"
      onsubmit={(e) => {
        e.preventDefault();
        close(result);
      }}
    >
      <Dialog.Header>
        <Dialog.Title>Create new Task</Dialog.Title>
      </Dialog.Header>

      <FormItem>
        <Label for="title">Title</Label>
        <Input id="title" bind:value={result.title} />
      </FormItem>

      <FormItem>
        <Label for="tags">Tags</Label>
        <Input id="tags" bind:value={result.tags} />
      </FormItem>

      <Dialog.Footer>
        <Button
          variant="outline"
          onclick={() => {
            close(null);
          }}
        >
          Close
        </Button>
        <Button type="submit">Create</Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
