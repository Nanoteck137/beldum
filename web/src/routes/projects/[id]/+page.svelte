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
    EllipsisVertical,
    Clipboard,
    Trash,
  } from "lucide-svelte";
  import { modals } from "svelte-modals";
  import NewTaskModal from "./NewTaskModal.svelte";
  import { getApiClient, handleApiError, openConfirm } from "$lib";
  import { invalidateAll } from "$app/navigation";
  import toast from "svelte-5-french-toast";

  const { data } = $props();
  const apiClient = getApiClient();
</script>

<p>{data.project.id}</p>
<p>{data.project.name}</p>

<ScrollArea orientation="horizontal">
  <div class="flex w-full gap-4 px-2 py-4">
    {#each data.boards as board, i}
      {@const hasPrev = data.boards.length > 1 && i >= 1}
      {@const hasNext = data.boards.length > 1 && i < data.boards.length - 1}
      <div class="flex min-h-[480px] w-96 min-w-96 flex-col rounded border">
        <div class="flex items-center justify-between border-b px-4 py-2">
          <p class="text-lg">{board.name}</p>

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
                    const modalData = await modals.open(NewTaskModal, {});
                    if (modalData) {
                      const res = await apiClient.createTask({
                        title: modalData.title,
                        tags: modalData.tags.split(","),
                        boardId: board.id,
                      });
                      if (!res.success) {
                        handleApiError(res.error);
                        invalidateAll();
                        return;
                      }

                      toast.success("Successfully created task");
                      invalidateAll();
                    }
                  }}
                >
                  <Clipboard />
                  New Task
                </DropdownMenu.Item>
              </DropdownMenu.Group>
            </DropdownMenu.Content>
          </DropdownMenu.Root>
        </div>

        <div class="flex flex-grow flex-col gap-2 px-4 py-2">
          {#each board.items as item}
            <div
              class="flex flex-col gap-2 rounded bg-primary p-4 text-primary-foreground"
            >
              <div class="flex items-center justify-between gap-4">
                <p class="text-sm">{item.name}</p>
              </div>

              <div class="flex flex-wrap gap-2">
                {#each item.tags as tag}
                  <p
                    class="flex items-center justify-center rounded-full bg-red-200 px-4 py-2 text-xs"
                  >
                    {tag}
                  </p>
                {/each}
              </div>

              <div class="flex justify-between">
                <div>
                  <DropdownMenu.Root>
                    <DropdownMenu.Trigger
                      class={cn(
                        buttonVariants({ variant: "ghost", size: "icon" }),
                      )}
                    >
                      <EllipsisVertical />
                    </DropdownMenu.Trigger>
                    <DropdownMenu.Content align="end">
                      <DropdownMenu.Group>
                        <DropdownMenu.Item
                          onSelect={async () => {
                            const confirmed = await openConfirm({
                              title: "Are you sure?",
                            });
                            if (confirmed) {
                              const res = await apiClient.deleteTask(item.id);
                              if (!res.success) {
                                handleApiError(res.error);
                                invalidateAll();
                                return;
                              }

                              toast.success("Successfully deleted task");
                              invalidateAll();
                            }
                          }}
                        >
                          <Trash />
                          Delete Task
                        </DropdownMenu.Item>
                      </DropdownMenu.Group>
                    </DropdownMenu.Content>
                  </DropdownMenu.Root>
                </div>

                <div class="flex gap-2">
                  <Button
                    variant="secondary"
                    size="icon"
                    disabled={!hasPrev}
                    onclick={async () => {
                      let newBoard = data.boards[i - 1];
                      const res = await apiClient.moveTask(
                        item.id,
                        newBoard.id,
                      );
                      if (!res.success) {
                        handleApiError(res.error);
                      }

                      invalidateAll();
                    }}
                  >
                    <ChevronLeft />
                  </Button>

                  <Button
                    variant="secondary"
                    size="icon"
                    disabled={!hasNext}
                    onclick={async () => {
                      let newBoard = data.boards[i + 1];
                      const res = await apiClient.moveTask(
                        item.id,
                        newBoard.id,
                      );
                      if (!res.success) {
                        handleApiError(res.error);
                      }

                      invalidateAll();
                    }}
                  >
                    <ChevronRight />
                  </Button>
                </div>
              </div>
            </div>
          {/each}
        </div>

        <div class="flex items-center justify-between border-t px-4 py-2">
          <p class="text-lg">{board.name}</p>

          <Button variant="ghost" size="icon">
            <EllipsisVertical />
          </Button>
        </div>
      </div>
    {/each}
  </div>
</ScrollArea>
