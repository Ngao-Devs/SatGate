import type {
  SatGoForm,
  SatGoMessage,
  InvoiceResponse,
  InvoiceStatusResponse,
} from "./types";

const API_BASE_URL = `${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api`;

export const API_URL = API_BASE_URL;

type CreateInvoicePayload = {
  sender_name: string;
  sender_email: string;
  body: string;
};

function getAuthHeaders(): Record<string, string> {
  if (typeof window === "undefined") {
    return {};
  }

  const token = localStorage.getItem("satgo_token");

  if (token) {
    return { Authorization: `Bearer ${token}` };
  }

  if (process.env.NODE_ENV === "development") {
    return { Authorization: "Bearer mock-owner-demo-user" };
  }

  return {};
}

async function getErrorMessage(response: Response, fallback: string): Promise<string> {
  const body = await response.json().catch(() => null);

  if (body && typeof body.error === "string") {
    return body.error;
  }

  return fallback;
}

export async function createInvoice(
  formId: string,
  payload: CreateInvoicePayload,
): Promise<InvoiceResponse> {
  const response = await fetch(`${API_BASE_URL}/invoices`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      form_id: formId,
      sender_name: payload.sender_name,
      sender_email: payload.sender_email,
      body: payload.body,
    }),
  });

  if (!response.ok) {
    throw new Error(await getErrorMessage(response, "Failed to create invoice"));
  }

  return response.json();
}

export async function checkInvoice(invoiceId: string): Promise<InvoiceStatusResponse> {
  const response = await fetch(`${API_BASE_URL}/invoices/${encodeURIComponent(invoiceId)}`);

  if (!response.ok) {
    throw new Error(await getErrorMessage(response, "Failed to check invoice status"));
  }

  return response.json();
}

export async function getForms(): Promise<SatGoForm[]> {
  const response = await fetch(`${API_BASE_URL}/forms`, {
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    throw new Error(await getErrorMessage(response, "Failed to load forms"));
  }

  return response.json();
}

export async function getMessages(): Promise<SatGoMessage[]> {
  const response = await fetch(`${API_BASE_URL}/messages`, {
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    throw new Error(await getErrorMessage(response, "Failed to load messages"));
  }

  return response.json();
}

export async function createForm(data: {
  name: string;
  domain: string;
  amount_sats: number;
}): Promise<SatGoForm> {
  const response = await fetch(`${API_BASE_URL}/forms`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...getAuthHeaders(),
    },
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    throw new Error(await getErrorMessage(response, "Failed to create form"));
  }

  return response.json();
}

export async function updateForm(
  id: string,
  data: {
    name: string;
    domain: string;
    amount_sats: number;
  },
): Promise<SatGoForm> {
  const response = await fetch(`${API_BASE_URL}/forms/${encodeURIComponent(id)}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
      ...getAuthHeaders(),
    },
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    throw new Error(await getErrorMessage(response, "Failed to update form"));
  }

  return response.json();
}

export async function deleteForm(id: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/forms/${encodeURIComponent(id)}`, {
    method: "DELETE",
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    throw new Error(await getErrorMessage(response, "Failed to delete form"));
  }
}
