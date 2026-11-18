"use client";

import { FormEvent, useEffect, useMemo, useState } from "react";
import Link from "next/link";
import {
  Check,
  Copy,
  ExternalLink,
  Inbox,
  Pencil,
  RefreshCw,
  ShieldCheck,
  Trash2,
  X,
  Globe,
  Coins,
  FileText,
  PlusCircle,
  Code,
  LayoutGrid,
  Zap,
} from "lucide-react";
import {
  API_URL,
  createForm,
  deleteForm,
  getForms,
  getMessages,
  updateForm,
  withdrawSats,
} from "@/lib/api";
import type { SatGoForm, SatGoMessage } from "@/lib/types";
import { PrimaryButton } from "./PrimaryButton";
import { StatusBadge } from "./StatusBadge";

function BitcoinPulseLoader() {
  return (
    <main className="mx-auto flex min-h-screen w-full max-w-7xl flex-col gap-6 px-6 py-12 sm:px-8">
      <div className="flex min-h-[60vh] flex-col items-center justify-center gap-6">
        <div className="relative flex items-center justify-center">
          <div
            style={{
              position: "absolute",
              width: "100px",
              height: "100px",
              borderRadius: "50%",
              background: "rgba(29,78,216,0.12)",
              animation: "ping 1.4s ease-out infinite",
            }}
          />
          <div
            style={{
              position: "absolute",
              width: "80px",
              height: "80px",
              borderRadius: "50%",
              background: "rgba(29,78,216,0.18)",
              animation: "ping 1.4s ease-out infinite 0.3s",
            }}
          />
          <div
            style={{
              width: "64px",
              height: "64px",
              borderRadius: "50%",
              background: "linear-gradient(135deg, #1d4ed8, #3b82f6)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              boxShadow: "0 0 30px rgba(29,78,216,0.45)",
              animation: "glow 1.4s ease-in-out infinite alternate",
              zIndex: 1,
            }}
          >
            <span
              style={{
                color: "white",
                fontSize: "30px",
                fontWeight: "bold",
                fontFamily: "var(--font-brand)",
                lineHeight: 1,
                marginTop: "-2px",
              }}
            >
              ₿
            </span>
          </div>
        </div>
        <div className="flex flex-col items-center gap-1.5 mt-2">
          <p className="brand-font text-xl font-bold text-satBlack">SatGo</p>
          <p className="text-sm text-slate-500">Loading your dashboard...</p>
        </div>
        <style>{`
          @keyframes ping {
            0% { transform: scale(0.8); opacity: 0.8; }
            100% { transform: scale(2); opacity: 0; }
          }
          @keyframes glow {
            from { box-shadow: 0 0 12px rgba(29,78,216,0.4); }
            to { box-shadow: 0 0 28px rgba(29,78,216,0.9), 0 0 48px rgba(59,130,246,0.3); }
          }
          @keyframes shimmer {
            0% { background-position: -600px 0; }
            100% { background-position: 600px 0; }
          }
          @keyframes progress {
            0% { width: 0%; }
            20% { width: 25%; }
            50% { width: 60%; }
            80% { width: 85%; }
            95% { width: 95%; }
          }
        `}</style>
      </div>
    </main>
  );
}

function SkeletonCard() {
  const shimmer = {
    background: "linear-gradient(90deg, #f1f5f9 25%, #e2e8f0 50%, #f1f5f9 75%)",
    backgroundSize: "600px 100%",
    animation: "shimmer 1.4s infinite linear",
  };
  return (
    <div className="rounded-2xl border border-slate-200/60 bg-white p-8 shadow-panel">
      <div
        style={{
          ...shimmer,
          height: "14px",
          width: "50%",
          borderRadius: "6px",
          marginBottom: "12px",
        }}
      />
      <div style={{ ...shimmer, height: "36px", width: "35%", borderRadius: "6px" }} />
    </div>
  );
}

