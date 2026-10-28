import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import {spawnSync} from 'node:child_process';
const root=path.resolve(import.meta.dirname,'../..');
const evidence=path.join(root,'artifacts/cli-review/848b4ec');
function fixture(t) {
 const base=fs.mkdtempSync(path.join(os.tmpdir(),'factoryctl-matrix-check-'));
 t.after(()=>fs.rmSync(base,{recursive:true,force:true}));
 const dir=path.join(base,'artifacts/cli-review/848b4ec');fs.mkdirSync(dir,{recursive:true});fs.mkdirSync(path.join(base,'artifacts/demo'),{recursive:true});
 for(const name of ['requirements-coverage.json','results.json','coverage.json','toolchain.json','test-events.jsonl','test-run.json','annotations.json'])fs.copyFileSync(path.join(evidence,name),path.join(dir,name));
 fs.symlinkSync(path.join(evidence,'demo.mp4'),path.join(dir,'demo.mp4'));
 fs.symlinkSync(path.join(root,'artifacts/demo/latest.mp4'),path.resolve(dir,'../../demo/latest.mp4'));
 return dir;
}
function change(dir,name,fn){const p=path.join(dir,name),v=JSON.parse(fs.readFileSync(p,'utf8'));fn(v);fs.writeFileSync(p,JSON.stringify(v));}
function run(dir,...args){return spawnSync(process.execPath,[path.join(root,'scripts/review-video/check-matrix.mjs'),path.join(dir,'requirements-coverage.json'),...args],{encoding:'utf8',timeout:15000});}
test('complete accounting does not authorize incomplete E2E qualification',t=>{const dir=fixture(t);const good=run(dir);assert.equal(good.status,0,good.stderr+good.stdout);assert.equal(JSON.parse(good.stdout).accounting_complete,true);const gated=run(dir,'--qualify');assert.equal(gated.status,1);assert.equal(JSON.parse(gated.stdout).fully_qualified,false);assert.ok(JSON.parse(gated.stdout).unmet_rows.length>0);});
test('missing and duplicate rows are rejected',t=>{const dir=fixture(t);change(dir,'requirements-coverage.json',m=>{m.rows.pop();m.rows.push(m.rows[0]);});const r=run(dir);assert.equal(r.status,1);assert.match(r.stdout,/Missing matrix row U095/);assert.match(r.stdout,/Duplicate row U001/);});
test('requirement text changes cannot reuse stale evidence',t=>{const dir=fixture(t);change(dir,'requirements-coverage.json',m=>{m.rows[0].requirement='Outdated requirement';});const r=run(dir);assert.equal(r.status,1);assert.match(r.stdout,/Stale requirement text U001/);});
test('source change cannot reuse old video qualification',t=>{const dir=fixture(t);change(dir,'toolchain.json',m=>{m.product_files[0].sha256='0'.repeat(64);});const r=run(dir);assert.equal(r.status,1);assert.match(r.stdout,/Stale product evidence/);});
test('missing video fails closed',t=>{const dir=fixture(t);fs.unlinkSync(path.join(dir,'demo.mp4'));const r=run(dir);assert.equal(r.status,1);assert.match(r.stdout,/video is absent or empty/);});
test('declared tests require actual passing events',t=>{const dir=fixture(t);change(dir,'requirements-coverage.json',m=>{m.rows[0].tests=['TestNeverExecuted'];});const r=run(dir);assert.equal(r.status,1);assert.match(r.stdout,/No executed passing named test TestNeverExecuted/);});
test('partial rows cannot become passes while retaining gaps',t=>{const dir=fixture(t);change(dir,'requirements-coverage.json',m=>{m.rows[0].status='pass';});const r=run(dir);assert.equal(r.status,1);assert.match(r.stdout,/Pass retains unmet requirement U001/);});

test('every real demo scene must carry requirement explanations',t=>{const dir=fixture(t);change(dir,'annotations.json',m=>{m.annotations.pop();});const r=run(dir);assert.equal(r.status,1);assert.match(r.stdout,/Not every demo scene has requirement annotations/);});
test('stable latest file must match the immutable demo',t=>{const dir=fixture(t);const p=path.resolve(dir,'../../demo/latest.mp4');fs.unlinkSync(p);fs.writeFileSync(p,'wrong video');const r=run(dir);assert.equal(r.status,1);assert.match(r.stdout,/Canonical demo differs/);});
test('single canonical demo link cannot silently change',t=>{const dir=fixture(t);change(dir,'requirements-coverage.json',m=>{m.stable_url='https://example.org/other.mp4';});const r=run(dir);assert.equal(r.status,1);assert.match(r.stdout,/Canonical demo link changed/);});

test('final annotation may cover only the verified short export tail',t=>{const dir=fixture(t);const a=JSON.parse(fs.readFileSync(path.join(dir,'annotations.json'),'utf8')).annotations.at(-1);const c=JSON.parse(fs.readFileSync(path.join(dir,'coverage.json'),'utf8')).scenes.at(-1);assert.ok(a.end_seconds-c.end_seconds>=1 && a.end_seconds-c.end_seconds<=2);const r=run(dir);assert.equal(r.status,0,r.stderr+r.stdout);});
test('intermediate annotation timing retains the strict bound',t=>{const dir=fixture(t);change(dir,'annotations.json',m=>{m.annotations[0].end_seconds+=1.5;});const r=run(dir);assert.equal(r.status,1);assert.match(r.stdout,/Annotation timing differs from actual scene 1/);});
test('final tail must match actual exported duration',t=>{const dir=fixture(t);change(dir,'annotations.json',m=>{m.annotations.at(-1).end_seconds+=0.2;});const r=run(dir);assert.equal(r.status,1);assert.match(r.stdout,/Annotation timing differs from actual scene 23/);});
test('final annotation start remains strict despite the export tail',t=>{const dir=fixture(t);change(dir,'annotations.json',m=>{m.annotations.at(-1).start_seconds+=1.5;});const r=run(dir);assert.equal(r.status,1);assert.match(r.stdout,/Annotation timing differs from actual scene 23/);});
