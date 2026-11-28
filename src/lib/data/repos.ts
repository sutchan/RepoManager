// src/lib/data/repos.ts v2.3.16
// 扫描结果与优化建议（参考界面 dashboard 演示数据）
// 注意：以下均为 mock 数据，用于展示界面结构，非真实磁盘扫描结果；容量数字为示意，不可作为事实引用。
// 优化建议文案随界面语言切换，保持简中 → 繁中 → 英文三语齐全（对齐 OPENSPEC §4.1）。

import type { Lang } from "@/lib/i18n";

export type Risk = "safe" | "warning";

export interface ScanResult {
  repo: string;
  path: string;
  risk: Risk;
}

export const scanResults: ScanResult[] = [
  { repo: "next-dashboard", path: "E:\\Github\\next-dashboard", risk: "safe" },
  { repo: "vue-admin", path: "E:\\Github\\vue-admin", risk: "warning" },
  { repo: "react-portfolio", path: "E:\\Github\\react-portfolio", risk: "safe" },
  { repo: "node-api", path: "E:\\Github\\node-api", risk: "safe" },
];

export interface Tweak {
  id: string;
  category: Record<Lang, string>;
  title: Record<Lang, string>;
  desc: Record<Lang, string>;
}

export const tweaks: Tweak[] = [
  {
    id: "global-node",
    category: { zh: "依赖", "zh-Hant": "依賴", en: "Dependencies" },
    title: { zh: "启用全局 node_modules", "zh-Hant": "啟用全域 node_modules", en: "Enable global node_modules" },
    desc: {
      zh: "为全部仓库写入 .npmrc 并配置 NODE_PATH，预计回收约 1.2 GB。",
      "zh-Hant": "為所有專案寫入 .npmrc 並設定 NODE_PATH，預計可回收約 1.2 GB。",
      en: "Write .npmrc and configure NODE_PATH for all repos; about 1.2 GB can be reclaimed.",
    },
  },
  {
    id: "next-cache",
    category: { zh: "缓存", "zh-Hant": "快取", en: "Cache" },
    title: { zh: "清理 .next 构建缓存", "zh-Hant": "清理 .next 建置快取", en: "Clear .next build cache" },
    desc: {
      zh: "2 个 Next.js 项目共占用 860 MB，禁用后热更新略慢但省空间。",
      "zh-Hant": "2 個 Next.js 專案共佔用 860 MB，停用後熱更新稍慢但節省空間。",
      en: "2 Next.js projects use 860 MB; disabling slows hot reload slightly but saves space.",
    },
  },
  {
    id: "vite-cache",
    category: { zh: "缓存", "zh-Hant": "快取", en: "Cache" },
    title: { zh: "清理 Vite 构建缓存", "zh-Hant": "清理 Vite 建置快取", en: "Clear Vite build cache" },
    desc: {
      zh: "3 个 Vite 项目 .vite 目录共占用 540 MB。",
      "zh-Hant": "3 個 Vite 專案的 .vite 目錄共佔用 540 MB。",
      en: "The .vite directories of 3 Vite projects use 540 MB.",
    },
  },
  {
    id: "lockfiles",
    category: { zh: "依赖", "zh-Hant": "依賴", en: "Dependencies" },
    title: { zh: "移除 lock 文件", "zh-Hant": "移除 lock 檔案", en: "Remove lock files" },
    desc: {
      zh: "6 个 package-lock.json / yarn.lock，共 4.2 MB。",
      "zh-Hant": "6 個 package-lock.json / yarn.lock，共 4.2 MB。",
      en: "6 package-lock.json / yarn.lock files, 4.2 MB in total.",
    },
  },
];
