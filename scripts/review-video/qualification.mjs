// Read-only replay of executed qualification evidence, not a product UI walkthrough.
import fs from 'node:fs';
import path from 'node:path';
import cp from 'node:child_process';
import crypto from 'node:crypto';
import assert from 'node:assert/strict';
import {createRequire} from 'node:module';
const require=createRequire(import.meta.url);
const root=path.resolve(import.meta.dirname,'../..');
const liveDir=path.resolve(process.env.QUALIFICATION_LIVE_DIR||path.join(root,'artifacts/interaction-runtime/20261028/host-isolation-retry2'));
const fixtureDir=path.join(root,'artifacts/interaction-review/20261028');
const out=path.resolve(process.env.REVIEW_OUTPUT_DIR||path.join(root,'artifacts/interaction-qualification/20261028'));
const hash=p=>crypto.createHash('sha256').update(fs.readFileSync(p)).digest('hex');
const read=p=>JSON.parse(fs.readFileSync(p,'utf8'));
const live=read(path.join(liveDir,'result.json')),manifest=read(path.join(liveDir,'source-inventory.json'));
assert.ok(live.passed&&live.automated_host_api_submissions&&!live.browser_ui_evidence);
for(const key of ['source_files_unchanged','fixture_files_unchanged','binary_unchanged'])assert.equal(live[key],true,key);
const checks=['distinct_project_native_threads','cross_scope_spoof_rejected','answer_one_leaves_other_pending','alpha_agent_received_answer','new_http_connection_retains_snapshot','decline_received_by_agent','actual_multiple_choice_ids_received','restart_forgets_old_conversation'];
for(const key of checks)assert.equal(live.checks[key],true,key);
for(const [file,sha]of Object.entries(manifest.files))assert.equal(hash(path.join(root,file)),sha,'Source changed: '+file);
const fixture=read(path.join(fixtureDir,'browser-fixture.json'));assert.equal(fixture.status,'pass');
const fixtureChecks=['native user-input exact response','foreign-window/null-source/wrong-origin messages ignored','answer draft and focus survive polling','host credential inaccessible to frame','renewed launch fragment erases secret and resets rich-view trust','view mount requires explicit host authorization','interrupting blocks follow-up and repeat stop','native approvals expose supported one-turn strings only'];
for(const check of fixtureChecks)assert.ok(fixture.checks.includes(check),check);
const events=fs.readFileSync(path.join(fixtureDir,'go-test.jsonl'),'utf8').trim().split('\n').map(line=>JSON.parse(line));
assert.ok(!events.some(e=>e.Action==='fail'),'Named test run must contain no failures');
const tests=[['runtime','TestConcurrentRequestsAndNativeReplyCorrelation'],['runtime','TestPreferenceSelectionAndAdapterIsolation'],['runtime','TestModelAndEffortRequireAdvertisedSupport'],['runtime','TestMalformedPreferencesFailWithoutExecution'],['web','TestNativeAnswersRequireExplicitSupportedDisposition'],['web','TestApprovalMustMatchOfferedDecision'],['web','TestLateCompletionCannotReleaseNewTurn'],['web','TestDelayedStartResponseCannotBindToNextTurn'],['interaction','TestWaitWakesOnAnswerAndShutdown']];
for(const [pkg,test]of tests)assert.ok(events.some(e=>e.Action==='pass'&&e.Test===test&&e.Package.endsWith('/'+pkg)),test);
const entries=[
 {title:'Live Codex: isolation, multiple choice and restart',ids:['U034','U076','U078','U079','U080'],kind:'ACTUAL RUNTIME · AUTOMATED HOST API SUBMISSIONS · NOT BROWSER UI',explanation:'Real native threads remain isolated; submitted options reach the agent, and restart rejects old conversation IDs.',lines:['PASS  Two projects use distinct native threads','PASS  Cross-scope answer spoof rejected; other question stays pending','PASS  Agent receives ALPHA_ONLY_real_input; decline reaches the other agent','PASS  New HTTP connection retains the pending snapshot','PASS  Agent acknowledges selected IDs: '+live.selected_ids.join(', '),'PASS  Restart forgets conversations and rejects the old ID'],source:path.join(liveDir,'result.json'),assertions:checks},
 {title:'Deterministic browser fixture: host boundaries',ids:['U036','U070','U076','U077','U081'],kind:'DETERMINISTIC MOCK RUNTIME FIXTURE · EXECUTED BROWSER CHECKS',explanation:'The fixture exercises native input, approval controls, frame spoof rejection and trusted host view mounting.',lines:['PASS  Native user-input response preserves exact submitted text','PASS  Foreign-window, null-source and wrong-origin messages ignored','PASS  Answer draft and focus survive polling; host credential stays private','PASS  Renewed launch clears secret and resets rich-view trust','PASS  View mount needs host authorization; interruption blocks new turns','PASS  Native approval controls offer supported one-turn values only'],source:path.join(fixtureDir,'browser-fixture.json'),assertions:fixtureChecks},
 {title:'Named regression tests: preferences and lifecycle',ids:['U034','U036','U070','U076','U079','U080','U081'],kind:'DETERMINISTIC GO TESTS · MOCK PROTOCOL / IN-PROCESS CHECKS',explanation:'Current named tests cover preference selection, reply correlation, approval limits and late-event lifecycle handling.',lines:tests.map(([pkg,test])=>'PASS  '+pkg+'/'+test),source:path.join(fixtureDir,'go-test.jsonl'),assertions:tests.map(([pkg,test])=>pkg+'/'+test)}
];
assert.ok(!fs.existsSync(out),'Preserve earlier attempts; select a fresh output directory');fs.mkdirSync(out,{recursive:true});
const sources=[path.join(liveDir,'result.json'),path.join(liveDir,'source-inventory.json'),path.join(fixtureDir,'browser-fixture.json'),path.join(fixtureDir,'go-test.jsonl')].map(file=>({path:path.relative(root,file),sha256:hash(file)}));
const {chromium}=require(process.env.PLAYWRIGHT_MODULE||'playwright');
const browser=await chromium.launch({headless:true,...(process.env.REVIEW_CHROMIUM?{executablePath:process.env.REVIEW_CHROMIUM}:{})});
const context=await browser.newContext({viewport:{width:1440,height:1000},recordVideo:{dir:path.join(out,'raw'),size:{width:1440,height:1000}}});
const page=await context.newPage(),video=page.video();
await page.setContent(`<style>*{box-sizing:border-box}body{margin:0;background:#101d28;color:#e8f4f2;font:23px/1.6 system-ui}header{padding:16px 40px;background:#23483f;height:135px}header b{color:#8de9ce}main{padding:30px 45px}pre{font:21px/1.75 monospace;white-space:pre-wrap;background:#192c3a;border:2px solid #6dcbb4;padding:28px;border-radius:12px}small{color:#a8c4d0;font-size:18px}h1{font-size:32px}footer{position:fixed;bottom:22px;left:45px;font-size:17px;color:#a8c4d0}</style><header><b id="ids"></b><div id="explanation"></div></header><main><small id="kind"></small><h1 id="title"></h1><pre id="output"></pre><small id="source"></small></main><footer>Replay of recorded qualification evidence · ${manifest.head.slice(0,12)} · Full requirements matrix incomplete</footer>`);
const began=Date.now(),scenes=[];
for(const [i,e]of entries.entries()){
 const start=(Date.now()-began)/1000;
 await page.evaluate(e=>{document.querySelector('#ids').textContent=e.ids.join(' · ');for(const k of ['title','explanation','kind'])document.getElementById(k).textContent=e[k];document.getElementById('output').textContent=e.lines.join('\n');document.getElementById('source').textContent='Source: '+e.source;},{...e,source:path.relative(root,e.source)});
 assert.ok(await page.locator('pre').evaluate(el=>el.scrollWidth<=el.clientWidth&&el.getBoundingClientRect().bottom<900),'Evidence text must fit without clipping');
 await page.waitForTimeout(20000);await page.screenshot({path:path.join(out,`scene-${i+1}.png`)});
 scenes.push({id:i+1,title:e.title,requirement_ids:e.ids,explanation:e.explanation,evidence_kind:e.kind,start_seconds:start,end_seconds:(Date.now()-began)/1000,assertions:e.assertions,commands:[{command:'Replay verified evidence '+path.relative(root,e.source),exit:0,expected:0}]});
}
await context.close();await browser.close();
const target=path.join(out,'demo.mp4');
cp.execFileSync('ffmpeg',['-v','error','-n','-i',await video.path(),'-c:v','libx264','-preset','fast','-crf','24','-pix_fmt','yuv420p','-movflags','+faststart','-an',target]);
cp.execFileSync('ffmpeg',['-v','error','-i',target,'-f','null','-']);
const media=JSON.parse(cp.execFileSync('ffprobe',['-v','error','-show_entries','format=duration,size:stream=codec_name,width,height,nb_frames','-of','json',target],{encoding:'utf8'}));
fs.writeFileSync(path.join(out,'results.json'),JSON.stringify({revision:manifest.head,scenes,sources,binary_sha256:live.binary_sha256,source_manifest:manifest,source_files_verified:true,stable_alias_updated:false,media,video_sha256:hash(target)},null,2));
console.log(JSON.stringify({out,scenes:scenes.length,sha256:hash(target),media}));
