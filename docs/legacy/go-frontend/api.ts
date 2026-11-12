import type { InvoiceResponse, InvoiceStatus, SatGateForm, SatGateMessage } from "./types";

export const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "https://satgate.onrender.com";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_URL}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...init?.headers,
    },
  });

  if (!response.ok) {
    const detail = await response.json().catch(() => ({ error: response.statusText }));
    throw new Error(detail.error ?? "SatGate request failed");
  }

  return response.json() as Promise<T>;
}

export function getForms() {
  return request<SatGateForm[]>("/api/forms", { cache: "no-store" });
}

export function createForm(payload: { name: string; domain: string; amount_sats: number }) {
  return request<SatGateForm>("/api/forms", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export async function deleteForm(id: string): Promise<void> {
  const response = await fetch(`${API_URL}/api/forms/${id}`, {
    method: "DELETE",
    headers: { "Content-Type": "application/json" },
  });
  if (!response.ok) {
    const detail = await response.json().catch(() => ({ error: response.statusText }));
    throw new Error(detail.error ?? "SatGate request failed");
  }
}

export function getMessages() {
  return request<SatGateMessage[]>("/api/messages", { cache: "no-store" });
}

export function createInvoice(payload: {
  form_id: string;
  sender_name: string;
  sender_email: string;
  body: string;
}) {
  return request<InvoiceResponse>("/api/invoices", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function getInvoice(invoiceId: string) {
  return request<InvoiceStatus>(`/api/invoices/${invoiceId}`, { cache: "no-store" });
}

export function mockPayInvoice(invoiceId: string) {
  return request<InvoiceStatus>(`/api/invoices/${invoiceId}/mock-pay`, {
    method: "POST",
  });
}


export async function updateForm(
  id: string,
  data: { name: string; domain: string; amount_sats: number }
): Promise<SatGateForm> {
  const res = await fetch(`${API_URL}/api/forms/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}