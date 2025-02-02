import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, params }) => {
  const boards = await locals.apiClient.getAllProjectBoards(params.id);
  if (!boards.success) {
    throw error(boards.error.code, { message: boards.error.message });
  }

  return {
    boards: boards.data.boards,
    hiddenBoards: boards.data.hiddenBoards,
  };
};
