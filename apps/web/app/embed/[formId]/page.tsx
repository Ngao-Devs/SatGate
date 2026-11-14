import { WidgetForm } from "@/components/WidgetForm";

export default async function EmbedPage({ params }: { params: Promise<{ formId: string }> }) {
  const { formId } = await params;

  return <WidgetForm formId={formId} />;
}
