"use client";

import { FormEvent, useEffect, useState } from "react";
import QRCode from "qrcode";
import { AlertTriangle, CheckCircle, Copy, Loader2, Send, Wallet } from "lucide-react";
import { createInvoice, getInvoice, mockPayInvoice } from "@/lib/api";
import type { InvoiceResponse } from "@/lib/types";
import { payWithWebLn } from "@/lib/webln";
import { PrimaryButton } from "./PrimaryButton";
import { StatusBadge } from "./StatusBadge";

type WidgetFormProps = {
  formId: string;
};

export function WidgetForm({ formId }: WidgetFormProps) {
  const [senderName, setSenderName] = useState("");
  const [senderEmail, setSenderEmail] = useState("");
  const [body, setBody] = useState("");
  const [invoice, setInvoice] = useState<InvoiceResponse | null>(null);
  const [qrDataUrl, setQrDataUrl] = useState("");
  const [status, setStatus] = useState<"idle" | "creating" | "pending" | "paid" | "failed">("idle");
  const [error, setError] = useState("");
  const [copied, setCopied] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

      console.log("FORM ID:", formId);

    console.log("submit fired")
    setError("");
    setStatus("creating");

    try {
      const nextInvoice = await createInvoice({
        form_id: formId,
        sender_name: senderName,
        sender_email: senderEmail,
        body,
      });
      setInvoice(nextInvoice);
      console.log("INVOICE FROM BACKEND:", nextInvoice);
console.log("PAYMENT REQUEST:", nextInvoice.payment_request);
      setStatus("pending");


      const nextQr = await QRCode.toDataURL(nextInvoice.payment_request, {
        margin: 1,
        width: 220,
        color: {
          dark: "#0b0f14",
          light: "#ffffff",
        },
      });
      setQrDataUrl(nextQr);
    } catch (submitError) {
      setStatus("failed");
      setError(submitError instanceof Error ? submitError.message : "Could not create invoice");
    }
  }

  async function handleWebLnPay() {
    if (!invoice) return;

    setError("");

    try {
      await payWithWebLn(invoice.payment_request);
      await pollUntilPaid(invoice.invoice_id);
    } catch (payError) {
      setError(
        payError instanceof Error
          ? `${payError.message}. Use the QR code or the demo payment button.`
          : "WebLN payment failed. Use the QR code or the demo payment button.",
      );
    }
  }

  async function handleMockPay() {
    if (!invoice) return;

    setError("");

    try {
      const paid = await mockPayInvoice(invoice.invoice_id);
      if (paid.status === "paid") {
        setStatus("paid");
      }
    } catch (mockError) {
      setError(mockError instanceof Error ? mockError.message : "Could not settle demo invoice");
    }
  }

  async function pollUntilPaid(invoiceId: string) {
    for (let attempt = 0; attempt < 10; attempt += 1) {
      const nextStatus = await getInvoice(invoiceId);

      if (nextStatus.status === "paid") {
        setStatus("paid");
        return;
      }

      await new Promise((resolve) => window.setTimeout(resolve, 1200));
    }
  }

  async function copyInvoice() {
    if (!invoice) return;

    await navigator.clipboard.writeText(invoice.payment_request);
    setCopied(true);
    window.setTimeout(() => setCopied(false), 1800);
  }

  useEffect(() => {
    if (status === "paid") {
      setSenderName("");
      setSenderEmail("");
      setBody("");
    }
  }, [status]);

  return (
    <div className="min-h-screen bg-white text-satBlack">
      <div className="mx-auto w-full max-w-xl px-4 py-5">
        <header className="flex items-center justify-between border-b border-slate-200 pb-3">
          <span className="brand-font text-xl font-bold">SatGate</span>
          <span className="rounded-full bg-blue-50 px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider text-satBlue">
            Protected Form
          </span>
        </header>

        {status === "paid" ? (
          <section className="mt-5 rounded-md border border-green-200 bg-green-50 p-5 text-satGreen">
            <div className="flex items-center gap-2">
              <CheckCircle size={22} aria-hidden="true" />
              <h2 className="brand-font text-lg font-bold">Message verified</h2>
            </div>
            <p className="mt-2 text-sm leading-6">
              Your paid message is now visible in the SatGate dashboard.
            </p>
            <PrimaryButton className="mt-4" onClick={() => {
              setInvoice(null);
              
              setQrDataUrl("");
              setStatus("idle");
            }}>
              Send another
            </PrimaryButton>
          </section>
        ) : (
          <form className="mt-5 space-y-4" onSubmit={handleSubmit}>
            <label className="block text-sm font-semibold">
              Name
              <input
                className="mt-1 h-11 w-full rounded-md border border-slate-300 px-3 text-sm outline-none transition focus:border-satBlue focus:ring-2 focus:ring-blue-100"
                value={senderName}
                onChange={(event) => setSenderName(event.target.value)}
                disabled={status === "creating" || status === "pending"}
                required
              />
            </label>
            <label className="block text-sm font-semibold">
              Email
              <input
                className="mt-1 h-11 w-full rounded-md border border-slate-300 px-3 text-sm outline-none transition focus:border-satBlue focus:ring-2 focus:ring-blue-100"
                type="email"
                value={senderEmail}
                onChange={(event) => setSenderEmail(event.target.value)}
                disabled={status === "creating" || status === "pending"}
                required
              />
            </label>
            <label className="block text-sm font-semibold">
              Message
              <textarea
                className="mt-1 min-h-32 w-full resize-y rounded-md border border-slate-300 px-3 py-2 text-sm leading-6 outline-none transition focus:border-satBlue focus:ring-2 focus:ring-blue-100"
                value={body}
                onChange={(event) => setBody(event.target.value)}
                disabled={status === "creating" || status === "pending"}
                required
              />
            </label>

            {error ? (
              <div className="flex gap-2 rounded-md border border-yellow-200 bg-yellow-50 p-3 text-sm font-semibold text-amber-800">
                <AlertTriangle className="mt-0.5 shrink-0" size={18} aria-hidden="true" />
                <span>{error}</span>
              </div>
            ) : null}

            {!invoice ? (
              <PrimaryButton className="w-full" disabled={status === "creating"} type="submit">
                {status === "creating" ? <Loader2 className="animate-spin" size={17} aria-hidden="true" /> : <Send size={17} aria-hidden="true" />}
                Request invoice
              </PrimaryButton>
            ) : (
              <PaymentPanel
                copied={copied}
                invoice={invoice}
                qrDataUrl={qrDataUrl}
                status={status}
                onCopy={copyInvoice}
                onMockPay={handleMockPay}
                onWebLnPay={handleWebLnPay}
              />
            )}
          </form>
        )}
      </div>
    </div>
  );
}

