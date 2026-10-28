import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import {spawnSync} from 'node:child_process';
import crypto from 'node:crypto';
const root=path.resolve(import.meta.dirname,'../..');
// Synthetic structural fixtures: no dependency on local, untracked recordings.
// The mock probe is test-only; actual delivery separately decodes real media.
const hash=b=>crypto.createHash('sha256').update(b).digest('hex');
function fixture(t) {
 const base=fs.mkdtempSync(path.join(os.tmpdir(),'factoryctl-matrix-check-'));
 t.after(()=>fs.rmSync(base,{recursive:true,force:true}));
 const dir=path.join(base,'artifacts/review/fixture');fs.mkdirSync(dir,{recursive:true});fs.mkdirSync(path.join(base,'artifacts/demo'),{recursive:true});fs.mkdirSync(path.join(dir,'bin'));
 const write=(name,data)=>fs.writeFileSync(path.join(dir,name),JSON.stringify(data));
 const video=Buffer.from('SYNTHETIC STRUCTURAL TEST; NOT DELIVERABLE VIDEO');fs.writeFileSync(path.join(dir,'demo.mp4'),video);fs.writeFileSync(path.resolve(dir,'../../demo/latest.mp4'),video);
 fs.writeFileSync(path.join(dir,'bin/ffprobe'),`#!/bin/sh
printf '{"format":{"duration":"3.5"}}\\n'
`,{mode:0o755});
 const rows=fs.readFileSync(path.join(root,'factory/specifications/requirements-traceability.md'),'utf8').split('\n').filter(l=>/^\| U\d+ \|/.test(l)).map(l=>{const c=l.split('|').map(s=>s.trim());return {id:c[1],requirement:c[3],status:'implemented-unverified',coverage:'partial',evidence_kind:'executed-browser',observed:'Synthetic structural fixture only',remaining_gap:'Not demonstrated by this test',scene_ids:[1],tests:[]};});
 const canonical='/Users/Shared/factoryctl/artifacts/demo/latest.mp4';
 write('requirements-coverage.json',{product_revision:'synthetic-test',main_video:'demo.mp4',results:'results.json',recording_coverage:'coverage.json',annotations:'annotations.json',stable_video:'../../demo/latest.mp4',stable_url:canonical,rows});
 write('results.json',{revision:'synthetic-test',scenes:[{id:1,commands:[{exit:0,expected:0}]},{id:2,commands:[{exit:0,expected:0}]}]});
 write('coverage.json',{revision:'synthetic-test',scenes:[{id:1,start_seconds:0,end_seconds:1},{id:2,start_seconds:1,end_seconds:2}]});
 write('annotations.json',{product_revision:'synthetic-test',stable_link:canonical,output_sha256:hash(video),stable_alias_sha256:hash(video),annotations:[{scene_id:1,start_seconds:0,end_seconds:1,requirement_ids:['U001'],explanation:'Fixture scene1'},{scene_id:2,start_seconds:1,end_seconds:3.5,requirement_ids:['U002'],explanation:'Fixture closing hold'}]});
 const events='{"Action":"pass","Test":"TestIntentDoesNotInferAcceptance"}\n';fs.writeFileSync(path.join(dir,'test-events.jsonl'),events);write('test-run.json',{exit_code:0,product_revision:'synthetic-test',events_sha256:hash(events)});
 const product=[];function walk(p){for(const e of fs.readdirSync(p,{withFileTypes:true})){const f=path.join(p,e.name);if(e.isDirectory())walk(f);else if(/\.(go|html|css|js)$/.test(e.name))product.push({path:path.relative(root,f),sha256:hash(fs.readFileSync(f))});}}walk(path.join(root,'cmd'));walk(path.join(root,'internal'));for(const name of ['go.mod','go.sum'])product.push({path:name,sha256:hash(fs.readFileSync(path.join(root,name)))});write('toolchain.json',{product_files:product});
 return dir;
}
function change(dir,name,fn){const p=path.join(dir,name),v=JSON.parse(fs.readFileSync(p,'utf8'));fn(v);fs.writeFileSync(p,JSON.stringify(v));}
function run(dir,...args){return spawnSync(process.execPath,[path.join(root,'scripts/review-video/check-matrix.mjs'),path.join(dir,'requirements-coverage.json'),...args],{encoding:'utf8',timeout:15000,env:{...process.env,PATH:path.join(dir,'bin')+path.delimiter+process.env.PATH}});}
test('complete accounting does not authorize incomplete E2E qualification',t=>{const dir=fixture(t);const good=run(dir);assert.equal(good.status,0,good.stderr+good.stdout);assert.equal(JSON.parse(good.stdout).accounting_complete,true);const gated=run(dir,'--qualify');assert.equal(gated.status,1);assert.equal(JSON.parse(gated.stdout).fully_qualified,false);assert.ok(JSON.parse(gated.stdout).unmet_rows.length>0);});
test('missing and duplicate rows are rejected',t=>{const dir=fixture(t);let removedID;change(dir,'requirements-coverage.json',m=>{removedID=m.rows.pop().id;m.rows.push(m.rows[0]);});const r=run(dir);assert.equal(r.status,1);assert.ok(r.stdout.includes('Missing matrix row '+removedID),r.stdout);assert.match(r.stdout,/Duplicate row U001/);});
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
test('final tail must match actual exported duration',t=>{const dir=fixture(t);change(dir,'annotations.json',m=>{m.annotations.at(-1).end_seconds+=0.2;});const r=run(dir);assert.equal(r.status,1);assert.match(r.stdout,/Annotation timing differs from actual scene 2/);});
test('final annotation start remains strict despite the export tail',t=>{const dir=fixture(t);change(dir,'annotations.json',m=>{m.annotations.at(-1).start_seconds+=1.5;});const r=run(dir);assert.equal(r.status,1);assert.match(r.stdout,/Annotation timing differs from actual scene 2/);});

