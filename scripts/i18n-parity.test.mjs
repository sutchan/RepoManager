// scripts/i18n-parity.test.mjs v2.3.16
// 三语键一致性校验（S14）：覆盖弃用参考界面 src/lib/i18n.ts 与对外展示页 landing/assets/dict.js，零依赖 node:test。
import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import vm from 'node:vm';

const SRC = 'src/lib/i18n.ts';
const LANDING_DICT = 'landing/assets/dict.js';
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

// landing 是实际对外页面：其字典失衡此前无门禁覆盖，新增沙箱求值后比对三语键集合。
function landingDict() {
	const sandbox = { window: {} };
	vm.runInNewContext(fs.readFileSync(LANDING_DICT, 'utf8'), sandbox);
	return sandbox.window.LANDING_DICT;
}

test('landing 三语字典键集合一致', () => {
	const dict = landingDict();
	assert.ok(dict && dict.zh && dict['zh-Hant'] && dict.en, 'landing 字典缺少语言块');
	const keys = Object.keys(dict.zh).sort();
	assert.ok(keys.length > 0, 'landing 简中字典不应为空');
	assert.deepEqual(Object.keys(dict['zh-Hant']).sort(), keys);
	assert.deepEqual(Object.keys(dict.en).sort(), keys);
});

test('landing 页面标题与描述三语齐全', () => {
	const dict = landingDict();
	for (const lang of ['zh', 'zh-Hant', 'en']) {
		assert.ok(dict[lang].pageTitle, `${lang} 缺少 pageTitle`);
		assert.ok(dict[lang].pageDesc, `${lang} 缺少 pageDesc`);
	}
});
