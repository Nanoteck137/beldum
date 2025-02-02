<script lang="ts">
  import { cn } from "$lib/utils";
  import {
    Button,
    buttonVariants,
    DropdownMenu,
    ScrollArea,
  } from "@nanoteck137/nano-ui";
  import {
    ChevronLeft,
    ChevronRight,
    Edit,
    EllipsisVertical,
  } from "lucide-svelte";
  import { modals } from "svelte-modals";
  import NewBoardModal from "./NewBoardModal.svelte";
  import { getApiClient, handleApiError } from "$lib";
  import { invalidateAll } from "$app/navigation";
  import toast from "svelte-5-french-toast";
  import EditBoardModal from "./EditBoardModal.svelte";
  import { ShallowBoard } from "$lib/api/types";

  const { data } = $props();
  const apiClient = getApiClient();
</script>

<Button
  onclick={async () => {
    const newModal = await modals.open(NewBoardModal, {});
    if (newModal) {
      const res = await apiClient.createBoard({
        name: newModal.name,
        hidden: newModal.hidden,
        projectId: data.project.id,
      });
      if (!res.success) {
        handleApiError(res.error);
        invalidateAll();
        return;
      }

      toast.success("Successfully created board");
      invalidateAll();
    }
  }}
>
  New Board
</Button>

<Button
  onclick={async () => {
    for (let i = 0; i < data.boards.length; i++) {
      const res = await apiClient.editBoard(data.boards[i].id, {
        order: i + 1,
      });
      if (!res.success) {
        handleApiError(res.error);
        invalidateAll();
        return;
      }
    }

    toast.success("Setting order done");
    invalidateAll();
  }}
>
  Set Order
</Button>

{#snippet dropdownMenu(board: ShallowBoard, isHidden: boolean)}
  <DropdownMenu.Root>
    <DropdownMenu.Trigger
      class={cn(buttonVariants({ variant: "ghost", size: "icon" }))}
    >
      <EllipsisVertical />
    </DropdownMenu.Trigger>
    <DropdownMenu.Content align="end">
      <DropdownMenu.Group>
        <DropdownMenu.Item
          onSelect={async () => {
            const editModal = await modals.open(EditBoardModal, {
              board,
              isHidden,
            });
            if (editModal) {
              const res = await apiClient.editBoard(board.id, {
                name: editModal.name,
                order: editModal.order,
              });
              if (!res.success) {
                handleApiError(res.error);
                invalidateAll();
                return;
              }

              toast.success("Successfully edited board");
              invalidateAll();
            }
          }}
        >
          <Edit />
          Edit Board
        </DropdownMenu.Item>
      </DropdownMenu.Group>
    </DropdownMenu.Content>
  </DropdownMenu.Root>
{/snippet}

<ScrollArea orientation="horizontal">
  <div class="flex w-full gap-4 px-2 py-4">
    {#each data.boards as board, i}
      {@const hasPrev = data.boards.length > 1 && i >= 1}
      {@const hasNext = data.boards.length > 1 && i < data.boards.length - 1}

      <div class="flex w-96 min-w-96 flex-col rounded border">
        <div class="flex gap-2 border-b p-2">
          <Button class="w-full" variant="ghost" disabled={!hasPrev}>
            <ChevronLeft />
          </Button>

          <Button class="w-full" variant="ghost" disabled={!hasNext}>
            <ChevronRight />
          </Button>
        </div>
        <div class="flex items-center justify-between px-4 py-2">
          <p class="text-lg">{board.order}. {board.name}</p>
          {@render dropdownMenu(board, false)}
        </div>
      </div>
    {/each}
  </div>
</ScrollArea>

<ScrollArea orientation="horizontal">
  <div class="flex w-full gap-4 px-2 py-4">
    {#each data.hiddenBoards as board, i}
      <div class="flex w-96 min-w-96 flex-col rounded border">
        <div class="flex items-center justify-between px-4 py-2">
          <p class="text-lg">{board.name}</p>
          {@render dropdownMenu(board, true)}
        </div>
      </div>
    {/each}
  </div>
</ScrollArea>
