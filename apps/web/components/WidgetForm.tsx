"use client";

import { FormEvent, useEffect, useState, useRef } from "react";
import QRCode from "qrcode";
import { AlertTriangle, CheckCircle, Copy, Loader2, Send, Wallet, CreditCard } from "lucide-react";
import { createInvoice, checkInvoice, API_URL } from "@/lib/api";
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
  const [message, setMessage] = useState("");
  const [invoice, setInvoice] = useState<InvoiceResponse | null>(null);
  const [qrDataUrl, setQrDataUrl] = useState("");
  const [status, setStatus] = useState<"idle" | "creating" | "pending" | "paid" | "failed">("idle");
  const [error, setError] = useState("");
  const [copied, setCopied] = useState(false);

  const isMock = process.env.NEXT_PUBLIC_LN_MODE === "mock";
  const pollingRef = useRef<NodeJS.Timeout | null>(null);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setStatus("creating");

    try {
      const nextInvoice = await createInvoice(formId, {
        sender_name: senderName,
        sender_email: senderEmail,
        body: message,
      });

      setInvoice(nextInvoice);
      setStatus("pending");

      const nextQr = await QRCode.toDataURL(nextInvoice.payment_request, {
        margin: 2,
        width: 300,
        errorCorrectionLevel: "M",
        color: {
          dark: "#000000",
          light: "#ffffff",
        },
      });

      setQrDataUrl(nextQr);
      startPolling(nextInvoice.invoice_id);

      try {
        await payWithWebLn(nextInvoice.payment_request);
      } catch (paymentError) {
        console.log("WebLN auto-pay skipped or failed", paymentError);
      }
    } catch (submitError) {
      setStatus("failed");
      setError(submitError instanceof Error ? submitError.message : "Could not create invoice");
    }
  }

  function startPolling(invoiceId: string) {
    if (pollingRef.current) {
      clearInterval(pollingRef.current);
    }

    pollingRef.current = setInterval(async () => {
      try {
        const invoiceStatus = await checkInvoice(invoiceId);

        if (invoiceStatus.status === "paid" || invoiceStatus.paid_at !== null) {
          setStatus("paid");

          if (pollingRef.current) {
            clearInterval(pollingRef.current);
            pollingRef.current = null;
          }
        }
      } catch (pollingError) {
        console.error("Polling error", pollingError);
      }
    }, 2000);
  }

  async function handleSimulatePayment() {
    if (!invoice) return;

    setError("");

    try {
      const response = await fetch(
        `${API_URL}/invoices/${encodeURIComponent(invoice.invoice_id)}/mock-pay`,
        { method: "POST" },
      );

      if (!response.ok) {
        const body = await response.json().catch(() => null);
        throw new Error(body?.error || "Failed to simulate payment");
      }

      const result = await response.json();

      if (result.status === "paid" || result.paid_at) {
        setStatus("paid");

        if (pollingRef.current) {
          clearInterval(pollingRef.current);
          pollingRef.current = null;
        }
      }
    } catch (paymentError) {
      setError(paymentError instanceof Error ? paymentError.message : "Failed to simulate payment");
    }
  }

  useEffect(() => {
    return () => {
      if (pollingRef.current) {
        clearInterval(pollingRef.current);
      }
    };
  }, []);

  const copyToClipboard = async () => {
    if (!invoice) return;

    try {
      await navigator.clipboard.writeText(invoice.payment_request);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (copyError) {
      console.error("Failed to copy invoice", copyError);
    }
  };

  if (status === "paid") {
    return (
      <div className="flex flex-col items-center justify-center p-12 text-center space-y-4 bg-white min-h-[400px]">
        <div className="w-16 h-16 bg-green-50 rounded-full flex items-center justify-center">
          <CheckCircle className="w-10 h-10 text-satGreen" />
        </div>
        <h2 className="text-2xl font-bold text-satBlack brand-font">Message Sent!</h2>
        <p className="text-slate-500 max-w-[240px]">
          Your payment is confirmed and your message has been sent.
        </p>
        <button
          onClick={() => window.location.reload()}
          className="text-satBlue font-semibold hover:underline mt-4"
        >
          Send another message
        </button>
      </div>
    );
  }

  return (
    <div className="max-w-md mx-auto bg-white border border-slate-200 rounded-xl overflow-hidden font-body">
      <div className="p-6 space-y-6">
        <div className="flex items-center justify-between">
          <h1 className="text-lg font-bold text-satBlack brand-font">SatGo</h1>
          <StatusBadge status={status} />
        </div>

        {status === "idle" || status === "creating" || (status === "failed" && !invoice) ? (
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-1.5">
              <label className="text-sm font-bold text-satBlack">Your name</label>
              <input
                type="text"
                required
                value={senderName}
                onChange={(event) => setSenderName(event.target.value)}
                placeholder="Your name"
                className="w-full bg-white border border-slate-200 rounded-lg p-3 text-satBlack placeholder:text-slate-400 focus:ring-2 focus:ring-satBlue/20 focus:border-satBlue outline-none transition-all shadow-sm"
              />
            </div>

            <div className="space-y-1.5">
              <label className="text-sm font-bold text-satBlack">Your email</label>
              <input
                type="email"
                required
                value={senderEmail}
                onChange={(event) => setSenderEmail(event.target.value)}
                placeholder="you@example.com"
                className="w-full bg-white border border-slate-200 rounded-lg p-3 text-satBlack placeholder:text-slate-400 focus:ring-2 focus:ring-satBlue/20 focus:border-satBlue outline-none transition-all shadow-sm"
              />
            </div>

            <div className="space-y-1.5">
              <label className="text-sm font-bold text-satBlack">Message</label>
              <textarea
                required
                minLength={10}
                value={message}
                onChange={(event) => setMessage(event.target.value)}
                placeholder="What's on your mind?"
                className="w-full h-32 bg-white border border-slate-200 rounded-lg p-3 text-satBlack placeholder:text-slate-400 focus:ring-2 focus:ring-satBlue/20 focus:border-satBlue outline-none transition-all resize-none shadow-sm"
              />
              <p className="text-xs text-slate-400">
                Your message must contain at least 10 characters.
              </p>
            </div>

            <div className="flex items-center justify-between p-4 bg-slate-50 rounded-lg border border-slate-100">
              <span className="text-sm font-semibold text-slate-500">Amount</span>
              <span className="font-bold text-satBlue">Set by form owner</span>
            </div>

            <PrimaryButton
              type="submit"
              disabled={status === "creating"}
              className="w-full shadow-sm"
            >
              {status === "creating" ? (
                <>
                  <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                  Generating Invoice...
                </>
              ) : (
                <>
                  <Send className="w-4 h-4 mr-2" />
                  Pay & Submit
                </>
              )}
            </PrimaryButton>
          </form>
        ) : (
          <div className="space-y-6 animate-in slide-in-from-bottom-2 duration-300">
            <div className="flex flex-col items-center space-y-4">
              <div className="text-center">
                <p className="text-sm font-semibold text-slate-500">Amount due</p>
                <p className="text-2xl font-bold text-satBlue">{invoice?.amount_sats} SATS</p>
              </div>

              {qrDataUrl && (
                <div className="p-3 bg-white border border-slate-100 rounded-xl shadow-sm">
                  <img src={qrDataUrl} alt="Lightning Invoice" className="w-44 h-44" />
                </div>
              )}

              <div className="w-full space-y-3">
                <button
                  onClick={copyToClipboard}
                  className="w-full flex items-center justify-between p-3 bg-slate-50 hover:bg-slate-100 border border-slate-200 rounded-lg transition-colors group"
                >
                  <span className="text-xs font-mono text-slate-500 truncate mr-4">
                    {invoice?.payment_request}
                  </span>
                  {copied ? (
                    <CheckCircle className="w-4 h-4 text-satGreen" />
                  ) : (
                    <Copy className="w-4 h-4 text-slate-400 group-hover:text-slate-600" />
                  )}
                </button>

                <div className="grid grid-cols-2 gap-3">
                  <button
                    onClick={() => payWithWebLn(invoice?.payment_request || "")}
                    className="flex items-center justify-center p-3 bg-satBlue hover:bg-blue-700 text-white rounded-lg font-bold transition-colors shadow-sm"
                  >
                    <Wallet className="w-4 h-4 mr-2" />
                    WebLN
                  </button>
                  <button
                    onClick={() => window.open(`lightning:${invoice?.payment_request}`, "_blank")}
                    className="flex items-center justify-center p-3 bg-white hover:bg-slate-50 text-satBlack border border-slate-200 rounded-lg font-bold transition-colors shadow-sm"
                  >
                    Open Wallet
                  </button>
                </div>

                {isMock && (
                  <button
                    onClick={handleSimulatePayment}
                    className="w-full flex items-center justify-center p-2 text-xs font-bold text-amber-700 bg-amber-50 hover:bg-amber-100 border border-amber-200 rounded-lg transition-colors mt-2"
                  >
                    <CreditCard className="w-3 h-3 mr-2" />
                    Simulate Demo Payment
                  </button>
                )}
              </div>
            </div>

            <div className="p-4 bg-satBlue/5 border border-satBlue/10 rounded-lg">
              <p className="text-xs text-slate-600 leading-relaxed text-center italic">
                Waiting for network confirmation...
              </p>
            </div>
          </div>
        )}

        {error && (
          <div className="flex items-start space-x-2 p-3 bg-red-50 border border-red-100 rounded-lg text-satRed text-xs font-semibold">
            <AlertTriangle className="w-4 h-4 shrink-0" />
            <p>{error}</p>
          </div>
        )}
      </div>
    </div>
  );
}
