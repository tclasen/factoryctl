// Check evidence accounting; semantic E2E sufficiency still needs independent review.
import fs from 'node:fs';
import path from 'node:path';
import crypto from 'node:crypto';
import {spawnSync} from 'node:child_process';
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
  const annotations=JSON.parse(fs.readFileSync(path.join(dir,matrix.annotations),'utf8'));
  const stable=path.resolve(dir,matrix.stable_video);
  const canonical='https://raw.githubusercontent.com/tclasen/factoryctl/main/artifacts/demo/latest.mp4';
  demand(matrix.stable_url===canonical && annotations.stable_link===canonical,'Canonical demo link changed');
  demand(matrix.stable_video==='../../demo/latest.mp4','Canonical demo file path changed');
  demand(fs.existsSync(stable),'Canonical latest demo is absent');
  if(fs.existsSync(video) && fs.existsSync(stable)) {
    const archiveHash=crypto.createHash('sha256').update(fs.readFileSync(video)).digest('hex');
    const latestHash=crypto.createHash('sha256').update(fs.readFileSync(stable)).digest('hex');
    demand(archiveHash===latestHash,'Canonical demo differs from immutable archive');
    demand(archiveHash===annotations.output_sha256 && latestHash===annotations.stable_alias_sha256,'Annotated video digest mismatch');
  }
  demand(annotations.product_revision===matrix.product_revision,'Annotation revision mismatch');
  demand(annotations.annotations.length===recording.scenes.length && new Set(annotations.annotations.map(a=>a.scene_id)).size===recording.scenes.length,'Not every demo scene has requirement annotations');
  for(const scene of recording.scenes) {
    const a=annotations.annotations.find(a=>a.scene_id===scene.id);
    demand(a && a.end_seconds>a.start_seconds && a.explanation && a.requirement_ids?.length,`Missing requirement explanation for scene ${scene.id}`);
    if(a) {
      let validEnd=Math.abs(a.end_seconds-scene.end_seconds)<1;
      // Context shutdown can append a brief final hold after the presentation clock.
      // Permit only the final scene to cover that hold, and verify the actual media end.
      if(!validEnd && scene===recording.scenes.at(-1)) {
        const tail=a.end_seconds-scene.end_seconds;
        if(tail>=0 && tail<=2 && fs.existsSync(video)) {
          const probe=spawnSync('ffprobe',['-v','error','-show_entries','format=duration','-of','json',video],{encoding:'utf8',timeout:10000});
          let exportedDuration=NaN;
          try { exportedDuration=Number(JSON.parse(probe.stdout).format.duration); } catch {}
          validEnd=probe.status===0 && Number.isFinite(exportedDuration) && Math.abs(a.end_seconds-exportedDuration)<0.05;
        }
      }
      demand(Math.abs(a.start_seconds-scene.start_seconds)<1 && validEnd,`Annotation timing differs from actual scene ${scene.id}`);
      for(const id of a.requirement_ids)demand(expected.includes(id),`Unknown on-video requirement ${id}`);
    }
  }
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
    for(const name of row.tests || []) { demand(testNames.has(name), `Unknown named test ${name} for ${row.id}`); demand(observedTests.has(name),`No executed passing named test ${name} for ${row.id}`); }
    if(row.status==='pass') {
      demand(row.coverage==='end-to-end' && ['executed-cli','executed-browser','executed-integration','verified-delivery'].includes(row.evidence_kind), `Pass lacks declared E2E evidence ${row.id}`);
      demand(row.remaining_gap==='', `Pass retains unmet requirement ${row.id}`);
    } else demand(typeof row.remaining_gap==='string' && row.remaining_gap.length>0, `Unfulfilled row lacks explicit gap ${row.id}`);
  }
  for(const id of expected)demand(seen.has(id), `Missing matrix row ${id}`);
  const unmet=matrix.rows.filter(r=>r.status!=='pass').map(r=>r.id);
  const summary={schema_version:1,product_revision:matrix.product_revision,row_count:expected.length,accounting_complete:failures.length===0,fully_qualified:failures.length===0&&unmet.length===0,pass_rows:matrix.rows.filter(r=>r.status==='pass').map(r=>r.id),unmet_rows:unmet,status_counts:Object.fromEntries([...new Set(matrix.rows.map(r=>r.status))].map(s=>[s,matrix.rows.filter(r=>r.status===s).length])),main_video:matrix.main_video,canonical_demo_link:canonical,main_video_sha256:fs.existsSync(video)?crypto.createHash('sha256').update(fs.readFileSync(video)).digest('hex'):null,failures,note:'Complete accounting is not full requirement satisfaction. Named backend cases and boundary refusals do not pass unimplemented E2E requirements.'};
  if(args.includes('--write')) {
    fs.writeFileSync(file,JSON.stringify(matrix,null,2)+'\n');
    fs.writeFileSync(path.join(dir,'matrix-summary.json'),JSON.stringify(summary,null,2)+'\n');
  }
  console.log(JSON.stringify(summary,null,2));
  process.exitCode=failures.length || (qualify && unmet.length) ? 1 : 0;
} catch(error) { console.error('matrix_evidence_error: '+error.message); process.exitCode=1; }
