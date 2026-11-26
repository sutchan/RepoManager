// scripts/check-asset-refs.mjs v2.3.11
// S18 冗余预防：扫描 landing/assets/ 静态资源是否被 landing/ 内源码引用。
// 告警模式：输出未引用清单但退出码恒为 0；基线建立后再转强制门禁。
import { readdirSync, readFileSync, statSync } from "node:fs";
import { join } from "node:path";

const ROOT = process.cwd();
const ASSET_DIR = join(ROOT, "landing", "assets");
const SCAN_DIR = join(ROOT, "landing");

function collectSources(dir) {
  const files = [];
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    if (entry.name === "node_modules" || entry.name === ".next") continue;
    const full = join(dir, entry.name);
    if (entry.isDirectory()) {
      files.push(...collectSources(full));
    } else if (/\.(html|js|mjs|cjs|css|tsx?)$/.test(entry.name)) {
      files.push(full);
    }
  }
  return files;
}

const sources = collectSources(SCAN_DIR)
  .map((file) => readFileSync(file, "utf8"))
  .join("\n");
const assets = readdirSync(ASSET_DIR).filter((name) => statSync(join(ASSET_DIR, name)).isFile());
const unreferenced = assets.filter((name) => !sources.includes(name));

if (unreferenced.length > 0) {
  console.log(`[asset-refs] 未引用资源（${unreferenced.length}）：${unreferenced.join(", ")}`);
} else {
  console.log(`[asset-refs] 全部 ${assets.length} 个资源均被引用`);
}
console.log("[asset-refs] 当前为告警模式（退出码 0）；基线稳定后可转强制。");
