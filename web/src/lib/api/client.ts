import { z } from "zod";
import * as api from "./types";
import { BaseApiClient, type ExtraOptions } from "./base-client";


export class ApiClient extends BaseApiClient {
  constructor(baseUrl: string) {
    super(baseUrl);
  }
  
  createProject(body: api.CreateProjectBody, options?: ExtraOptions) {
    return this.request("/api/v1/projects", "POST", api.CreateProject, z.any(), body, options)
  }
  
  getUserProjects(options?: ExtraOptions) {
    return this.request("/api/v1/projects", "GET", api.GetProjects, z.any(), undefined, options)
  }
  
  getProjectById(projectId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/projects/${projectId}`, "GET", api.GetProjectById, z.any(), undefined, options)
  }
  
  getProjectBoards(projectId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/projects/${projectId}/boards`, "GET", api.GetProjectBoards, z.any(), undefined, options)
  }
  
  createTask(boardId: string, body: api.CreateTaskBody, options?: ExtraOptions) {
    return this.request(`/api/v1/boards/${boardId}/tasks`, "POST", api.CreateTask, z.any(), body, options)
  }
  
  moveTask(taskId: string, boardId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/tasks/${taskId}/move/${boardId}`, "POST", z.undefined(), z.any(), undefined, options)
  }
  
  signup(body: api.SignupBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/signup", "POST", api.Signup, z.any(), body, options)
  }
  
  signin(body: api.SigninBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/signin", "POST", api.Signin, z.any(), body, options)
  }
  
  changePassword(body: api.ChangePasswordBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/password", "PATCH", z.undefined(), z.any(), body, options)
  }
  
  getMe(options?: ExtraOptions) {
    return this.request("/api/v1/auth/me", "GET", api.GetMe, z.any(), undefined, options)
  }
  
  getSystemInfo(options?: ExtraOptions) {
    return this.request("/api/v1/system/info", "GET", api.GetSystemInfo, z.any(), undefined, options)
  }
  
  updateUserSettings(body: api.UpdateUserSettingsBody, options?: ExtraOptions) {
    return this.request("/api/v1/user/settings", "PATCH", z.undefined(), z.any(), body, options)
  }
  
  createApiToken(body: api.CreateApiTokenBody, options?: ExtraOptions) {
    return this.request("/api/v1/user/apitoken", "POST", api.CreateApiToken, z.any(), body, options)
  }
  
  getAllApiTokens(options?: ExtraOptions) {
    return this.request("/api/v1/user/apitoken", "GET", api.GetAllApiTokens, z.any(), undefined, options)
  }
  
  deleteApiToken(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/user/apitoken/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
}