function SkeletonMessage() {
  const shimmer = {
    background: "linear-gradient(90deg, #f1f5f9 25%, #e2e8f0 50%, #f1f5f9 75%)",
    backgroundSize: "600px 100%",
    animation: "shimmer 1.4s infinite linear",
  };
  return (
    <div className="rounded-xl border border-slate-200 bg-white p-5">
      <div className="flex items-start justify-between">
        <div style={{ flex: 1 }}>
          <div
            style={{
              ...shimmer,
              height: "16px",
              width: "30%",
              borderRadius: "6px",
              marginBottom: "8px",
            }}
          />
          <div style={{ ...shimmer, height: "13px", width: "45%", borderRadius: "6px" }} />
        </div>
        <div style={{ ...shimmer, height: "24px", width: "70px", borderRadius: "20px" }} />
      </div>
      <div
        style={{
          ...shimmer,
          height: "13px",
          width: "100%",
          borderRadius: "6px",
          marginTop: "16px",
        }}
      />
      <div
        style={{ ...shimmer, height: "13px", width: "75%", borderRadius: "6px", marginTop: "8px" }}
      />
    </div>
  );
}

function FormProgressBar({ active }: { active: boolean }) {
  if (!active) return null;
  return (
    <div
      style={{
        position: "absolute",
        top: 0,
        left: 0,
        right: 0,
        height: "4px",
        background: "#e2e8f0",
        borderRadius: "4px 4px 0 0",
        overflow: "hidden",
      }}
    >
      <div
        style={{
          height: "100%",
          background: "linear-gradient(90deg, #1d4ed8, #3b82f6)",
          animation: "progress 2.5s ease-out forwards",
          boxShadow: "0 0 8px rgba(29,78,216,0.6)",
        }}
      />
    </div>
  );
}

function Metric({ label, value, icon }: { label: string; value: string; icon: React.ReactNode }) {
  return (
    <div className="rounded-2xl border border-slate-200/60 bg-white p-6 md:p-8 shadow-panel hover:shadow-lg transition-all duration-300 flex items-center justify-between group">
      <div>
        <p className="text-sm font-semibold text-slate-500 tracking-wide">{label}</p>
        <p className="brand-font mt-2.5 text-3xl md:text-4xl font-bold text-satBlack tracking-tight">
          {value}
        </p>
      </div>
      <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-slate-50 text-slate-400 group-hover:bg-blue-50 group-hover:text-satBlue transition duration-300">
        {icon}
      </div>
    </div>
  );
}

interface EditState {
  id: string;
  name: string;
  domain: string;
  amount_sats: number;
}

