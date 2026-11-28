// landing/assets/app.js v2.3.13
// 语言切换：偏好持久化并做枚举校验；文案来自 dict.js，版本来自构建期生成的 version.js。
const SUPPORTED = ["zh", "zh-Hant", "en"];
const HTML_LANG = { zh: "zh", "zh-Hant": "zh-Hant", en: "en" };
const STORAGE_KEY = "repo-manager-landing-lang";
const DICT = window.LANDING_DICT;

function readLang() {
  try {
    const saved = localStorage.getItem(STORAGE_KEY);
    return SUPPORTED.includes(saved) ? saved : "zh";
  } catch {
    return "zh";
  }
}

function applyLang(lang) {
  const chosen = SUPPORTED.includes(lang) ? lang : "zh";
  const dict = DICT[chosen] || DICT.zh;
  document.documentElement.lang = HTML_LANG[chosen];
  document.querySelectorAll("[data-i18n]").forEach((el) => {
    const value = dict[el.dataset.i18n];
    if (value) el.textContent = value;
  });
  document.querySelectorAll("#landing-lang-switch button").forEach((btn) => {
    const active = btn.dataset.lang === chosen;
    btn.setAttribute("aria-pressed", String(active));
    btn.classList.toggle("is-active", active);
  });
  const label = window.LANDING_VERSION && window.LANDING_VERSION.label;
  document.getElementById("landing-version").textContent = label || dict.versionMissing;
  try {
    localStorage.setItem(STORAGE_KEY, chosen);
  } catch {
    // 隐私模式下无法写入偏好时静默降级，保持页面可用
  }
}

document.addEventListener("DOMContentLoaded", () => {
  document.querySelectorAll("#landing-lang-switch button").forEach((btn) => {
    btn.addEventListener("click", () => applyLang(btn.dataset.lang));
  });
  applyLang(readLang());
});
