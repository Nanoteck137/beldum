import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, params }) => {
  const tasks = await locals.apiClient.getProjectTasks(params.id);
  if (!tasks.success) {
    throw error(tasks.error.code, { message: tasks.error.message });
  }

  return {
    tasks: tasks.data.tasks,
  };
};