export function DashboardClient() {
  const [forms, setForms] = useState<SatGoForm[]>([]);
  const [messages, setMessages] = useState<SatGoMessage[]>([]);
  const [name, setName] = useState("Portfolio Contact");
  const [domain, setDomain] = useState("localhost");
  const [amount, setAmount] = useState(5);
  const [error, setError] = useState("");
  const [copied, setCopied] = useState(false);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [selectedFormId, setSelectedFormId] = useState<string | null>(null);
  const [editState, setEditState] = useState<EditState | null>(null);
  const [saving, setSaving] = useState(false);

  const [withdrawInvoice, setWithdrawInvoice] = useState("");
  const [withdrawing, setWithdrawing] = useState(false);
  const [withdrawStatus, setWithdrawStatus] = useState<{
    success: boolean;
    message: string;
  } | null>(null);

  const handleWithdraw = async () => {
    if (!selectedForm || !withdrawInvoice) return;
    setWithdrawing(true);
    setWithdrawStatus(null);
    try {
      const res = await withdrawSats(selectedForm.id, withdrawInvoice);
      setWithdrawStatus({
        success: true,
        message: `Successfully withdrew ${res.amount_paid_sats} sats! Preimage: ${res.payment_preimage.slice(0, 16)}...`,
      });
      setWithdrawInvoice("");
      await loadData();
    } catch (e) {
      setWithdrawStatus({
        success: false,
        message: e instanceof Error ? e.message : "Withdrawal failed.",
      });
    } finally {
      setWithdrawing(false);
    }
  };

  useEffect(() => {
    async function init() {
      setLoading(true);
      try {
        const [f, m] = await Promise.all([getForms(), getMessages()]);
        setForms(f);
        setMessages(m);
        if (f.length > 0) setSelectedFormId(f[0].id);
      } catch (e) {
        console.error(e);
        setError("Failed to load data");
      } finally {
        setLoading(false);
      }
    }
    init();
  }, []);

  const selectedForm = forms.find((f) => f.id === (selectedFormId || forms[0]?.id));

  const filteredMessages = useMemo(() => {
    if (!selectedForm) return [];
    return messages.filter((m) => m.form_id === selectedForm.id);
  }, [messages, selectedForm]);

  const totalSats = filteredMessages.reduce((sum, m) => sum + m.amount_sats, 0);

  const embedUrl =
    typeof window !== "undefined" && selectedForm
      ? `${window.location.origin}/embed/${selectedForm.id}`
      : "";

  const embedSnippet = selectedForm
    ? `<iframe src="${embedUrl}" title="SatGo contact form" width="100%" height="620" style="border:0"></iframe>`
    : "";

  async function loadData() {
    setLoading(true);
    try {
      const [f, m] = await Promise.all([getForms(), getMessages()]);
      setForms(f);
      setMessages(m);
    } catch {
      setError("Refresh failed");
    } finally {
      setLoading(false);
    }
  }

  async function handleCreateForm(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSubmitting(true);
    try {
      const newForm = await createForm({ name, domain, amount_sats: amount });
      setForms((prev) => [newForm, ...prev]);
      setSelectedFormId(newForm.id);
    } catch {
      setError("Form creation failed");
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDeleteForm(id: string) {
    if (!confirm("Are you sure? All associated messages will be deleted.")) return;
    try {
      await deleteForm(id);
      setForms((prev) => prev.filter((f) => f.id !== id));
      if (selectedFormId === id) setSelectedFormId(forms.find((f) => f.id !== id)?.id || null);
    } catch {
      setError("Deletion failed");
    }
  }

  function startEdit(form: SatGoForm) {
    setEditState({
      id: form.id,
      name: form.name,
      domain: form.domain,
      amount_sats: form.amount_sats,
    });
  }

  function cancelEdit() {
    setEditState(null);
  }

  async function handleSaveEdit() {
    if (!editState) return;
    setSaving(true);
    setError("");
    try {
      const updated = await updateForm(editState.id, {
        name: editState.name,
        domain: editState.domain,
        amount_sats: editState.amount_sats,
      });
      setForms((prev) => prev.map((f) => (f.id === updated.id ? updated : f)));
      setEditState(null);
    } catch {
      setError("Update failed");
    } finally {
      setSaving(false);
    }
  }

  async function copyEmbed() {
    if (!embedSnippet) return;
    await navigator.clipboard.writeText(embedSnippet);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }

  if (loading && forms.length === 0) {
    return <BitcoinPulseLoader />;
  }

  return (
    <main className="mx-auto flex min-h-screen w-full max-w-7xl flex-col gap-8 px-6 py-8 sm:py-12">
      <header className="flex flex-col gap-5 border-b border-slate-200/80 pb-6 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <p className="text-xs font-bold uppercase tracking-wider text-satBlue">
            Economic proof of intent
          </p>
          <h1 className="brand-font mt-2.5 text-3xl sm:text-4xl font-extrabold text-satBlack tracking-tight">
            SatGo Dashboard
          </h1>
          <p className="mt-2 max-w-2xl text-base text-slate-500 leading-relaxed">
            Lightning-powered contact forms that make abuse costly and keep verified messages easy
            to read.
          </p>
        </div>
        <div className="flex flex-wrap gap-3">
          <Link
            className="inline-flex h-11 items-center justify-center gap-2 rounded-xl border border-slate-200 bg-white px-5 text-sm font-semibold text-slate-700 shadow-sm transition hover:bg-slate-50 hover:text-satBlue hover:border-slate-300"
            href="/demo"
          >
            <ExternalLink size={16} aria-hidden="true" />
            Demo site
          </Link>
          <PrimaryButton onClick={loadData} disabled={loading} className="h-11 rounded-xl">
            <RefreshCw className={loading ? "animate-spin" : ""} size={16} aria-hidden="true" />
            Refresh
          </PrimaryButton>
        </div>
      </header>

      {error && (
        <section className="rounded-xl border border-red-200 bg-red-50 p-4 text-sm font-bold text-satRed flex items-center gap-2">
          <span className="w-2 h-2 rounded-full bg-satRed animate-pulse" />
          {error}
        </section>
      )}

      {/* Metrics Section */}
      <section className="grid gap-6 md:grid-cols-3">
        {loading ? (
          <>
            <SkeletonCard />
            <SkeletonCard />
            <SkeletonCard />
          </>
        ) : (
          <>
            <Metric
              label="Verified messages"
              value={filteredMessages.length.toString()}
              icon={<Inbox size={22} />}
            />
            <Metric label="Sats collected" value={totalSats.toString()} icon={<Zap size={22} />} />
            <Metric
              label="Active forms"
              value={forms.length.toString()}
              icon={<LayoutGrid size={22} />}
            />
          </>
        )}
      </section>

      {/* Workspace Section */}
      <section className="grid gap-8 lg:grid-cols-[420px_1fr]">
        {/* Left Column: Manage & Create Forms */}
        <div className="flex flex-col gap-8">
          {/* Card 1: Select/Edit Form */}
          <div className="rounded-2xl border border-slate-200/60 bg-white p-6 md:p-8 shadow-panel">
            <h2 className="brand-font text-xl font-bold text-satBlack tracking-tight">
              Active Forms
            </h2>
            <p className="mt-1 text-xs text-slate-500 leading-normal">
              Choose a form to configure embed details and view submissions.
            </p>
            <div className="mt-6 flex flex-col gap-3">
              {forms.length === 0 && !loading && (
                <p className="text-sm italic text-slate-400 py-4 text-center border border-dashed border-slate-200 rounded-xl">
                  No forms found. Create one below to begin.
                </p>
              )}
              {forms.map((form) =>
                editState?.id === form.id ? (
                  // ── inline edit row ──
                  <div
                    key={form.id}
                    className="flex flex-col gap-3 rounded-xl border border-satBlue bg-blue-50/20 p-4 shadow-sm"
                  >
                    <div className="space-y-2">
                      <input
                        className="h-10 w-full rounded-lg border border-slate-200 bg-white px-3 text-sm outline-none focus:border-satBlue focus:ring-4 focus:ring-blue-50 transition"
                        placeholder="Form name"
                        value={editState.name}
                        onChange={(e) => setEditState({ ...editState, name: e.target.value })}
                      />
                      <input
                        className="h-10 w-full rounded-lg border border-slate-200 bg-white px-3 text-sm outline-none focus:border-satBlue focus:ring-4 focus:ring-blue-50 transition"
                        placeholder="Allowed domain"
                        value={editState.domain}
                        onChange={(e) => setEditState({ ...editState, domain: e.target.value })}
                      />
                      <input
                        className="h-10 w-full rounded-lg border border-slate-200 bg-white px-3 text-sm outline-none focus:border-satBlue focus:ring-4 focus:ring-blue-50 transition"
                        type="number"
                        min={1}
                        max={10000}
                        placeholder="Sats"
                        value={editState.amount_sats}
                        onChange={(e) =>
                          setEditState({ ...editState, amount_sats: Number(e.target.value) })
                        }
                      />
                    </div>
                    <div className="flex gap-2">
                      <button
                        onClick={handleSaveEdit}
                        disabled={saving}
                        className="flex flex-1 items-center justify-center gap-1.5 rounded-lg bg-satBlue py-2 text-xs font-bold text-white transition hover:bg-blue-700 disabled:opacity-50"
                      >
                        <Check size={13} />
                        {saving ? "Saving…" : "Save"}
                      </button>
                      <button
                        onClick={cancelEdit}
                        className="flex flex-1 items-center justify-center gap-1.5 rounded-lg border border-slate-200 bg-white py-2 text-xs font-bold text-slate-600 transition hover:bg-slate-50"
                      >
                        <X size={13} />
                        Cancel
                      </button>
                    </div>
                  </div>
                ) : (
                  // ── normal form row ──
                  <div key={form.id} className="group relative">
                    <button
                      onClick={() => setSelectedFormId(form.id)}
                      className={`flex w-full flex-col rounded-xl border p-4 text-left transition duration-200 ${
                        selectedForm?.id === form.id
                          ? "border-satBlue bg-blue-50/20 shadow-sm"
                          : "border-slate-200 hover:border-slate-300 hover:bg-slate-50/50"
                      }`}
                    >
                      <span className="pr-16 text-sm font-bold text-satBlack">{form.name}</span>
                      <span className="mt-1 text-xs text-slate-500 flex items-center gap-1.5">
                        <Globe size={11} className="text-slate-400" />
                        {form.domain} · {form.amount_sats} sats
                      </span>
                    </button>
                    <div className="absolute right-3 top-1/2 flex -translate-y-1/2 gap-1.5 opacity-0 transition group-hover:opacity-100">
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          startEdit(form);
                        }}
                        className="flex h-8 w-8 items-center justify-center rounded-lg text-slate-400 bg-white border border-slate-200 shadow-sm hover:bg-blue-50 hover:text-satBlue transition hover:border-blue-200"
                        title="Edit form"
                      >
                        <Pencil size={14} />
                      </button>
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleDeleteForm(form.id);
                        }}
                        className="flex h-8 w-8 items-center justify-center rounded-lg text-slate-400 bg-white border border-slate-200 shadow-sm hover:bg-red-50 hover:text-satRed transition hover:border-red-200"
                        title="Delete form"
                      >
                        <Trash2 size={14} />
                      </button>
                    </div>
                  </div>
                ),
              )}
            </div>
          </div>

          {/* Card 2: Create a New Form */}
          <div className="relative overflow-hidden rounded-2xl border border-slate-200 bg-white p-6 md:p-8 shadow-panel">
            <FormProgressBar active={submitting} />
            <div className="flex items-start justify-between gap-3">
              <div>
                <h2 className="brand-font text-xl font-bold text-satBlack tracking-tight">
                  New Form
                </h2>
                <p className="mt-1 text-xs text-slate-500 leading-normal">
                  Deploy a spam-resistant Lightning gate widget.
                </p>
              </div>
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-blue-50 text-satBlue">
                <ShieldCheck size={20} aria-hidden="true" />
              </div>
            </div>

            <form className="mt-6 space-y-5" onSubmit={handleCreateForm}>
              <div>
                <label className="block text-xs font-bold uppercase tracking-wider text-slate-500 mb-2">
                  Form Name
                </label>
                <div className="relative rounded-xl shadow-sm">
                  <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5">
                    <FileText className="h-4 w-4 text-slate-400" />
                  </div>
                  <input
                    className="block w-full rounded-xl border border-slate-200/80 bg-white py-3 pl-10 pr-4 text-sm text-satBlack placeholder-slate-400 outline-none transition duration-200 focus:border-satBlue focus:ring-4 focus:ring-blue-50"
                    placeholder="e.g., Portfolio Contact"
                    required
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                  />
                </div>
              </div>

              <div>
                <label className="block text-xs font-bold uppercase tracking-wider text-slate-500 mb-2">
                  Allowed Domain
                </label>
                <div className="relative rounded-xl shadow-sm">
                  <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5">
                    <Globe className="h-4 w-4 text-slate-400" />
                  </div>
                  <input
                    className="block w-full rounded-xl border border-slate-200/80 bg-white py-3 pl-10 pr-4 text-sm text-satBlack placeholder-slate-400 outline-none transition duration-200 focus:border-satBlue focus:ring-4 focus:ring-blue-50"
                    placeholder="e.g., myportfolio.com"
                    required
                    value={domain}
                    onChange={(e) => setDomain(e.target.value)}
                  />
                </div>
              </div>

              <div>
                <label className="block text-xs font-bold uppercase tracking-wider text-slate-500 mb-2">
                  Price in Sats
                </label>
                <div className="relative rounded-xl shadow-sm">
                  <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5">
                    <Coins className="h-4 w-4 text-slate-400" />
                  </div>
                  <input
                    className="block w-full rounded-xl border border-slate-200/80 bg-white py-3 pl-10 pr-4 text-sm text-satBlack placeholder-slate-400 outline-none transition duration-200 focus:border-satBlue focus:ring-4 focus:ring-blue-50"
                    min={1}
                    max={10000}
                    type="number"
                    required
                    value={amount}
                    onChange={(e) => setAmount(Number(e.target.value))}
                  />
                </div>
              </div>

              <PrimaryButton
                className="w-full h-12 rounded-xl text-sm font-bold shadow-md shadow-satBlue/10 hover:shadow-satBlue/20"
                type="submit"
                disabled={submitting}
              >
                {submitting ? "Creating..." : "Create Form"}
              </PrimaryButton>
            </form>
          </div>
        </div>

        {/* Right Column: Form Workspace, Embed Code & Inbox Submissions */}
        <div className="flex flex-col gap-8">
          {selectedForm ? (
            <div className="rounded-2xl border border-slate-200/60 bg-white p-6 md:p-8 shadow-panel flex flex-col gap-8">
              {/* Workspace Header */}
              <div className="border-b border-slate-200/80 pb-5">
                <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                  <div>
                    <span className="inline-flex items-center gap-1.5 rounded-full bg-emerald-50 px-2.5 py-1 text-xs font-bold text-satGreen border border-emerald-100">
                      <span className="w-1.5 h-1.5 rounded-full bg-satGreen" />
                      Active Form
                    </span>
                    <h2 className="brand-font mt-2.5 text-2xl font-bold text-satBlack tracking-tight">
                      {selectedForm.name}
                    </h2>
                    <p className="mt-1 text-sm text-slate-500">
                      Domain:{" "}
                      <span className="font-semibold text-slate-700">{selectedForm.domain}</span> ·
                      Cost:{" "}
                      <span className="font-semibold text-slate-700">
                        {selectedForm.amount_sats} sats
                      </span>
                    </p>
                  </div>
                  <div className="text-left sm:text-right rounded-xl bg-slate-50 p-4 border border-slate-100">
                    <p className="text-[10px] font-bold uppercase tracking-wider text-slate-500">
                      Form Balance
                    </p>
                    <p className="brand-font mt-1 text-2xl font-extrabold text-satBlue">
                      {selectedForm.sat_balance} sats
                    </p>
                  </div>
                </div>
              </div>

              {/* Embed Code Snippet Card */}
              <div className="rounded-xl border border-slate-200/60 bg-slate-50/50 p-5 md:p-6 transition hover:bg-slate-50">
                <div className="flex items-center justify-between">
                  <h3 className="text-sm font-bold text-satBlack flex items-center gap-2">
                    <Code size={16} className="text-satBlue" />
                    Embed Code Snippet
                  </h3>
                </div>
                <p className="mt-1 text-xs text-slate-500 leading-normal">
                  Copy this iframe code snippet and paste it into any web page or content block.
                </p>
                <div className="relative mt-4">
                  <pre className="overflow-x-auto rounded-xl bg-slate-950 p-4 text-[11px] font-mono leading-relaxed text-slate-300">
                    {embedSnippet || "No form selected to generate snippet."}
                  </pre>
                  <PrimaryButton
                    className="absolute right-3 top-3 h-8 px-3 rounded-lg text-xs"
                    disabled={!embedSnippet}
                    onClick={copyEmbed}
                    tone="neutral"
                  >
                    {copied ? (
                      <Check size={13} className="mr-1" />
                    ) : (
                      <Copy size={13} className="mr-1" />
                    )}
                    {copied ? "Copied" : "Copy"}
                  </PrimaryButton>
                </div>
              </div>

              {/* Withdraw Collected Sats Card */}
              <div className="rounded-2xl border border-slate-200/60 bg-white p-6 md:p-8 shadow-panel flex flex-col gap-4">
                <div>
                  <h3 className="brand-font text-xl font-bold text-satBlack tracking-tight flex items-center gap-2">
                    <Coins size={20} className="text-satBlue" />
                    Withdraw Collected Sats
                  </h3>
                  <p className="mt-1 text-xs text-slate-500 leading-normal">
                    Enter a Lightning Invoice (BOLT11) to withdraw funds from this form&apos;s
                    balance.
                  </p>
                </div>
                <div className="flex flex-col sm:flex-row gap-3">
                  <input
                    className="flex-1 h-11 rounded-xl border border-slate-200 bg-white px-4 py-2 text-sm text-satBlack placeholder-slate-400 outline-none focus:border-satBlue focus:ring-4 focus:ring-blue-50 transition"
                    placeholder="lnbc..."
                    value={withdrawInvoice}
                    onChange={(e) => setWithdrawInvoice(e.target.value)}
                  />
                  <PrimaryButton
                    onClick={handleWithdraw}
                    disabled={withdrawing || !withdrawInvoice}
                    className="h-11 px-6 rounded-xl font-bold brand-font text-sm shrink-0 shadow-md shadow-satBlue/10 hover:shadow-satBlue/20"
                  >
                    {withdrawing ? "Processing..." : "Withdraw"}
                  </PrimaryButton>
                </div>
                {withdrawStatus && (
                  <p
                    className={`text-xs font-bold ${withdrawStatus.success ? "text-satGreen" : "text-satRed"}`}
                  >
                    {withdrawStatus.message}
                  </p>
                )}
              </div>

              {/* Verified Inbox */}
              <div>
                <div className="flex flex-col gap-3 border-b border-slate-200/80 pb-4 sm:flex-row sm:items-center sm:justify-between">
                  <div>
                    <h3 className="brand-font text-lg font-bold text-satBlack tracking-tight">
                      Verified Inbox
                    </h3>
                    <p className="mt-1 text-xs text-slate-500 leading-normal">
                      Only paid contact form submissions display here.
                    </p>
                  </div>
                </div>

                <div className="mt-6 space-y-4">
                  {loading ? (
                    <>
                      <SkeletonMessage />
                      <SkeletonMessage />
                    </>
                  ) : filteredMessages.length === 0 ? (
                    <div className="flex min-h-[300px] flex-col items-center justify-center rounded-xl border border-dashed border-slate-200/80 p-8 text-center bg-slate-50/30">
                      <div className="flex h-14 w-14 items-center justify-center rounded-full bg-slate-50 text-slate-400 border border-slate-100 shadow-sm">
                        <Inbox className="text-slate-400" size={24} aria-hidden="true" />
                      </div>
                      <h3 className="brand-font mt-4 text-lg font-bold text-satBlack">
                        No verified messages yet
                      </h3>
                      <p className="mt-2 max-w-sm text-sm text-slate-500 leading-relaxed">
                        Open the demo site, submit the widget form, and use the mock pay button to
                        settle the invoice.
                      </p>
                    </div>
                  ) : (
                    filteredMessages.map((message) => (
                      <article
                        className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm transition hover:border-blue-200 hover:bg-blue-50/10 hover:shadow-md duration-200"
                        key={message.id}
                      >
                        <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                          <div>
                            <h3 className="text-sm font-bold text-slate-800">
                              Verified Submission
                            </h3>
                            <p className="text-xs text-slate-400 mt-1">
                              {new Date(message.timestamp).toLocaleString()}
                            </p>
                          </div>
                          <div className="flex items-center gap-2">
                            <StatusBadge status={message.status} />
                            <span className="rounded-full bg-emerald-50 border border-emerald-100 px-3 py-1 text-xs font-bold text-satGreen">
                              +{message.paid_sats} sats
                            </span>
                          </div>
                        </div>
                        <p className="mt-4 whitespace-pre-wrap text-sm leading-relaxed text-slate-600">
                          {message.message}
                        </p>
                      </article>
                    ))
                  )}
                </div>
              </div>
            </div>
          ) : (
            <div className="rounded-2xl border border-slate-200/60 bg-white p-12 shadow-panel text-center flex flex-col items-center justify-center min-h-[450px]">
              <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-blue-50 text-satBlue mb-4">
                <PlusCircle size={32} />
              </div>
              <h2 className="brand-font text-xl font-bold text-satBlack tracking-tight">
                Welcome to SatGo
              </h2>
              <p className="mt-2 max-w-md text-sm text-slate-500 leading-relaxed">
                Create a protected form configuration using the generator on the left, and embed it
                into your site to start filtering contact form spam.
              </p>
            </div>
          )}
        </div>
      </section>
    </main>
  );
}
