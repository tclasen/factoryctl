// Record explorer and actual runtime/browser interactions; all output stays in ignored artifacts/.
import fs from 'node:fs';
import path from 'node:path';
import crypto from 'node:crypto';
import cp from 'node:child_process';
import {createRequire} from 'node:module';
import assert from 'node:assert/strict';
import os from 'node:os';
const require=createRequire(import.meta.url);
const {chromium}=require(process.env.PLAYWRIGHT_MODULE||'playwright');
const root=path.resolve(import.meta.dirname,'../..');
const dir=path.resolve(process.env.REVIEW_OUTPUT_DIR||path.join(root,'artifacts/interaction-web/20261028'));
assert.ok(['demo.mp4','results.json','failure.json'].every(name=>!fs.existsSync(path.join(dir,name))),'Use a fresh archive directory; never overwrite successful or failed evidence');
fs.mkdirSync(dir,{recursive:true});
assert.ok(!process.env.REVIEW_BASE_URL,'External server capture is not supported: this recorder builds and owns the exact binary it records');
const revision=cp.execFileSync('git',['rev-parse','HEAD'],{cwd:root,encoding:'utf8'}).trim();
const gitStatus=cp.execFileSync('git',['status','--porcelain'],{cwd:root,encoding:'utf8'}).trim();
const productStatus=cp.execFileSync('git',['status','--porcelain','--','cmd','internal','go.mod','go.sum'],{cwd:root,encoding:'utf8'}).trim();
const hash=file=>crypto.createHash('sha256').update(fs.readFileSync(file)).digest('hex');
function inventory(directory,predicate=()=>true){const files=[];function walk(p){for(const e of fs.readdirSync(p,{withFileTypes:true})){const f=path.join(p,e.name);if(e.isDirectory())walk(f);else if(predicate(f))files.push({path:path.relative(root,f),sha256:hash(f)});}}walk(directory);return files.sort((a,b)=>a.path.localeCompare(b.path));}
const productManifest=()=>[...inventory(path.join(root,'cmd'),f=>f.endsWith('.go')),...inventory(path.join(root,'internal'),f=>f.endsWith('.go')||f.includes(path.sep+'assets'+path.sep)),...['go.mod','go.sum'].map(p=>({path:p,sha256:hash(path.join(root,p))}))].sort((a,b)=>a.path.localeCompare(b.path));
const before=inventory(path.join(root,'factory'),f=>f.endsWith('.md'));
const productBefore=productManifest();
const temporary=fs.mkdtempSync(path.join(os.tmpdir(),'factoryctl-web-capture-'));
const binary=path.join(temporary,process.platform==='win32'?'factoryctl.exe':'factoryctl');
let base='',browser,context,page,video,server;
const serverOutput={stdout:'',stderr:''};
const provenance={revision,repository_dirty:gitStatus!=='',product_dirty:productStatus!=='',product_status:productStatus,product_files_before_build:productBefore};
const errors=[],external=[];
async function launchServer(){
 server=cp.spawn(binary,['--path',root,'web','--runtime','codex','--listen','127.0.0.1:0'],{cwd:root,stdio:['ignore','pipe','pipe']});
 server.stdout.on('data',chunk=>{serverOutput.stdout+=chunk;});server.stderr.on('data',chunk=>{serverOutput.stderr+=chunk;});
 return await new Promise((resolve,reject)=>{
  const timer=setTimeout(()=>finish(Error('Owned browser server did not become ready within 15 seconds')),15000);
  function ready(){const match=serverOutput.stdout.match(/http:\/\/127\.0\.0\.1:\d+/);if(match && /#host=([a-f0-9]+)/.test(serverOutput.stdout))finish(null,match[0]);}
  function failed(error){finish(error);}
  function ended(code,signal){finish(Error('Owned browser server exited before readiness: '+code+' '+signal+' '+serverOutput.stderr));}
  function finish(error,url){clearTimeout(timer);server.stdout.off('data',ready);server.off('error',failed);server.off('exit',ended);error?reject(error):resolve(url);}
  server.stdout.on('data',ready);server.once('error',failed);server.once('exit',ended);ready();
 });
}
async function stopServer(){
 if(!server?.pid||server.exitCode!==null||server.signalCode!==null)return;
 await new Promise(resolve=>{const timer=setTimeout(()=>server.kill('SIGKILL'),3000);server.once('exit',()=>{clearTimeout(timer);resolve();});server.once('error',()=>{clearTimeout(timer);resolve();});server.kill('SIGTERM');});
}
const scenes=[];let began;
async function check(label,fn){await fn();scenes.at(-1).commands.push({command:label,exit:0,expected:0,evidence_kind:'executed-browser'});}
async function scene(title,ids,explanation,fn){
 const start=(Date.now()-began)/1000;const item={id:scenes.length+1,title,start_seconds:start,requirement_ids:ids,explanation,commands:[]};scenes.push(item);
 await page.evaluate(({title,ids,explanation})=>{const box=document.getElementById('review-annotation');box.replaceChildren();const h=document.createElement('strong');h.textContent=ids.join(' · ')+' — '+title;const p=document.createElement('div');p.textContent=explanation;box.append(h,p);},{title,ids,explanation});
 console.log('SCENE '+item.id+' '+title);await page.waitForTimeout(1800);await fn();await page.waitForTimeout(6500);item.end_seconds=(Date.now()-began)/1000;console.log('PASS '+item.id);
}
async function focus(selector){const node=page.locator(selector).first();await node.scrollIntoViewIfNeeded();await node.evaluate(e=>{e.style.outline='3px solid #d99d27';e.style.outlineOffset='4px';});await page.waitForTimeout(1800);await node.evaluate(e=>{e.style.outline='';});}
async function selectDocument(query,title){await page.locator('#search').fill(query);await page.locator('#document-list .document-link').filter({hasText:title}).first().click();await page.waitForFunction(t=>document.querySelector('#document-title').textContent===t,title);await page.waitForTimeout(1800);}
try {
 provenance.build_command=['go','build','-o','<task-temporary-binary>','./cmd/factoryctl'];
 provenance.build_environment={GOENV:'off',GOFLAGS:'',GOWORK:'off'};
 provenance.go_version=cp.execFileSync('go',['version'],{cwd:root,encoding:'utf8'}).trim();
 cp.execFileSync('go',['build','-o',binary,'./cmd/factoryctl'],{cwd:root,encoding:'utf8',env:{...process.env,GOENV:'off',GOFLAGS:'',GOWORK:'off',GOCACHE:process.env.GOCACHE||path.join(temporary,'go-cache')},timeout:600000});
 provenance.product_files_after_build=productManifest();
 assert.deepEqual(provenance.product_files_after_build,productBefore,'Product source changed during build');
 provenance.binary_sha256=hash(binary);
 provenance.binary_build_info=cp.execFileSync('go',['version','-m',binary],{cwd:root,encoding:'utf8'}).trim();
 base=await launchServer();provenance.server_url=base;provenance.server_arguments=['--path',root,'web','--runtime','codex','--listen','127.0.0.1:0'];
 browser=await chromium.launch({headless:true,...(process.env.REVIEW_CHROMIUM?{executablePath:process.env.REVIEW_CHROMIUM}:{})});
 context=await browser.newContext({viewport:{width:1440,height:1000},recordVideo:{dir:path.join(dir,'raw'),size:{width:1440,height:1000}}});
 page=await context.newPage();video=page.video();
 page.on('pageerror',e=>errors.push(String(e)));page.on('request',r=>{if(new URL(r.url()).origin!==base)external.push(r.url());});
 await page.goto(base+'/#host='+serverOutput.stdout.match(/#host=([a-f0-9]+)/)[1]);await page.waitForFunction(()=>document.querySelectorAll('.document-link').length>0);
 await page.evaluate(()=>{const box=document.createElement('div');box.id='review-annotation';Object.assign(box.style,{position:'fixed',top:'0',left:'0',right:'0',height:'86px',boxSizing:'border-box',padding:'14px 24px',background:'#142d28',color:'#fff',zIndex:'9999',font:'17px/1.6 system-ui'});document.body.prepend(box);document.querySelector('.topbar').style.marginTop='86px';document.querySelector('.workspace').style.height='calc(100vh - 150px)';});
 began=Date.now();
 await scene('Explore real planning documents',['U007','U026','U100'],'Actual product UI: search and open local Markdown. This scene exercises navigation; editing and boards are demonstrated in the companion product scenes.',async()=>{
  await check('Loaded real document collection',async()=>assert.ok(await page.locator('.document-link').count()>50));
  await focus('#open-plan');await page.locator('#open-plan').click();await page.waitForFunction(()=>document.querySelector('#document-title').textContent==='Current implementation plan');
  await check('Current plan opens and missing intent is not inferred',async()=>assert.equal(await page.locator('#document-state').textContent(),'Unresolved'));await focus('.intent-panel');
 });
 await scene('Inspect accepted scope',['U096','U100'],'Filter the recorded accepted decision and inspect exactly what was authorized. A label is a claim, not automatic approval.',async()=>{
  await page.locator('#state-filter').selectOption('accepted');await selectDocument('intent','Preserve intent and review tangible outcomes');
  await check('Explicit accepted method with bounded authority',async()=>{assert.equal(await page.locator('#document-state').textContent(),'Accepted');assert.match(await page.locator('#intent-authorization').textContent(),/does not authorize publishing/);});await focus('#intent-authorization');
 });
 await scene('Follow the source and read original text',['U029','U096','U100'],'Open the linked source, then raw Markdown. The request stays separate from the assistant’s interpretation.',async()=>{
  await page.locator('#source-links a').first().click();await page.waitForFunction(()=>document.querySelector('#document-title').textContent==='Operator requirements and design authority');
  await check('Source navigation reaches operator record',async()=>assert.match(await page.locator('#document-body').textContent(),/please implement all of that/));
  await page.locator('#raw-tab').click();await check('Raw Markdown preserves original content',async()=>assert.match(await page.locator('#raw-body').textContent(),/LLM-authored extraction/));await focus('#raw-tab');
 });
 await scene('Keep future designs proposed',['U031','U096','U100'],'The MCP Apps protocol remains a proposal to qualify. First-party views are implemented; external MCP Apps interoperability remains unqualified.',async()=>{
  await page.locator('#clear-filters').click();await page.locator('#state-filter').selectOption('proposed');
  await selectDocument('13-user','Qualify user questions and generative web views');
  await check('Future protocol remains proposed',async()=>{assert.equal(await page.locator('#document-state').textContent(),'Proposed');assert.match(await page.locator('#intent-authorization').textContent(),/remain proposals/);});await focus('#intent-authorization');
 });
 await scene('See what changed in the plan',['U097','U100'],'The historical CLI-first sequence is explicitly superseded and attributed to the assistant.',async()=>{
  await page.locator('#clear-filters').click();await page.locator('#state-filter').selectOption('superseded');await selectDocument('','Historical CLI-first sequence');
  await check('Superseded sequence retains assistant attribution',async()=>{assert.equal(await page.locator('#document-state').textContent(),'Superseded');assert.equal(await page.locator('#intent-owner').textContent(),'assistant');});await focus('#intent-authorization');
 });
 await scene('Distinguish a request from an approach',['U096','U100'],'Inspect the user’s requested experience. A desired outcome does not select an implementation order or hosting destination.',async()=>{
  await page.locator('#clear-filters').click();await page.locator('#state-filter').selectOption('requested');await selectDocument('','Explore the evolving planning documents');
  await check('Requested experience stays distinct from solution',async()=>{assert.equal(await page.locator('#document-state').textContent(),'Requested');assert.match(await page.locator('#intent-authorization').textContent(),/does not select/);});await focus('#document-body');
 });
 await scene('Return to one current plan',['U097','U099','U100'],'Follow the current-plan link. Reading and watching submit no acceptance; external MCP Apps interoperability remains unqualified.',async()=>{
  await page.locator('#document-body a').filter({hasText:'current plan'}).click();await page.waitForFunction(()=>document.querySelector('#document-title').textContent==='Current implementation plan');
  await check('Current plan links are navigable',async()=>assert.match(await page.locator('#document-body').textContent(),/Current work and proposed continuation/));
  await page.locator('#document-body h2').filter({hasText:'Current work and proposed continuation'}).scrollIntoViewIfNeeded();await page.waitForTimeout(3500);
  await check('Document view never infers an acceptance control',async()=>assert.equal(await page.getByRole('button',{name:/^(approve|accept)$/i}).count(),0));
 });
 const safety='This is a synthetic local review. Do not edit files, read files, execute commands, or use any tools except factory_user_ask, factory_user_wait and factory_ui_present. ';
 const waitIdle=()=>page.waitForFunction(()=>document.querySelector('#conversation-status').textContent==='idle',null,{timeout:180000});
 async function snapshot(){return await page.evaluate(async()=>{const headers={'X-Factory-Host':sessionStorage.getItem('factoryctl-host')};const list=await(await fetch('/api/conversations',{headers})).json();return await(await fetch('/api/conversations/'+list.conversations[0].id,{headers})).json();});}
 async function send(text){await page.locator('#conversation-prompt').fill(text);await page.locator('#send-conversation').click();}
 const questions=page.locator('.question-card');
 let threadID;
 await scene('Ask the real agent; answer in the browser',['U031','U032','U033','U035','U066','U076','U077','U080'],'Live Codex tool call: a free-text question stays pending until the trusted Submit answer action. Synthetic answer: Autumn.',async()=>{
  await page.locator('#action-buttons button').filter({hasText:'Discuss'}).first().click();await page.locator('#action-preview').waitFor({state:'visible'});await page.locator('#action-discuss').click();await page.waitForFunction(()=>!document.querySelector('#send-conversation').disabled);
  const actionDraft=await page.locator('#conversation-prompt').inputValue();assert.ok(actionDraft.includes('Discuss the selected record'));
  await send(actionDraft+'\n\n'+safety+'Use factory_user_ask to ask exactly "Which season do you prefer?" as a text question. Use factory_user_wait until answered. Then reply exactly "Recorded preference: <answer>".');
  await questions.locator('textarea').waitFor({timeout:180000});await focus('.question-card');
  await check('Question is pending and answer is initially empty',async()=>assert.equal(await questions.locator('textarea').inputValue(),''));
  await questions.locator('textarea').fill('Autumn');await page.waitForTimeout(2000);await questions.getByRole('button',{name:'Submit answer',exact:true}).click();
  await page.waitForFunction(()=>document.querySelector('#conversation-messages').textContent.includes('Recorded preference: Autumn'),null,{timeout:180000});await waitIdle();
  threadID=(await snapshot()).thread_id;
  await check('Real agent consumes submitted answer',async()=>assert.match(await page.locator('#conversation-messages').innerText(),/Recorded preference: Autumn/));await focus('#conversation-messages');
 });
 await scene('Generated form proposes a draft',['U077','U081','U082'],'The real agent presents a form. Explicit session trust mounts its isolated frame; the form can draft an answer but cannot submit it.',async()=>{
  await send(safety+'Use factory_user_ask for the text question "What should we name the reading session?". Then use factory_ui_present title "Reading session form", fallback "Type a session name in the plain answer above.", request_id equal to that question id, and generate your own self-contained attractive HTML form. Include input id session-name and a button id draft-session labeled Copy to answer draft. The button must parent.postMessage({type:\'factoryctl:draft\',text:document.getElementById(\'session-name\').value},\'*\'). Do not submit the answer directly. Then factory_user_wait until answered and reply "Session name: <answer>".');
  await page.locator('.generated-card').waitFor({timeout:180000});
  await check('Generated content is initially unmounted; plain fallback available',async()=>{assert.equal(await page.locator('iframe.generated-frame').count(),0);assert.match(await page.locator('.view-fallback').innerText(),/session name/);});
  await page.getByRole('button',{name:'Allow generated views for this session',exact:true}).click();
  const frame=page.frameLocator('iframe.generated-frame').first();await frame.locator('#session-name').fill('Autumn reading');await page.locator('iframe.generated-frame').first().scrollIntoViewIfNeeded();await frame.locator('#session-name').evaluate(e=>e.style.outline='3px solid #d99d27');await page.waitForTimeout(6500);await frame.locator('#draft-session').click();
  const pending=questions.filter({hasText:'What should we name'});await page.waitForFunction(()=>[...document.querySelectorAll('.question-card textarea')].some(e=>e.value==='Autumn reading'));
  await check('Frame draft does not answer the pending question',async()=>{assert.equal(await pending.locator('textarea').inputValue(),'Autumn reading');assert.ok(await pending.getByRole('button',{name:'Submit answer',exact:true}).isVisible());});await focus('.draft-notice:not([hidden])');await pending.getByRole('button',{name:'Submit answer',exact:true}).click();
  await page.waitForFunction(()=>document.querySelector('#conversation-messages').textContent.includes('Session name: Autumn reading'),null,{timeout:180000});await waitIdle();
  await check('Form response returns to the same native conversation',async()=>assert.equal((await snapshot()).thread_id,threadID));
 });
 await scene('Explore a chart and add a custom answer',['U078','U081','U082'],'A distinct generated chart drafts a choice. The trusted host also accepts custom text; the agent receives only the explicit submission.',async()=>{
  await send(safety+'Ask exactly "When would you like to read?" with kind single_choice and choices [{"id":"morning","label":"Morning"},{"id":"evening","label":"Evening"}]. Present title "Reading time chart", fallback "Morning: 8 hours. Evening: 5 hours. Choose or type a custom answer above.", request_id that question id and generate your own self-contained horizontal bar chart with Morning 8 hours and Evening 5 hours. Make the Evening bar a button id evening that calls parent.postMessage({type:\'factoryctl:draft\',option_ids:[\'evening\']},\'*\'); make a Morning bar as well. Wait for the answer then repeat the submitted custom text verbatim in your final response.');
  await page.locator('iframe.generated-frame').nth(1).waitFor({timeout:180000});const frame=page.frameLocator('iframe.generated-frame').nth(1);await page.locator('iframe.generated-frame').nth(1).scrollIntoViewIfNeeded();await frame.locator('#evening').scrollIntoViewIfNeeded();await frame.locator('#evening').evaluate(e=>e.style.outline='3px solid #d99d27');await page.waitForTimeout(8000);await frame.locator('#evening').click();
  const q=questions.filter({hasText:'When would you like'});await page.waitForFunction(()=>document.querySelector('input[value="evening"]')?.checked);
  await q.locator('textarea').fill('Saturday afternoon instead');await q.scrollIntoViewIfNeeded();await q.locator('textarea').evaluate(e=>e.style.outline='3px solid #d99d27');await page.waitForTimeout(6500);
  await check('Chart selects a draft; custom answer stays editable',async()=>{assert.ok(await q.locator('input[value="evening"]').isChecked());assert.equal(await q.locator('textarea').inputValue(),'Saturday afternoon instead');});
  await q.getByRole('button',{name:'Submit answer',exact:true}).click();await page.waitForFunction(()=>document.querySelector('#conversation-messages').textContent.includes('Saturday afternoon instead'),null,{timeout:180000});await waitIdle();
 });
 await scene('Verify the generated frame boundary',['U070','U081','U082'],'Adversarial browser checks: a mounted view cannot read the host, its credential storage, or fetch data. Foreign-window drafts are ignored; plain controls remain usable.',async()=>{
  const frame=page.frames().find(f=>f.url().includes('/view/'));
  const checks=await frame.evaluate(async()=>{let parentBlocked=false,storageBlocked=false,fetchBlocked=false;try{void parent.document.body}catch{parentBlocked=true}try{void sessionStorage.length}catch{storageBlocked=true}try{await fetch('/api/documents')}catch{fetchBlocked=true}return {parentBlocked,storageBlocked,fetchBlocked};});
  await check('Opaque frame blocks host document, credential storage and network',async()=>assert.deepEqual(checks,{parentBlocked:true,storageBlocked:true,fetchBlocked:true}));
  await check('Mounted resource cannot be opened as a top-level page',async()=>assert.equal((await page.request.get(frame.url())).status(),404));
  const draftBefore=await page.locator('#conversation-prompt').inputValue();
  await page.evaluate(()=>window.dispatchEvent(new MessageEvent('message',{origin:'null',source:window,data:{type:'factoryctl:draft',text:'spoofed answer'}})));
  await check('Foreign source-window cannot draft a response',async()=>assert.equal(await page.locator('#conversation-prompt').inputValue(),draftBefore));
  await frame.evaluate(()=>setTimeout(()=>{throw Error('Intentional qualification view failure')},0));
  await page.waitForFunction(()=>[...document.querySelectorAll('.view-notice')].some(e=>e.textContent.includes('reported an error')));
  await check('Generated code failure preserves plain fallback',async()=>assert.ok(await page.locator('.view-fallback').first().isVisible()));
  const navigationRequests=[];await context.route('https://factoryctl-blocked.invalid/**',route=>{navigationRequests.push(route.request().url());route.abort();});
  await frame.evaluate(()=>{location.href='https://factoryctl-blocked.invalid/navigation-check'}).catch(()=>{});await page.waitForTimeout(1000);
  await check('Frame self-navigation is blocked before external network dispatch',async()=>assert.deepEqual(navigationRequests,[]));
  await focus('.view-fallback');
 });
 await scene('Follow up without losing the conversation',['U079','U080'],'The next agent turn recalls the actual submitted season and session name on the same native thread.',async()=>{
  const before=await snapshot();
  await send(safety+'Without calling any tools, reply with the season and reading-session name I already submitted, prefixed FOLLOWUP_RECALL.');
  await page.waitForFunction(()=>document.querySelector('#conversation-messages').textContent.includes('FOLLOWUP_RECALL')&&document.querySelector('#conversation-status').textContent==='idle',null,{timeout:180000});
  const after=await snapshot();const newAssistant=after.messages.slice(before.messages.length).filter(m=>m.role==='assistant').map(m=>m.text).join('\n');
  await check('Fresh assistant turn recalls answers on exact same native thread',async()=>{assert.equal(after.thread_id,threadID);assert.equal(after.status,'idle');assert.match(newAssistant,/FOLLOWUP_RECALL/);assert.match(newAssistant,/Autumn/);assert.match(newAssistant,/Autumn reading/);});await focus('#conversation-messages');
 });
 await scene('Interrupt an unanswered request',['U076','U079'],'Stop turn cancels the unanswered request. The host waits for native interruption before allowing another message; cancellation is not consent.',async()=>{
  await send(safety+'Ask the text question "Should we schedule this now?" and wait for my explicit response before proceeding.');
  const q=questions.filter({hasText:'Should we schedule this now?'});await q.waitFor({timeout:180000});await page.waitForTimeout(2000);await page.click('#interrupt-conversation');
  await page.waitForFunction(()=>document.querySelector('#conversation-status').textContent==='interrupted',null,{timeout:60000});
  await check('Native interruption completes and question is cancelled',async()=>{assert.match(await q.innerText(),/cancelled/);assert.ok(await page.locator('#send-conversation').isEnabled());});await focus('#conversation-status');
 });

 await scene('Two native conversations stay independent',['U034','U076','U077','U079'],'Two real conversations own simultaneous pending questions. Each synthetic answer returns only to its requesting native thread.',async()=>{
  await send(safety+'Ask the text question "First concurrent question?" and wait. When answered reply exactly "FIRST_ACK <answer>".');
  await questions.filter({hasText:'First concurrent question?'}).locator('textarea').waitFor({timeout:180000});
  await page.click('#new-conversation');await page.waitForFunction(()=>!document.querySelector('#send-conversation').disabled);
  await send(safety+'Ask the text question "Second concurrent question?" and wait. When answered reply exactly "SECOND_ACK <answer>".');
  const second=questions.filter({hasText:'Second concurrent question?'});await second.locator('textarea').waitFor({timeout:180000});
  const readAll=()=>page.evaluate(async()=>{const headers={'X-Factory-Host':sessionStorage.getItem('factoryctl-host')};const l=await(await fetch('/api/conversations',{headers})).json();return await Promise.all(l.conversations.map(async c=>(await fetch('/api/conversations/'+c.id,{headers})).json()));});
  const pending=await readAll();assert.equal(pending.length,2);assert.equal(new Set(pending.map(c=>c.thread_id)).size,2);assert.equal(pending.filter(c=>c.questions.some(q=>q.status==='pending')).length,2);
  await second.locator('textarea').fill('BETA_ONLY_demo');await second.getByRole('button',{name:'Submit answer',exact:true}).click();await page.waitForFunction(()=>document.getElementById('conversation-messages').textContent.includes('SECOND_ACK BETA_ONLY_demo'));await waitIdle();
  await page.locator('.conversation-tab[aria-pressed="false"]').click();const first=questions.filter({hasText:'First concurrent question?'});await first.locator('textarea').waitFor();assert.equal(await first.locator('textarea').inputValue(),'');await first.locator('textarea').fill('ALPHA_ONLY_demo');await first.getByRole('button',{name:'Submit answer',exact:true}).click();await page.waitForFunction(()=>document.getElementById('conversation-messages').textContent.includes('FIRST_ACK ALPHA_ONLY_demo'));await waitIdle();
  const completed=await readAll();const firstThread=completed.find(c=>c.thread_id===threadID);assert(firstThread.messages.some(m=>m.role==='assistant'&&m.text.includes('FIRST_ACK ALPHA_ONLY_demo')));const secondThread=completed.find(c=>c.thread_id!==threadID);assert(secondThread.messages.some(m=>m.role==='assistant'&&m.text.includes('SECOND_ACK BETA_ONLY_demo')));
  await check('Both real pending scopes receive only their correlated answers',async()=>assert.equal(completed.every(c=>c.status==='idle'),true));
 });

 assert.deepEqual(errors.filter(e=>!e.includes('Intentional qualification view failure')),[]);assert.deepEqual(external.filter(u=>!u.startsWith('https://factoryctl-blocked.invalid/')),[]);
 assert.deepEqual(inventory(path.join(root,'factory'),f=>f.endsWith('.md')),before,'Browser source inventory or bytes changed during capture');
 provenance.product_files_after_capture=productManifest();assert.deepEqual(provenance.product_files_after_capture,productBefore,'Product source changed during capture');
 assert.equal(hash(binary),provenance.binary_sha256,'Owned binary changed during capture');
 assert.ok(server.exitCode===null&&server.signalCode===null,'Owned server exited during capture');
 await page.screenshot({path:path.join(dir,'review-final.png')});
 await context.close();await browser.close();const raw=await video.path();
 await stopServer();provenance.owned_server_stopped=true;
 cp.execFileSync('ffmpeg',['-v','error','-i',raw,'-c:v','libx264','-preset','fast','-crf','24','-pix_fmt','yuv420p','-movflags','+faststart','-an',path.join(dir,'demo.mp4')]);
 const media=JSON.parse(cp.execFileSync('ffprobe',['-v','error','-show_entries','format=duration,size:stream=codec_name,width,height,nb_frames','-of','json',path.join(dir,'demo.mp4')],{encoding:'utf8'}));
 fs.writeFileSync(path.join(dir,'results.json'),JSON.stringify({revision,base,provenance,server_output:{stdout:serverOutput.stdout.replace(/#host=[a-f0-9]+/g,"#host=<redacted>"),stderr:serverOutput.stderr},scenes,page_errors:errors,external_requests:external,source_files_unchanged:true,source_files:before},null,2));
 fs.writeFileSync(path.join(dir,'coverage.json'),JSON.stringify({revision,scenes:scenes.map(({commands,...scene})=>scene)},null,2));
 fs.writeFileSync(path.join(dir,'media.json'),JSON.stringify(media,null,2));
 console.log(JSON.stringify({directory:dir,revision,scenes:scenes.length,checks:scenes.reduce((n,s)=>n+s.commands.length,0),duration:media.format.duration,sha256:hash(path.join(dir,'demo.mp4'))}));
} catch(e){if(page){await page.screenshot({path:path.join(dir,'failure.png')}).catch(()=>{});fs.writeFileSync(path.join(dir,'failure-ui.txt'),await page.locator('#conversation-panel').innerText().catch(()=>''));}fs.writeFileSync(path.join(dir,'failure.json'),JSON.stringify({revision,provenance,server_output:{stdout:serverOutput.stdout.replace(/#host=[a-f0-9]+/g,"#host=<redacted>"),stderr:serverOutput.stderr},error:String(e),scenes,page_errors:errors,external_requests:external},null,2));throw e;} finally {
 await context?.close().catch(()=>{});await browser?.close().catch(()=>{});await stopServer();fs.rmSync(temporary,{recursive:true,force:true});
}
