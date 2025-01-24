// place files you want to import through the `$lib` alias in this folder.

import { ApiClient } from "$lib/api/client";
import ConfirmModal, {
  type Props as ConfirmModalProps,
} from "$lib/components/modals/ConfirmModal.svelte";
import InputModal, {
  type Props as InputModalProps,
} from "$lib/components/modals/InputModal.svelte";
import { getContext, setContext } from "svelte";
import toast from "svelte-5-french-toast";
import { modals } from "svelte-modals";

export function handleApiError(err: {
  code: number;
  type: string;
  message: string;
}) {
  toast.error(`API Error: ${err.type} (${err.code}): ${err.message}`);
  console.error("API Error", err);
}

const API_CLIENT_KEY = Symbol("API_CLIENT");

export function setApiClient(baseUrl: string, token?: string) {
  const apiClient = new ApiClient(baseUrl);
  apiClient.setToken(token);
  return setContext(API_CLIENT_KEY, apiClient);
}

export function getApiClient() {
  return getContext<ReturnType<typeof setApiClient>>(API_CLIENT_KEY);
}

export function openConfirm(props: ConfirmModalProps) {
  return modals.open(ConfirmModal, props);
}

export function openInput(props: InputModalProps) {
  return modals.open(InputModal, props);
}