function PaymentPanel({
  copied,
  invoice,
  qrDataUrl,
  status,
  onCopy,
  onMockPay,
  onWebLnPay,
}: {
  copied: boolean;
  invoice: InvoiceResponse;
  qrDataUrl: string;
  status: string;
  onCopy: () => void;
  onMockPay: () => void;
  onWebLnPay: () => void;
}) {
  return (
    <section className="rounded-md border border-slate-200 bg-slate-50 p-4">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h2 className="brand-font text-lg font-bold">Pay {invoice.amount_sats} sats</h2>
          <p className="mt-1 text-sm text-slate-600">Invoice expires at {new Date(invoice.expires_at).toLocaleTimeString()}.</p>
        </div>
        <StatusBadge status={status} />
      </div>

      {qrDataUrl ? (
        <div className="mt-4 flex justify-center">
          <img
            alt="Lightning invoice QR code"
            className="h-48 w-48 rounded-md border border-slate-200 bg-white p-2"
            height={192}
            src={qrDataUrl}
            width={192}
          />
        </div>
      ) : null}

      <p className="mt-4 break-all rounded-md bg-white p-3 text-xs leading-5 text-slate-700">
        {invoice.payment_request}
      </p>

      <div className="mt-4 grid gap-2 sm:grid-cols-3">
        <PrimaryButton onClick={onWebLnPay} type="button">
          <Wallet size={17} aria-hidden="true" />
          WebLN
        </PrimaryButton>
        <PrimaryButton onClick={onCopy} tone="neutral" type="button">
          <Copy size={17} aria-hidden="true" />
          {copied ? "Copied" : "Copy"}
        </PrimaryButton>
        <PrimaryButton onClick={onMockPay} tone="success" type="button">
          Demo pay
        </PrimaryButton>
      </div>
    </section>
  );
}
