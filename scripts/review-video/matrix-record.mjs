import fs from 'node:fs';
import path from 'node:path';
import os from 'node:os';
import crypto from 'node:crypto';
import cp from 'node:child_process';
import assert from 'node:assert/strict';
import {createRequire} from 'node:module';
const require = createRequire(import.meta.url);
const root = path.resolve(import.meta.dirname, '../..');
const out = path.resolve(process.env.REVIEW_OUTPUT_DIR || path.join(root, 'artifacts/cli-review/848b4ec'));
const input = path.resolve(process.env.MATRIX_INPUT || path.join(root, 'artifacts/cli-review/848b4ec/requirements-coverage.json'));
const outputs = ['matrix-review.webm', 'matrix-review.mp4', 'matrix-video-coverage.json', 'matrix-poster.png'];
for (const name of outputs) assert.ok(!fs.existsSync(path.join(out, name)), `Refusing overwrite: ${name}; choose a fresh REVIEW_OUTPUT_DIR`);
const bytes = fs.readFileSync(input), matrix = JSON.parse(bytes);
const hash = data => crypto.createHash('sha256').update(data).digest('hex');
assert.equal(matrix.rows.length, 94);
assert.equal(new Set(matrix.rows.map(r => r.id)).size, 94);
assert.ok(matrix.rows.every((r, i) => r.id === `U${String(i + 1).padStart(3, '0')}`));
const passes = matrix.rows.filter(r => r.status === 'pass' && r.coverage === 'end-to-end').length;
assert.ok(passes >= 0 && passes <= matrix.rows.length);
const scratch = fs.mkdtempSync(path.join(os.tmpdir(), 'factoryctl-matrix-video-'));
fs.mkdirSync(out, {recursive: true});
function command(executable, args) {
  const r = cp.spawnSync(executable, args, {encoding: 'utf8', timeout: 900000});
  if (r.error || r.status !== 0) throw new Error(`${executable}: ${r.error || r.stderr}`);
  return r.stdout;
}
const ffmpeg = process.env.FFMPEG || 'ffmpeg', ffprobe = process.env.FFPROBE || 'ffprobe';
const {chromium} = require(process.env.PLAYWRIGHT_MODULE || 'playwright');
const browser = await chromium.launch({headless: true, ...(process.env.REVIEW_CHROMIUM ? {executablePath: process.env.REVIEW_CHROMIUM} : {})});
try {
  // Verify the actual browser encoder and exported file before the long capture.
  const smoke = await browser.newContext({recordVideo: {dir: path.join(scratch, 'smoke')}});
  const sp = await smoke.newPage(); await sp.setContent('<h1>Matrix encoder smoke check</h1>'); await sp.waitForTimeout(1000);
  const sv = sp.video(); await smoke.close(); const smokePath = await sv.path();
  assert.ok(fs.statSync(smokePath).size > 1000);
  assert.ok(Number(command(ffprobe, ['-v', 'error', '-show_entries', 'format=duration', '-of', 'default=nw=1:nk=1', smokePath])) > 0);
  console.log('PASS exported browser-video encoder smoke');
  if (process.argv.includes('--smoke-only')) {await browser.close(); process.exit(0);}
  const context = await browser.newContext({viewport: {width: 1440, height: 1000}, deviceScaleFactor: 1, recordVideo: {dir: path.join(scratch, 'raw'), size: {width: 1440, height: 1000}}});
  const layoutOnly = process.argv.includes('--layout-only');
  const origin = Date.now();
  const page = await context.newPage();
  await page.setContent(`<!doctype html><meta charset="utf-8"><style>*{box-sizing:border-box}body{margin:0;background:#0b1221;color:#eef3ff;font-family:Arial,sans-serif}header{padding:24px 42px;background:#15223a;height:130px;border-bottom:2px solid #b98b42}h1{font-size:30px;margin:8px 0}.eyebrow{font-size:17px;letter-spacing:1px;color:#ffcc7a}.meta{position:absolute;right:42px;top:28px;font:16px/1.7 monospace;color:#becbe1;text-align:right}main{padding:22px 42px}#caption{font-size:21px;color:#bdcce5;height:62px;line-height:1.4}.row{border:2px solid #354661;border-radius:10px;background:#132037;padding:15px 18px;margin:10px 0;min-height:151px;transition:border-color .4s,background .4s}.row.active{border-color:#ffcc7a;background:#1e2d44}.line{display:flex;gap:18px;align-items:baseline}.id{font: bold 23px monospace;color:#8fcfff}.requirement{font-size:22px;line-height:1.3;flex:1}.status{font-size:16px;color:#ffd28f;white-space:nowrap}.pass{color:#8ce0bf}.gap{font-size:19px;line-height:1.35;margin-top:10px;color:#e0e7f3}.refs{font-size:16px;color:#91afd4;margin-top:8px}.intro{padding:45px 30px;font-size:27px;line-height:1.6}.intro strong{color:#ffcc7a}footer{position:fixed;bottom:0;background:#15223a;width:100%;height:64px;padding:18px 42px;font-size:17px;color:#b8c8e2;display:flex;justify-content:space-between}#progress{position:fixed;bottom:0;height:4px;background:#e9b966;transition:width .7s}</style><header><div class="eyebrow">PART 2 · REQUIREMENTS ACCOUNTING · FULL MATRIX INCOMPLETE</div><h1>Accounting, not functional demonstration</h1><div class="meta">factoryctl · 0.1.0-dev<br>${matrix.product_revision.slice(0,12)}</div></header><main><div id="caption"></div><div id="rows"></div></main><footer><span>Part 1: cli-review.mp4 · executed CLI results and assertions · evaluation panel, not product GUI</span><span id="counter"></span></footer><div id="progress"></div>`);
  const intro = async (caption, lines) => page.evaluate(({caption, lines}) => {document.querySelector('#caption').textContent=caption;const el=document.querySelector('#rows');el.replaceChildren();const d=document.createElement('div');d.className='intro';for(const text of lines){const p=document.createElement('p');p.textContent=text;d.append(p);}el.append(d);}, {caption, lines});
  await intro('Read alongside Part 1 and requirements-coverage.json. No new product functionality is tested here.', [`${matrix.rows.length} requirements accounted for. ${passes} mapped end-to-end passes; ${matrix.rows.length-passes} requirements have unmet or partial coverage.`, 'The following pages show every exact requirement, its recorded status, and remaining gap. CLI scene numbers refer to Part 1.', 'GUI, MCP and agent runtime are not implemented. Unsupported refusals demonstrate boundaries; they do not pass those capability requirements.']);
  await page.screenshot({path: path.join(out, 'matrix-poster.png')});
  await page.waitForTimeout(layoutOnly ? 50 : 10000);
  const pages = [], rows = [];
  for (let offset=0; offset<matrix.rows.length; offset+=4) {
    const group=matrix.rows.slice(offset,offset+4), index=offset/4+1;
    await page.evaluate(({group,index})=>{
      document.querySelector('#caption').textContent='Requirement → current status → remaining gap. Referenced scenes demonstrate only the implemented CLI portion.';
      document.querySelector('#counter').textContent=`PAGE ${index} / 24`;
      document.querySelector('#progress').style.width=`${index/24*100}%`;
      const list=document.querySelector('#rows');list.replaceChildren();
      for(const r of group){const row=document.createElement('div');row.className='row';
        const line=document.createElement('div');line.className='line';
        for(const [cls,text] of [['id',r.id],['requirement',r.requirement],['status'+(r.status==='pass'?' pass':''),r.status+' · '+r.coverage]]){const d=document.createElement('div');d.className=cls;d.textContent=text;line.append(d);}row.append(line);
        const gap=document.createElement('div');gap.className='gap';gap.textContent=r.remaining_gap?'GAP · '+r.remaining_gap:'MAPPED PASS · '+r.observed;row.append(gap);
        const refs=document.createElement('div');refs.className='refs';refs.textContent=r.scene_ids.length?'Part 1 CLI scene references: '+r.scene_ids.join(', ')+' · '+r.evidence_kind:'No functional scene · '+r.evidence_kind;row.append(refs);list.append(row);
      }
    },{group,index});
    await page.waitForTimeout(150);
    // Fail rather than silently crop long requirement or gap text.
    assert.ok(await page.evaluate(()=>document.querySelector('#rows').getBoundingClientRect().bottom<935),'Matrix page must fit above footer');
    if (layoutOnly) await page.screenshot({path:path.join(scratch, `layout-${index}.png`)});
    const start=(Date.now()-origin)/1000;
    for(let i=0;i<group.length;i++) {await page.locator('.row').evaluateAll((nodes,active)=>nodes.forEach((n,j)=>n.classList.toggle('active',j===active)),i);await page.waitForTimeout(layoutOnly ? 50 : 2250);}
    if(group.length<4) await page.waitForTimeout(layoutOnly ? 50 : (4-group.length)*2250);
    const end=(Date.now()-origin)/1000;
    pages.push({page:index,start_seconds:start,end_seconds:end,ids:group.map(r=>r.id)});
    rows.push(...group.map(r=>({id:r.id,page:index,start_seconds:start,end_seconds:end,status:r.status,coverage:r.coverage,evidence_kind:'matrix-accounting-only',functional_video:matrix.main_video,scene_ids:r.scene_ids})));
    console.log(`RECORDED matrix page ${index}/24 ${group[0].id}–${group.at(-1).id} ${end.toFixed(1)}s`);
  }
  await intro('Full-matrix qualification remains incomplete.', [`All ${matrix.rows.length} rows are accounted for; ${passes} have mapped end-to-end CLI passes. The other ${matrix.rows.length-passes} remain partial or unmet.`, 'Use cli-review.mp4 with coverage.json and results.json for executed CLI evidence. Use requirements-coverage.json for complete gaps and named regression cases.', 'This addendum is a review of the matrix, not an end-to-end demonstration of missing GUI, MCP, runtime or release capabilities.']);
  await page.waitForTimeout(layoutOnly ? 50 : 10000);
  const video=page.video();await context.close();await video.saveAs(path.join(out,'matrix-review.webm'));
  command(ffmpeg,['-hide_banner','-loglevel','error','-i',path.join(out,'matrix-review.webm'),'-c:v','libx264','-preset','slow','-crf','30','-pix_fmt','yuv420p','-movflags','+faststart','-an',path.join(out,'matrix-review.mp4')]);
  const duration=Number(command(ffprobe,['-v','error','-show_entries','format=duration','-of','default=nw=1:nk=1',path.join(out,'matrix-review.mp4')]));
  assert.ok(fs.statSync(path.join(out,'matrix-review.mp4')).size<=10000000,'MP4 must be at most 10MB');
  assert.ok(duration>rows.at(-1).end_seconds);
  fs.writeFileSync(path.join(out,'matrix-video-coverage.json'),JSON.stringify({schema_version:1,layout_check_only:layoutOnly,product_revision:matrix.product_revision,matrix_input:path.basename(input),matrix_sha256:hash(bytes),recorder_sha256:hash(fs.readFileSync(import.meta.filename)),created_at:new Date().toISOString(),qualification:'incomplete',purpose:'accounting, not functional demonstration',row_count:rows.length,end_to_end_passes:passes,partial_or_unmet:rows.length-passes,main_video:matrix.main_video,video:'matrix-review.mp4',duration_seconds:duration,timestamp_basis:'Wall clock from recorded page creation; page intervals exclude transition settling. Browser capture startup may differ by a few frames.',toolchain:{node:process.version,playwright:require((process.env.PLAYWRIGHT_MODULE||'playwright')+'/package.json').version,chromium:browser.version()},artifacts:outputs.filter(x=>x!=='matrix-video-coverage.json').map(name=>({path:name,sha256:hash(fs.readFileSync(path.join(out,name))),bytes:fs.statSync(path.join(out,name)).size})),pages,rows},null,2)+'\n');
  console.log(`DONE ${duration}s ${out}; temporary frames: ${scratch}`);
} finally {await browser.close();}
