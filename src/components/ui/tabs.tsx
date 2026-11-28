// src/components/ui/tabs.tsx v2.3.14
"use client";
import * as React from "react";
import { cn } from "@/lib/utils";

interface TabsCtx { value: string; setValue: (v: string) => void; }
const Ctx = React.createContext<TabsCtx | null>(null);

/**
 * 标签页容器：同时支持受控（value + onValueChange）与非受控（defaultValue）用法。
 */
export function Tabs({
  defaultValue = "", value, onValueChange, id, className, children,
}: {
  defaultValue?: string;
  value?: string;
  onValueChange?: (v: string) => void;
  id?: string;
  className?: string;
  children: React.ReactNode;
}) {
  const [internal, setInternal] = React.useState(defaultValue);
  const isControlled = value !== undefined;
  const current = isControlled ? value : internal;
  const setValue = (v: string) => {
    if (!isControlled) setInternal(v);
    onValueChange?.(v);
  };
  return (
    <div id={id} className={cn(className)}>
      <Ctx.Provider value={{ value: current, setValue }}>{children}</Ctx.Provider>
    </div>
  );
}

export function TabsList({ id, className, children }: { id?: string; className?: string; children: React.ReactNode }) {
  const ctx = React.useContext(Ctx)!;
  const values = React.Children.toArray(children)
    .filter(React.isValidElement)
    .map((child) => (child.props as { value?: string }).value)
    .filter((v): v is string => typeof v === "string");
  const onKeyDown = (event: React.KeyboardEvent) => {
    if (event.key !== "ArrowLeft" && event.key !== "ArrowRight") return;
    const index = values.indexOf(ctx.value);
    if (index < 0) return;
    event.preventDefault();
    const step = event.key === "ArrowRight" ? 1 : values.length - 1;
    const next = values[(index + step) % values.length];
    ctx.setValue(next);
    document.getElementById(`tab-${next}`)?.focus();
  };
  return (
    <div id={id} role="tablist" onKeyDown={onKeyDown} className={cn("inline-flex rounded-[8px] bg-muted p-1", className)}>
      {children}
    </div>
  );
}

export function TabsTrigger({ id, value, children }: { id?: string; value: string; children: React.ReactNode }) {
  const ctx = React.useContext(Ctx)!;
  const active = ctx.value === value;
  return (
    <button
      id={id ?? `tab-${value}`}
      role="tab"
      aria-selected={active}
      aria-controls={`tab-content-${value}`}
      tabIndex={active ? 0 : -1}
      onClick={() => ctx.setValue(value)}
      className={cn(
        "rounded-[6px] px-3 py-1.5 text-sm font-medium transition-colors",
        active ? "bg-card text-foreground shadow-sm" : "text-muted-foreground hover:text-foreground"
      )}
    >
      {children}
    </button>
  );
}

export function TabsContent({ id, value, children, className }: { id?: string; value: string; children: React.ReactNode; className?: string }) {
  const ctx = React.useContext(Ctx)!;
  if (ctx.value !== value) return null;
  return (
    <div id={id ?? `tab-content-${value}`} role="tabpanel" aria-labelledby={`tab-${value}`} className={cn("mt-4", className)}>
      {children}
    </div>
  );
}
