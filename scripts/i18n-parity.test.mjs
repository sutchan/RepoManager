// scripts/i18n-parity.test.mjs v2.3.11
// 弃用参考界面 src/lib/i18n.ts 的三语键一致性校验（S14，零依赖 node:test）。
import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';

const SRC = 'src/lib/i18n.ts';
const text = fs.readFileSync(SRC, 'utf8');

function keysOf(lang) {
	const token = lang === 'zh-Hant' ? '"zh-Hant"' : lang;
	const block = text.match(new RegExp(`${token}:\\s*\\{([\\s\\S]*?)\\n  \\},`));
	assert.ok(block, `缺少语言块 ${lang}`);
	return [...block[1].matchAll(/^\s{4}([A-Za-z]\w*):/gm)].map((m) => m[1]).sort();
}

test('三语键集合一致', () => {
	const zh = keysOf('zh');
	assert.ok(zh.length > 0, '简中字典不应为空');
	assert.deepEqual(keysOf('zh-Hant'), zh);
	assert.deepEqual(keysOf('en'), zh);
});

test('SUPPORTED_LANGS 顺序为简中 → 繁中 → 英文', () => {
	const block = text.match(/SUPPORTED_LANGS:\s*Lang\[\]\s*=\s*\[([^\]]+)\]/);
	assert.ok(block, '缺少 SUPPORTED_LANGS 声明');
	const langs = [...block[1].matchAll(/"([^"]+)"/g)].map((m) => m[1]);
	assert.deepEqual(langs, ['zh', 'zh-Hant', 'en']);
});

test('HTML_LANG 映射正确', () => {
	const block = text.match(/HTML_LANG:[\s\S]*?=\s*\{([\s\S]*?)\};/);
	assert.ok(block, '缺少 HTML_LANG 声明');
	assert.match(block[1], /zh:\s*"zh"/);
	assert.match(block[1], /"zh-Hant":\s*"zh-Hant"/);
	assert.match(block[1], /en:\s*"en"/);
});
