// Record actual read-only browser interactions; all output stays in ignored artifacts/.
import fs from 'node:fs';
import path from 'node:path';
import crypto from 'node:crypto';
import cp from 'node:child_process';
import {createRequire} from 'node:module';
import assert from 'node:assert/strict';
const require=createRequire(import.meta.url);
const {chromium}=require(process.env.PLAYWRIGHT_MODULE||'playwright');
const root=path.resolve(import.meta.dirname,'../..');
const dir=path.resolve(process.env.REVIEW_OUTPUT_DIR||path.join(root,'artifacts/intent-review/20261028'));
assert.ok(!fs.existsSync(path.join(dir,'demo.mp4'))&&!fs.existsSync(path.join(dir,'results.json')),'Use a fresh archive directory; never overwrite evidence');
fs.mkdirSync(dir,{recursive:true});
const base=process.env.REVIEW_BASE_URL||'http://127.0.0.1:8080';
assert.ok(['127.0.0.1','localhost','[::1]'].includes(new URL(base).hostname),'Record only a local server');
const revision=cp.execFileSync('git',['rev-parse','HEAD'],{cwd:root,encoding:'utf8'}).trim();
const hash=file=>crypto.createHash('sha256').update(fs.readFileSync(file)).digest('hex');
const files=[];function inventory(p){for(const e of fs.readdirSync(p,{withFileTypes:true})){const f=path.join(p,e.name);if(e.isDirectory())inventory(f);else files.push({path:path.relative(root,f),sha256:hash(f)});}}
inventory(path.join(root,'factory'));const before=files.filter(f=>f.path.endsWith('.md'));
const browser=await chromium.launch({headless:true,...(process.env.REVIEW_CHROMIUM?{executablePath:process.env.REVIEW_CHROMIUM}:{})});
const context=await browser.newContext({viewport:{width:1440,height:1000},recordVideo:{dir:path.join(dir,'raw'),size:{width:1440,height:1000}}});
const page=await context.newPage();const video=page.video();const errors=[];const external=[];
page.on('pageerror',e=>errors.push(String(e)));page.on('request',r=>{if(!r.url().startsWith(base))external.push(r.url());});
const scenes=[];let began;
async function check(label,fn){await fn();scenes.at(-1).commands.push({command:label,exit:0,expected:0,evidence_kind:'executed-browser'});}
async function scene(title,ids,explanation,fn){
 const start=(Date.now()-began)/1000;const item={id:scenes.length+1,title,start_seconds:start,requirement_ids:ids,explanation,commands:[]};scenes.push(item);
 await page.evaluate(({title,ids,explanation})=>{const box=document.getElementById('review-annotation');box.replaceChildren();const h=document.createElement('strong');h.textContent=ids.join(' · ')+' — '+title;const p=document.createElement('div');p.textContent=explanation;box.append(h,p);},{title,ids,explanation});
 await page.waitForTimeout(1800);await fn();await page.waitForTimeout(6500);item.end_seconds=(Date.now()-began)/1000;
}
async function focus(selector){const node=page.locator(selector).first();await node.scrollIntoViewIfNeeded();await node.evaluate(e=>{e.style.outline='3px solid #d99d27';e.style.outlineOffset='4px';});await page.waitForTimeout(1800);await node.evaluate(e=>{e.style.outline='';});}
async function selectDocument(query,title){await page.locator('#search').fill(query);await page.locator('#document-list .document-link').filter({hasText:title}).first().click();await page.waitForFunction(t=>document.querySelector('#document-title').textContent===t,title);await page.waitForTimeout(1800);}
try {
 await page.goto(base);await page.waitForFunction(()=>document.querySelectorAll('.document-link').length>0);
 await page.evaluate(()=>{const box=document.createElement('div');box.id='review-annotation';Object.assign(box.style,{position:'fixed',top:'0',left:'0',right:'0',height:'86px',boxSizing:'border-box',padding:'14px 24px',background:'#142d28',color:'#fff',zIndex:'9999',font:'17px/1.6 system-ui'});document.body.prepend(box);document.querySelector('.topbar').style.marginTop='86px';document.querySelector('.workspace').style.height='calc(100vh - 150px)';});
 began=Date.now();
 await scene('Explore real planning documents',['U007','U026','U100'],'Actual product UI: search and open local Markdown. This is a read-only browser slice, not full GUI qualification.',async()=>{
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
 await scene('Keep future designs proposed',['U031','U096','U100'],'The MCP Apps protocol remains a proposal to qualify. This explorer does not claim MCP Apps integration.',async()=>{
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
 await scene('Return to one current plan',['U097','U099','U100'],'Follow the current-plan link. Reading and watching submit no acceptance; editing, conversations and MCP Apps remain outstanding.',async()=>{
  await page.locator('#document-body a').filter({hasText:'current plan'}).click();await page.waitForFunction(()=>document.querySelector('#document-title').textContent==='Current implementation plan');
  await check('Current plan links are navigable',async()=>assert.match(await page.locator('#document-body').textContent(),/Current work and proposed continuation/));
  await page.locator('#document-body h2').filter({hasText:'Current work and proposed continuation'}).scrollIntoViewIfNeeded();await page.waitForTimeout(3500);
  await check('Read-only surface exposes no approval or save action',async()=>assert.equal(await page.getByRole('button',{name:/^(approve|accept|save)$/i}).count(),0));
 });
 assert.deepEqual(errors,[]);assert.deepEqual(external,[]);
 for(const f of before)assert.equal(hash(path.join(root,f.path)),f.sha256,'Browser changed source '+f.path);
 await page.screenshot({path:path.join(dir,'review-final.png')});
 await context.close();await browser.close();const raw=await video.path();
 cp.execFileSync('ffmpeg',['-v','error','-i',raw,'-c:v','libx264','-preset','fast','-crf','24','-pix_fmt','yuv420p','-movflags','+faststart','-an',path.join(dir,'demo.mp4')]);
 const media=JSON.parse(cp.execFileSync('ffprobe',['-v','error','-show_entries','format=duration,size:stream=codec_name,width,height,nb_frames','-of','json',path.join(dir,'demo.mp4')],{encoding:'utf8'}));
 fs.writeFileSync(path.join(dir,'results.json'),JSON.stringify({revision,base,scenes,page_errors:errors,external_requests:external,source_files_unchanged:true,source_files:before},null,2));
 fs.writeFileSync(path.join(dir,'coverage.json'),JSON.stringify({revision,scenes:scenes.map(({commands,...scene})=>scene)},null,2));
 fs.writeFileSync(path.join(dir,'media.json'),JSON.stringify(media,null,2));
 console.log(JSON.stringify({directory:dir,revision,scenes:scenes.length,checks:scenes.reduce((n,s)=>n+s.commands.length,0),duration:media.format.duration,sha256:hash(path.join(dir,'demo.mp4'))}));
} catch(e){fs.writeFileSync(path.join(dir,'failure.json'),JSON.stringify({revision,error:String(e),scenes,page_errors:errors,external_requests:external},null,2));await context.close();await browser.close();throw e;}
