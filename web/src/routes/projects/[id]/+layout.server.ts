import { error } from "@sveltejs/kit";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ locals, params }) => {
  const project = await locals.apiClient.getProjectById(params.id);
  if (!project.success) {
    throw error(project.error.code, { message: project.error.message });
  }

  return {
    project: project.data,
  };
};
