<script lang="ts">
  import { Board, type ShallowBoard } from "$lib/api/types";
  import FormItem from "$lib/components/FormItem.svelte";
  import {
    Button,
    Checkbox,
    Dialog,
    Input,
    Label,
  } from "@nanoteck137/nano-ui";
  import type { ModalProps } from "svelte-modals";

  export type Props = {
    board: ShallowBoard;
    order: number;
  };

  export type Result = {
    name: string;
    order: number;
  };

  const { board, isHidden, isOpen, close }: Props & ModalProps<Result | null> =
    $props();

  let result = $state<Result>({
    name: board.name,
    order: board.order,
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
        <Label for="name">Name</Label>
        <Input id="name" bind:value={result.name} />
      </FormItem>

      <FormItem>
        <Label for="order">Order</Label>
        <Input id="order" type="number" bind:value={result.order} />
      </FormItem>

      <!-- <FormItem class="flex-row items-center">
        <Checkbox id="hidden" bind:checked={result.hidden} />
        <Label for="hidden">Hidden</Label>
      </FormItem> -->

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
