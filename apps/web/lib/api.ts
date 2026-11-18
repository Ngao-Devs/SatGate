const API_BASE_URL = (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080") + "/api";

export const API_URL = API_BASE_URL;

function getAuthHeaders() {
  const headers: Record<string, string> = {};
  if (typeof window !== "undefined") {
    const token = localStorage.getItem("satgo_token");
    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    } else {
      // In development / demo, fallback to a mock token to avoid setup friction
      headers["Authorization"] = "Bearer mock-owner-demo-user";
    }
  } else {
    headers["Authorization"] = "Bearer mock-owner-demo-user";
  }
  return headers;
}

export async function createInvoice(site_id: string, message: string, amount_sats: number) {
  const response = await fetch(`${API_BASE_URL}/invoices`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ site_id, message, amount_sats }),
  });

  if (!response.ok) {
    throw new Error("Failed to create invoice");
  }

  return response.json();
}

export async function checkInvoice(payment_hash: string) {
  const response = await fetch(`${API_BASE_URL}/invoices/${payment_hash}`);
  if (!response.ok) {
    throw new Error("Failed to check invoice status");
  }
  return response.json();
}

export async function getSite(site_id: string) {
  const response = await fetch(`${API_BASE_URL}/sites/${site_id}`, {
    headers: getAuthHeaders(),
  });
  if (!response.ok) {
    throw new Error("Site not found");
  }
  return response.json();
}

export async function getSubmissions(site_id: string) {
  const response = await fetch(`${API_BASE_URL}/sites/${site_id}/submissions`, {
    headers: getAuthHeaders(),
  });
  if (!response.ok) {
    throw new Error("Failed to load submissions");
  }
  return response.json();
}

// Aliases for the Dashboard UI
export async function getForms() {
  const response = await fetch(`${API_BASE_URL}/sites`, {
    headers: getAuthHeaders(),
  });
  if (!response.ok) {
    throw new Error("Failed to load forms");
  }
  return response.json();
}

export async function getMessages() {
  const response = await fetch(`${API_BASE_URL}/submissions`, {
    headers: getAuthHeaders(),
  });
  if (!response.ok) {
    throw new Error("Failed to load messages");
  }
  return response.json();
}

export async function createForm(data: { name: string; domain: string; amount_sats: number }) {
  const response = await fetch(`${API_BASE_URL}/sites`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...getAuthHeaders(),
    },
    body: JSON.stringify({
      owner_id: "demo-user",
      name: data.name,
      domain: data.domain,
      min_amount_sats: data.amount_sats,
    }),
  });

  if (!response.ok) {
    throw new Error("Failed to create form");
  }

  return response.json();
}

export async function updateForm(
  id: string,
  data: { name: string; domain: string; amount_sats: number },
) {
  const response = await fetch(`${API_BASE_URL}/sites/${id}`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
      ...getAuthHeaders(),
    },
    body: JSON.stringify({
      name: data.name,
      domain: data.domain,
      min_amount_sats: data.amount_sats,
    }),
  });

  if (!response.ok) {
    throw new Error("Failed to update form");
  }

  return response.json();
}

export async function deleteForm(id: string) {
  const response = await fetch(`${API_BASE_URL}/sites/${id}`, {
    method: "DELETE",
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    throw new Error("Failed to delete form");
  }
}

export async function createSite(owner_id: string, name: string) {
  const response = await fetch(`${API_BASE_URL}/sites`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...getAuthHeaders(),
    },
    body: JSON.stringify({ owner_id, name }),
  });

  if (!response.ok) {
    throw new Error("Failed to create site");
  }

  return response.json();
}

export async function withdrawSats(site_id: string, invoice: string) {
  const response = await fetch(`${API_BASE_URL}/withdraw`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...getAuthHeaders(),
    },
    body: JSON.stringify({ site_id, invoice }),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || "Failed to process withdrawal");
  }

  return response.json();
}
