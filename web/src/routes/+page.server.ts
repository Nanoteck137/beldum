import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals }) => {
  const projects = await locals.apiClient.getUserProjects();
  if (!projects.success) {
    throw error(projects.error.code, { message: projects.error.message });
  }

  return {
    projects: projects.data.projects,
  };
};
