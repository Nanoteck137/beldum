import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, params }) => {
  const project = await locals.apiClient.getProjectById(params.id);
  if (!project.success) {
    throw error(project.error.code, { message: project.error.message });
  }

  const boards = await locals.apiClient.getProjectBoards(params.id);
  if (!boards.success) {
    throw error(boards.error.code, { message: boards.error.message });
  }

  return {
    project: project.data,
    boards: boards.data.boards,
  };
};
