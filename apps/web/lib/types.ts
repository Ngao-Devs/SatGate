export type SatGoForm = {
  id: string;
  name: string;
  domain: string;
  amount_sats: number;
  created_at: string;
};

export type SatGoMessage = {
  id: string;
  form_id: string;
  sender_name: string;
  sender_email: string;
  body: string;
  status: string;
  amount_sats: number;
  created_at: string;
  paid_at: string | null;
};

export type InvoiceResponse = {
  invoice_id: string;
  message_id: string;
  payment_request: string;
  amount_sats: number;
  status: string;
  expires_at: string;
};

export type InvoiceStatusResponse = {
  invoice_id: string;
  message_id: string;
  status: string;
  paid_at: string | null;
};