test('embedded browser asset changes invalidate old capture',t=>{const dir=fixture(t);change(dir,'toolchain.json',m=>{m.product_files.find(f=>f.path==='internal/web/assets/app.js').sha256='0'.repeat(64);});const r=run(dir);assert.equal(r.status,1);assert.match(r.stdout,/Stale product evidence: internal\/web\/assets\/app.js/);});

function qualification(dir) {
 const files=[];function walk(p){for(const e of fs.readdirSync(p,{withFileTypes:true})){const f=path.join(p,e.name);if(e.isDirectory())walk(f);else if(e.name.endsWith('.go'))files.push({path:path.relative(root,f),sha256:hash(fs.readFileSync(f))});}}walk(path.join(root,'scripts/qualification'));
 change(dir,'toolchain.json',m=>{m.qualification_files=files;});
 const events=fs.readFileSync(path.join(dir,'test-events.jsonl'),'utf8')+'{"Action":"pass","Test":"TestSyntheticKnowledgeWorkCorpus"}\n';fs.writeFileSync(path.join(dir,'test-events.jsonl'),events);change(dir,'test-run.json',m=>{m.events_sha256=hash(events);});
 change(dir,'requirements-coverage.json',m=>{m.rows[0].tests=['TestSyntheticKnowledgeWorkCorpus'];});
}
test('qualification tests need source binding and actual passing events',t=>{const dir=fixture(t);qualification(dir);const r=run(dir);assert.equal(r.status,0,r.stdout+r.stderr);change(dir,'toolchain.json',m=>{delete m.qualification_files;});const absent=run(dir);assert.equal(absent.status,1);assert.match(absent.stdout,/Unknown named test TestSyntheticKnowledgeWorkCorpus/);});
test('qualification source changes invalidate its test evidence',t=>{const dir=fixture(t);qualification(dir);change(dir,'toolchain.json',m=>{m.qualification_files[0].sha256='0'.repeat(64);});const r=run(dir);assert.equal(r.status,1);assert.match(r.stdout,/Stale qualification evidence/);});
test('qualification source inventories must be complete',t=>{const dir=fixture(t);qualification(dir);change(dir,'toolchain.json',m=>{m.qualification_files.pop();});const r=run(dir);assert.equal(r.status,1);assert.match(r.stdout,/Qualification source inventory changed/);});
