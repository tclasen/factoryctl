// Check evidence accounting; semantic E2E sufficiency still needs independent review.
import fs from 'node:fs';
import path from 'node:path';
import crypto from 'node:crypto';
const root = path.resolve(import.meta.dirname, '../..');
const args = process.argv.slice(2);
const qualify = args.includes('--qualify');
const input = args.find(a => !a.startsWith('--')) || 'artifacts/cli-review/848b4ec/requirements-coverage.json';
const file = path.resolve(root, input);
const dir = path.dirname(file);
const failures = [];
const demand = (condition, message) => { if (!condition) failures.push(message); };
try {
  const matrix = JSON.parse(fs.readFileSync(file, 'utf8'));
  const spec = fs.readFileSync(path.join(root, 'factory/specifications/requirements-traceability.md'), 'utf8');
  const sourceRows = spec.split('\n').filter(l => /^\| U\d+ \|/.test(l)).map(l => l.split('|').map(c => c.trim()));
  const expected = sourceRows.map(r => r[1]);
  const sourceText = new Map(sourceRows.map(r => [r[1], r[3]]));
  const results = JSON.parse(fs.readFileSync(path.join(dir, matrix.results), 'utf8'));
  const recording = JSON.parse(fs.readFileSync(path.join(dir, matrix.recording_coverage), 'utf8'));
  const video = path.join(dir, matrix.main_video);
  demand(fs.existsSync(video) && fs.statSync(video).size > 0, 'Main Playwright video is absent or empty');
  demand(results.revision === matrix.product_revision && recording.revision === matrix.product_revision, 'Product revisions differ');
  demand(results.scenes.every(s => s.commands.every(c => c.exit === c.expected)), 'An executed command has an unexpected outcome');
  const testNames = new Set();
  const observedTests = new Set();
  const events = fs.readFileSync(path.join(dir,'test-events.jsonl'),'utf8').trim().split('\n').map(l=>JSON.parse(l));
  for(const event of events) { if(event.Action==='pass' && event.Test) observedTests.add(event.Test); demand(event.Action!=='fail','Current named-test run contains a failure'); }
  const testRun=JSON.parse(fs.readFileSync(path.join(dir,'test-run.json'),'utf8'));
  demand(testRun.exit_code===0 && testRun.product_revision===matrix.product_revision,'Named-test provenance mismatch');
  demand(testRun.events_sha256===crypto.createHash('sha256').update(fs.readFileSync(path.join(dir,'test-events.jsonl'))).digest('hex'),'Named-test event digest mismatch');
  const toolchain=JSON.parse(fs.readFileSync(path.join(dir,'toolchain.json'),'utf8'));
  const productPaths=[];
  const inventory=p=>{for(const e of fs.readdirSync(p,{withFileTypes:true})){const f=path.join(p,e.name);if(e.isDirectory())inventory(f);else if(e.name.endsWith('.go'))productPaths.push(path.relative(root,f));}};
  inventory(path.join(root,'cmd'));inventory(path.join(root,'internal'));productPaths.push('go.mod','go.sum');
  demand(JSON.stringify(productPaths.sort())===JSON.stringify(toolchain.product_files.map(f=>f.path).sort()),'Product source inventory changed since capture');
  for(const f of toolchain.product_files)demand(fs.existsSync(path.join(root,f.path)) && crypto.createHash('sha256').update(fs.readFileSync(path.join(root,f.path))).digest('hex')===f.sha256,'Stale product evidence: '+f.path);
  const scan = p => { for (const e of fs.readdirSync(p, {withFileTypes:true})) { const f=path.join(p,e.name); if(e.isDirectory())scan(f); else if(e.name.endsWith('_test.go')) for(const m of fs.readFileSync(f,'utf8').matchAll(/^func (Test\w+)\(/gm))testNames.add(m[1]); }};
  scan(path.join(root,'cmd')); scan(path.join(root,'internal'));
  const matrixVideo=JSON.parse(fs.readFileSync(path.join(dir,'matrix-video-coverage.json'),'utf8'));
  const matrixInput=fs.readFileSync(path.join(dir,'matrix-input.json'));
  demand(crypto.createHash('sha256').update(matrixInput).digest('hex')===matrixVideo.matrix_sha256,'Frozen matrix video input digest mismatch');
  demand(matrixVideo.product_revision===matrix.product_revision,'Matrix video revision mismatch');
  demand(matrixVideo.rows.length===expected.length && new Set(matrixVideo.rows.map(r=>r.id)).size===expected.length,'Matrix video row coverage is incomplete');
  demand(fs.existsSync(path.join(dir,matrixVideo.video)),'Matrix review video is absent');
  const frozenRows=new Map(JSON.parse(matrixInput).rows.map(r=>[r.id,r]));
  const seen = new Set();
  for (const row of matrix.rows) {
    demand(expected.includes(row.id), `Unknown row ${row.id}`);
    demand(!seen.has(row.id), `Duplicate row ${row.id}`); seen.add(row.id);
    demand(['pass','planned','implemented-unverified','fail','blocked','deprecated-and-removed'].includes(row.status), `Invalid status ${row.id}`);
    demand(typeof row.requirement === 'string' && row.requirement.length > 0, `Missing requirement ${row.id}`);
    demand(row.requirement===sourceText.get(row.id),`Stale requirement text ${row.id}`);
    demand(row.observed && row.evidence_kind && row.coverage, `Missing observed evidence/disposition ${row.id}`);
    demand(Array.isArray(row.scene_ids) && row.scene_ids.length > 0, `No scene/gap mapping ${row.id}`);
    const ranges=[];
    for (const id of row.scene_ids || []) {
      const scene=recording.scenes.find(s => s.id === id);
      demand(scene && Number.isFinite(scene.start_seconds) && scene.end_seconds > scene.start_seconds, `Missing valid scene ${id} for ${row.id}`);
      if(scene)ranges.push({scene_id:id,start_seconds:scene.start_seconds,end_seconds:scene.end_seconds,kind:row.evidence_kind});
    }
    row.video_ranges=ranges;
    const shown=matrixVideo.rows.find(r=>r.id===row.id);
    demand(shown && shown.end_seconds>shown.start_seconds && shown.status===row.status,`No matching matrix video row ${row.id}`);
    const frozen=frozenRows.get(row.id);
    demand(frozen && frozen.requirement===row.requirement && frozen.status===row.status && frozen.remaining_gap===row.remaining_gap,`Matrix video input is stale for ${row.id}`);
    for(const name of row.tests || []) { demand(testNames.has(name), `Unknown named test ${name} for ${row.id}`); demand(observedTests.has(name),`No executed passing named test ${name} for ${row.id}`); }
    if(row.status==='pass') {
      demand(row.coverage==='end-to-end' && row.evidence_kind==='executed-cli', `Pass lacks declared E2E evidence ${row.id}`);
      demand(row.remaining_gap==='', `Pass retains unmet requirement ${row.id}`);
    } else demand(typeof row.remaining_gap==='string' && row.remaining_gap.length>0, `Unfulfilled row lacks explicit gap ${row.id}`);
  }
  for(const id of expected)demand(seen.has(id), `Missing matrix row ${id}`);
  const unmet=matrix.rows.filter(r=>r.status!=='pass').map(r=>r.id);
  const summary={schema_version:1,product_revision:matrix.product_revision,row_count:expected.length,accounting_complete:failures.length===0,fully_qualified:failures.length===0&&unmet.length===0,pass_rows:matrix.rows.filter(r=>r.status==='pass').map(r=>r.id),unmet_rows:unmet,status_counts:Object.fromEntries([...new Set(matrix.rows.map(r=>r.status))].map(s=>[s,matrix.rows.filter(r=>r.status===s).length])),main_video:matrix.main_video,main_video_sha256:fs.existsSync(video)?crypto.createHash('sha256').update(fs.readFileSync(video)).digest('hex'):null,failures,note:'Complete accounting is not full requirement satisfaction. Named backend cases and boundary refusals do not pass unimplemented E2E requirements.'};
  if(args.includes('--write')) {
    fs.writeFileSync(file,JSON.stringify(matrix,null,2)+'\n');
    fs.writeFileSync(path.join(dir,'matrix-summary.json'),JSON.stringify(summary,null,2)+'\n');
  }
  console.log(JSON.stringify(summary,null,2));
  process.exitCode=failures.length || (qualify && unmet.length) ? 1 : 0;
} catch(error) { console.error('matrix_evidence_error: '+error.message); process.exitCode=1; }
