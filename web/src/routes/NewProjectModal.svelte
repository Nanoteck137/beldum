<script lang="ts">
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Dialog, Input, Label } from "@nanoteck137/nano-ui";
  import type { ModalProps } from "svelte-modals";

  export type Props = {};

  export type Result = {
    name: string;
  };

  const { isOpen, close }: Props & ModalProps<Result | null> = $props();

  let result = $state<Result>({
    name: "",
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
        <Dialog.Title>Create new Project</Dialog.Title>
      </Dialog.Header>

      <FormItem>
        <Label for="name">Name</Label>
        <Input id="name" bind:value={result.name} />
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
