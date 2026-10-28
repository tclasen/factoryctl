// Add requirement annotations to the immutable, continuous Playwright capture.
import fs from 'node:fs';
import path from 'node:path';
import crypto from 'node:crypto';
import cp from 'node:child_process';
import assert from 'node:assert/strict';
import {createRequire} from 'node:module';
const require=createRequire(import.meta.url);
const root=path.resolve(import.meta.dirname,'../..');
const dir=path.resolve(process.env.REVIEW_OUTPUT_DIR||path.join(root,'artifacts/cli-review/848b4ec'));
const source=path.join(dir,'cli-review.mp4');
const target=path.join(dir,'demo.mp4');
const alias=path.join(root,'artifacts/demo/latest.mp4');
const skipAlias=process.env.REVIEW_SKIP_ALIAS==='1';
assert.ok(!fs.existsSync(target),'Archived demo already exists; select a new output directory');
const hash=p=>crypto.createHash('sha256').update(fs.readFileSync(p)).digest('hex');
const sourceHash=hash(source);
const run=(args)=>{const r=cp.spawnSync('ffmpeg',args,{encoding:'utf8',maxBuffer:8*1024*1024});if(r.status!==0)throw Error(r.stderr);return r;};
const probe=p=>JSON.parse(cp.execFileSync('ffprobe',['-v','error','-show_entries','format=duration,size:stream=codec_name,width,height,nb_frames','-of','json',p],{encoding:'utf8'}));
const sourceProbe=probe(source),duration=Number(sourceProbe.format.duration);
const coverage=JSON.parse(fs.readFileSync(path.join(dir,'coverage.json')));
const matrix=JSON.parse(fs.readFileSync(path.join(dir,'requirements-coverage.json')));
const rows=new Map(matrix.rows.map(r=>[r.id,r]));
// Each explanation is intentionally limited to the actual command scene.
const descriptions=[
 [['U045','U048','U084'],'CLI PORTION / CURRENT MAC','Build, version and doctor run locally; portability is unqualified; web conversations are shown in the companion segment.'],
 [['U004','U008','U050'],'DEMONSTRATED CLI SCOPE','Initialize one explicit Git project with software defaults; repeated init preserves existing files.'],
 [['U004','U015','U050'],'DEMONSTRATED CLI SCOPE','Discover one project and create a real Markdown Task from its starter template.'],
 [['U058','U030'],'PARTIAL / CLI WRITE SAFETY','Missing template input and an existing destination fail safely; no GUI editing is shown.'],
 [['U026','U027'],'PARTIAL / CLI LISTING','Filter and paginate CLI records; the separate web review covers browser interactions.'],
 [['U015','U026','U029'],'PARTIAL / CLI READ MODEL','Read exact Markdown, revision and link metadata; this is CLI output, not a product detail view.'],
 [['U015','U030','U059'],'PARTIAL / CLI EDITING','Update a valid state and preserve unrelated bytes; structural validity does not grant approval.'],
 [['U058','U030'],'PARTIAL / CLI WRITE SAFETY','Reject stale revisions and invalid states without changing bytes; reject an expired list cursor.'],
 [['U015','U030'],'PARTIAL / CLI REPLACEMENT','Replace the complete Markdown document from stdin using the current revision; this scene exercises CLI editing.'],
 [['U029','U058'],'PARTIAL / LINKS AND SAFE DELETION','Resolve a backlink, reject referenced-record deletion, then delete the incoming record safely.'],
 [['U050','U021','U032'],'PARTIAL / TEMPLATES AND PROMPTS','Inspect starter templates and applicable prompts; runtime conversations are covered in the web segment.'],
 [['U025','U058'],'PARTIAL / STRUCTURAL VALIDATION','Run both validation levels; external links remain unchecked and strict completeness fails honestly.'],
 [['U058','U030'],'PARTIAL / CLI WRITE BOUNDARIES','Reject traversal and coordination-anchor deletion; the protected file remains on disk.'],
 [['U005','U029','U024'],'S2 DEMONSTRATED / OTHER ROWS PARTIAL','Two projects share one Git root; CLI links and backlinks resolve across independent local definitions.'],
 [['U005','U024','U058'],'S2 DEMONSTRATED / OTHER ROWS PARTIAL','Select alpha, enforce its local enum and prove beta stays byte-identical after accepted/rejected edits.'],
 [['U005','U022','U058'],'CUSTOM TYPES DEMONSTRATED / VALIDATION PARTIAL','Create a custom Widget from its template; reject removal of a heading referenced by another record.'],
 [['U006','U029','U026'],'S3 DEMONSTRATED / OTHER ROWS PARTIAL','Discover sibling Git repositories; CLI cross-project links and backlinks resolve using their local definitions.'],
 [['U006','U024','U058'],'S3 DEMONSTRATED / OTHER ROWS PARTIAL','Choose one sibling project, enforce its schema and preserve the other project unchanged.'],
 [['U006','U022','U058'],'CUSTOM TYPES DEMONSTRATED / VALIDATION PARTIAL','Create a new type without rebuilding; preserve the target heading when proposed replacement is invalid.'],
 [['U004','U008','U030'],'PATH SCOPE DEMONSTRATED / EDITING PARTIAL','Select a duplicate checkout; replace complex YAML explicitly and verify CWD versus explicit-path scope.'],
 [['U067','U058'],'MCP DOCUMENT TOOLS / HEADLESS BOUNDARY','Initialize real MCP; list tools, read and write with revisions, reject stale writes and headless questions.'],
 [['U001','U031'],'CAPABILITY BOUNDARIES / WEB SHOWN SEPARATELY','Headless chat refuses; legacy gui, actions run, view export and network validation remain unsupported.'],
 [['U048','U058','U092'],'CURRENT MAC TESTS / FULL QUALIFICATION INCOMPLETE','Fresh tests, race tests and vet pass on macOS; backend results do not establish full end-to-end coverage.'],
 [['U084','U088','U094'],'BOUNDED CLI DEMO / FULL MATRIX INCOMPLETE','Requirement IDs explain the scenes shown here; external MCP Apps and full release qualification remain incomplete.']
];
assert.equal(descriptions.length,coverage.scenes.length);
// Detect actual title changes in the exported source, then verify them against the scene clock.
const detection=run(['-hide_banner','-i',source,'-vf',"crop=1320:60:40:120,select='gt(scene,0.08)',showinfo",'-an','-f','null','-']);
// The initial blank-to-title paint can create a transition within the first 200ms.
// Scene one already starts at zero; retain every subsequent detected title change.
const starts=[0,...[...detection.stderr.matchAll(/pts_time:([0-9.]+)/g)].map(m=>Number(m[1])).filter(t=>t>0.2)];
assert.equal(starts.length,coverage.scenes.length,'Every source title transition must be identified');
const annotations=coverage.scenes.map((s,i)=>{
 assert.ok(Math.abs(starts[i]-s.start_seconds)<0.2,'Source title timing diverges from executed scene');
 const [ids,label,explanation]=descriptions[i];
 return {scene_id:s.id,start_seconds:starts[i],end_seconds:starts[i+1]??duration,requirement_ids:ids,label,explanation,coverage:i===21?'boundary':([1,2].includes(i)?'full':'partial'),original_caption:s.caption,requirements:ids.map(id=>{assert.ok(rows.has(id),id+' missing from matrix');const r=rows.get(id);return {id,requirement:r.requirement,matrix_status:matrix.product_revision===coverage.revision?r.status:'not-currently-qualified',reference_matrix_status:r.status,remaining_gap:r.remaining_gap};})};
});
const assets=path.join(dir,'annotation-bands');fs.mkdirSync(assets,{recursive:true});
const {chromium}=require(process.env.PLAYWRIGHT_MODULE||'playwright');
const browser=await chromium.launch({headless:true,...(process.env.REVIEW_CHROMIUM?{executablePath:process.env.REVIEW_CHROMIUM}:{})});
try{
 const page=await browser.newPage({viewport:{width:1440,height:100},deviceScaleFactor:1});
 await page.setContent(`<style>*{box-sizing:border-box}body{margin:0;background:#111b31;color:#edf1fa;font-family:Arial,sans-serif;border-bottom:2px solid #5edcbd;height:100px;padding:10px 32px}.top{display:flex;gap:18px;align-items:baseline;white-space:nowrap}#ids{font-size:25px;font-weight:700;color:#7ce7cb}#label{font-size:17px;color:#c7d5eb;letter-spacing:.5px}#explanation{font-size:21px;margin-top:5px;white-space:nowrap}#context{font-size:13px;margin-top:7px;color:#95adc9;letter-spacing:.2px}</style><div class="top"><span id="ids"></span><span id="label"></span></div><div id="explanation"></div><div id="context">REQUIREMENT ANNOTATION · CLI review harness / recorded executions · 0.1.0-dev · ${coverage.revision.slice(0,12)} · Full requirements matrix incomplete</div>`);
 for(const a of annotations){
  await page.evaluate(a=>{document.querySelector('#ids').textContent=a.requirement_ids.join(' · ');document.querySelector('#label').textContent=a.label;document.querySelector('#explanation').textContent=a.explanation;},a);
  assert.ok(await page.locator('#explanation').evaluate(el=>el.scrollWidth<=el.clientWidth),'Explanation must fit header without clipping');
  assert.ok(await page.locator('.top').evaluate(el=>el.scrollWidth<=el.clientWidth),'IDs and scope label must fit header');
  await page.screenshot({path:path.join(assets,`scene-${String(a.scene_id).padStart(2,'0')}.png`)});
 }
}finally{await browser.close();}
const concat=['ffconcat version 1.0'];for(const a of annotations){concat.push(`file 'scene-${String(a.scene_id).padStart(2,'0')}.png'`,`duration ${(a.end_seconds-a.start_seconds).toFixed(3)}`);}concat.push(`file 'scene-${String(annotations.at(-1).scene_id).padStart(2,'0')}.png'`);fs.writeFileSync(path.join(assets,'bands.ffconcat'),concat.join('\n')+'\n');
run(['-hide_banner','-loglevel','error','-n','-i',source,'-f','concat','-safe','0','-i',path.join(assets,'bands.ffconcat'),'-filter_complex','[0:v][1:v]overlay=0:0:eof_action=repeat[v]','-map','[v]','-t',String(duration),'-c:v','libx264','-preset','slow','-crf','24','-pix_fmt','yuv420p','-movflags','+faststart','-an',target]);
run(['-v','error','-i',target,'-f','null','-']);
const targetProbe=probe(target);assert.equal(targetProbe.format.duration,sourceProbe.format.duration,'No cuts or timing changes');assert.equal(targetProbe.streams[0].nb_frames,sourceProbe.streams[0].nb_frames,'Frame count must be preserved');assert.equal(hash(source),sourceHash,'Source capture must remain immutable');
if(!skipAlias){fs.mkdirSync(path.dirname(alias),{recursive:true});const tempAlias=alias+'.tmp';fs.copyFileSync(target,tempAlias);fs.renameSync(tempAlias,alias);assert.equal(hash(alias),hash(target),'Stable alias must be byte-identical');}
fs.writeFileSync(path.join(dir,'annotations.json'),JSON.stringify({schema_version:1,product_revision:coverage.revision,source:'cli-review.mp4',source_sha256:sourceHash,source_video_hash:sourceHash,output:'demo.mp4',output_sha256:hash(target),outputhash:hash(target),stable_alias:skipAlias?null:'artifacts/demo/latest.mp4',stable_alias_sha256:skipAlias?null:hash(alias),stable_link:skipAlias?null:'/Users/Shared/factoryctl/artifacts/demo/latest.mp4',alias_updated:!skipAlias,requirements_reference_revision:matrix.product_revision,requirements_reference_scope:matrix.product_revision===coverage.revision?'current mapping':'historical text reference only; current row qualification not claimed',duration_seconds:duration,frames:targetProbe.streams[0].nb_frames,annotation_region:{x:0,y:0,width:1440,height:100},timing:'Actual exported title transitions; original frames and pacing preserved',scope:'Scene-specific requirement evidence; partial and unmet rows are not promoted to pass',annotations},null,2)+'\n');
console.log(JSON.stringify({target,alias,sha256:hash(target),probe:targetProbe},null,2));
