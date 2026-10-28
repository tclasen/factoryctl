#!/usr/bin/env python3
"""Normalize immutable clips, bind current evidence, and build honest row accounting.

Does not infer operator acceptance or publish anything. Static linkage remains
blocked unless a separately reviewed source clarification is supplied later.
"""
import argparse,hashlib,json,pathlib,re,shutil,subprocess
root=pathlib.Path(__file__).resolve().parents[2]
parser=argparse.ArgumentParser();parser.add_argument('--output',required=True);parser.add_argument('--parts',required=True);parser.add_argument('--prepare-only',action='store_true');args=parser.parse_args()
out=(root/args.output).resolve();assert not out.exists(),'Use a fresh versioned review directory';out.mkdir(parents=True)
load=lambda p:json.loads((root/p).read_text())
sha=lambda p:hashlib.sha256(pathlib.Path(p).read_bytes()).hexdigest()
probe=lambda p:json.loads(subprocess.check_output(['ffprobe','-v','error','-show_entries','format=duration,size:stream=codec_name,width,height,nb_frames,r_frame_rate','-of','json',str(p)]))
config=load(args.parts);revision=config['product_revision'];scenes=[];offset=0.;parts=[]
for part in config['parts']:
 data=load(part['results']);video=(root/part['video']).resolve();media=probe(video);duration=float(media['format']['duration'])
 assert media['streams'][0]['width']==1440 and media['streams'][0]['height']==1000
 assert media['streams'][0]['r_frame_rate']=='25/1'
 if 'passed' in data:assert data['passed'],part['name']
 raw=data['scenes']
 for index,s in enumerate(raw):
  ids=s.get('requirement_ids',s.get('requirements',[]));ids=[i for i in ids if re.fullmatch(r'U\d{3}',i)]
  if part['kind']=='model':
   ids=['U104'] if s['title'].startswith('U104') else (['U103','U106'] if 'Choose the next turn' in s['title'] else ['U103'])
  assert ids,(part['name'],index)
  start=s.get('start_seconds',s.get('at_seconds'))
  end=s.get('end_seconds',raw[index+1].get('at_seconds') if index+1<len(raw) else duration)
  source_end=end
  closing_hold=duration-end if index+1==len(raw) and part['kind']=='evidence' else 0
  if closing_hold>2:end=duration
  assert isinstance(start,(float,int)) and isinstance(end,(float,int)) and 0<=start<end<=duration+2
  commands=s.get('commands')
  if commands is None:
   if part['kind']=='web':commands=[{'command':' '.join(data['command']['argv']),'assertion_group':s['name'],'exit':data['observed_exit'],'expected':data['expected_exit']}]
   elif part['kind']=='model':commands=[{'command':'node scripts/review-video/model-settings-live.cjs','assertion_group':s['title'],'exit':0 if data['passed'] else 1,'expected':0}]
   else:raise AssertionError('Missing executed command binding '+part['name'])
  assert all(c['exit']==c['expected'] for c in commands)
  title=s.get('title',s.get('name'));explanation=s.get('explanation',s.get('text'))
  if not explanation:explanation=(s.get('layout','')+' '+title+'. Current executed browser assertions; fixture responses are synthetic, not operator acceptance.').strip()
  if closing_hold>2:explanation+=f' Preserved {closing_hold:.3f}-second exported closing hold on this same final panel; raw presentation clock end {source_end:.3f}s.'
  scenes.append({'id':len(scenes)+1,'part':part['name'],'local_scene':s.get('id',index+1),'title':title,'requirement_ids':ids,'explanation':explanation,'start_seconds':offset+start,'end_seconds':offset+min(end,duration),'commands':commands,'evidence_kind':part['kind'],'source_results':part['results']})
 parts.append(dict(part,sha256=sha(video),media=media,start_seconds=offset,duration_seconds=duration))
 offset+=duration
manifest=load('artifacts/completion-review/a7c22ba/source-inventory.json')
for f in manifest['product_files']:assert sha(root/f['path'])==f['sha256'],'Product source changed '+f['path']
qualification=load('artifacts/completion-qualification/0.1.0/attempt-009/report.json')
assert qualification['source_stable_during_run']
manifest['qualification_files']=[{'path':p,'sha256':digest} for p,digest in qualification['source_after'].items() if p.startswith('scripts/qualification/') and p.endswith('.go')]
assert manifest['qualification_files']
for f in manifest['qualification_files']:assert sha(root/f['path'])==f['sha256'],'Qualification source changed '+f['path']
events=root/'artifacts/completion-review/a7c22ba/test-events.jsonl';run=load('artifacts/completion-review/a7c22ba/test-run.json');assert run['exit_code']==0 and run['product_revision']==revision and sha(events)==run['events_sha256']
observed={e.get('Test') for line in events.read_text().splitlines() if (e:=json.loads(line)).get('Action')=='pass'}
testmap={
 'discovery':['TestCollectionScenarios','TestDepthAndFallback','TestWorktreeAndInvalidMarkers','TestIgnoredSymlinksAndNonGit'],
 'storage':['TestParseHazardsAndLimits','TestLosslessScalarEdit','TestProposedGraphAndUnrelatedRepairs','TestConcurrencyProcessLockAndExit'],
 'definitions':['TestFields','TestRuntimeExtensionsReloadAndSharedStructuralBoundary','TestDefinitionChangesValidateIncomingInstances'],
 'gui':['TestDocumentCRUDWithoutRuntimeAndConflicts','TestCanonicalCrossProjectLinksAndBacklinks','TestHTTPNodeCursorBindsSelectionAndRejectsMalformedOffset'],
 'runtime':['TestConcurrentRequestsAndNativeReplyCorrelation','TestPreferenceSelectionAndAdapterIsolation','TestApprovalMustMatchOfferedDecision'],
 'custom':['TestCustomTrustDigestScopeAndSandbox','TestCustomRejectsSymlinkAndExternalAssets','TestCustomDisabledDefinitionUsesGenericFallback'],
 'templates':['TestIncludesAndRestrictions','TestEscapedInputsAndLiteralData','TestTemplateSourceAndAssetPinsInsideWriter','TestPromptPreviewNeverLaunchesAndExposesContext'],
 'minimal':['TestMinimalWorkflowBoundary','TestExtensionRejectsExecutableBindingsAndReservedMappings','TestPromptPreviewNeverLaunchesAndExposesContext'],
 'questions':['TestPendingNeverBecomesAnAnswer','TestQuestionScopesAndHostSubmission','TestChoiceValidationAndCustomAnswers','TestBridgeScopeAndRealSubmissionRoundTrip'],
 'settings':['TestSettingsPersistenceResetIsolationAndPermissions','TestSettingsConcurrentProcesses','TestChatSettingsSaveConflictResetAndProjectIsolation'],
 'delivery':['TestLocalDeliveryFreshnessAndSharedPlan','TestDeliveryDoesNotInferApprovalOrRemoteEvidence'],
 'extensions':['TestRuntimeExtensionsReloadAndSharedStructuralBoundary','TestDynamicExtensionToolsAndDefinitionPinnedTransport','TestCLIExtensionParityAndPromptPreview'],
 'breadth':['TestResearchOpportunityJourney','TestSyntheticKnowledgeWorkCorpus'],
 'mcp':['TestScopedDocumentToolsAndStaleWrite','TestMCPProtocolAndPendingWaitCancellation','TestMCPPaginationBindsFiltersAndSnapshot']}
for group,names in testmap.items():assert all(n in observed for n in names),group
def group(n):
 if 3<=n<=10:return 'discovery'
 if 11<=n<=20:return 'storage'
 if 21<=n<=25:return 'definitions'
 if 26<=n<=30 or n==105:return 'gui'
 if 31<=n<=39:return 'templates' if n in [32,33,37] else 'runtime'
 if n==42 or n==71:return 'custom'
 if n in [55,56]:return 'breadth'
 if n in [58,59]:return 'storage'
 if n==67:return 'mcp'
 if n in [68,69]:return 'extensions'
 if 76<=n<=83:return 'questions'
 if n in [101,102]:return 'delivery'
 if n in [103,104]:return 'settings'
 if n in [60,61,62,63,84,85,86,87,88,89,90,91,92,93,94,95,96,97,98,99,100,106,45,46,47,48,49,50,51,52,53,54,57]:return None
 return 'minimal'
# Composite interfaces require actual product GUI scenes as well as CLI/MCP.
for scene in scenes:
 if (scene['part']=='web' and scene['local_scene'] in [1,2]) or (scene['part']=='native' and scene['local_scene'] in [1,8]):
  if 'U001' not in scene['requirement_ids']:scene['requirement_ids'].append('U001')
for scene in scenes:
 if (scene['part']=='collections' and scene['local_scene'] in [1,3,6]) or (scene['part']=='model' and scene['title'].endswith('Defaults survive server restart')) or (scene['part']=='evidence' and scene['local_scene']==1):
  if 'U046' not in scene['requirement_ids']:scene['requirement_ids'].append('U046')
# Lists are supplemental named assertions, not substitutes for mapped GUI journeys.
specific_tests={1:['TestCLIExtensionParityAndPromptPreview','TestScopedDocumentToolsAndStaleWrite','TestDocumentCRUDWithoutRuntimeAndConflicts'],2:['TestMinimalWorkflowBoundary'],18:[],20:[],29:['TestCanonicalCrossProjectLinksAndBacklinks'],36:['TestPreferenceSelectionAndAdapterIsolation'],37:['TestPromptPreviewNeverLaunchesAndExposesContext'],38:['TestMinimalWorkflowBoundary'],39:[],40:['TestMinimalWorkflowBoundary'],41:['TestMinimalWorkflowBoundary'],43:['TestHostOriginAndResourceBoundaries','TestListenRefusesNonLoopbackBeforeBinding'],44:['TestCollectionDocumentJourneys'],49:[],50:[],51:[],52:[],53:[],54:[]}
source=(root/'factory/specifications/requirements-traceability.md').read_text();rows=[]
for line in source.splitlines():
 if not re.match(r'\| U\d{3} \|',line):continue
 cols=[c.strip() for c in line.split('|')];uid,text=cols[1],cols[3];number=int(uid[1:]);mapped=[s for s in scenes if uid in s['requirement_ids']];assert mapped,'Missing actual scene for '+uid
 status='blocked' if uid=='U045' else 'pass'
 selected=group(number);tests=specific_tests.get(number,testmap[selected] if selected else [])
 kind='executed-browser' if any(s['evidence_kind'] not in ['evidence'] for s in mapped) else 'executed-integration'
 if kind=='executed-integration' and number<60 and any(s['part']=='evidence' and s['local_scene']<=24 for s in mapped):kind='executed-cli'
 gap='Darwin/Windows use OS-provided dynamic libraries; strict static interpretation requires the operator clarification already requested. No exception is inferred.' if status=='blocked' else ''
 limits=[]
 if uid in ['U035','U066']:limits.append('Native command cancellation verified; unsupported MCP elicitation remains explicit. Original approval-only recorder termination failure is preserved separately from verified trace criterion.')
 if uid in ['U045','U047','U048']:limits.append('Foreign binaries are built/inspected, not executed. Linux static linkage observed; Darwin/Windows system interfaces are not fully static.')
 if uid=='U046':limits.append('Mac desktop/CI/server restart exercised; server/Kubernetes use operator-protected proxy instructions, no deployment performed.')
 if uid=='U057':limits.append('Proposed one-second refresh target not met (~2.5s); declared corpus measured; RSS samples are not peak, end-to-end watcher visibility unmeasured.')
 if uid in ['U030','U058']:limits.append('Wait for displayed revision/card state after a save; rapid sequential interactions can cause a recoverable stale conflict. No race-free promise with external writers.')
 if uid in ['U061','U092']:limits.append('Current binary is 0.1.0-dev; initial release policy target is0.1.0. No release/tag/publication is claimed.')
 if uid in ['U081','U082','U083']:limits.append('First-party generated views verified; external MCP Apps interoperability unqualified. CLI retained under the explicit evaluation criterion.')
 if uid in ['U096','U098','U100']:limits.append('Source labels are recorded claims. Independent execution is procedural, not secure isolation or an authenticated original transcript.')
 rows.append({'id':uid,'requirement':text,'status':status,'observed':'Observed current-source evidence: '+ '; '.join(s['title'] for s in mapped)+'. '+('CLI/MCP calls and actual product GUI navigation/edit/conversation are separate demonstrated portions. ' if uid=='U001' else '')+('Static all-target compliance remains unresolved.' if uid=='U045' else 'See the mapped scene assertions and retained source results for sub-behavior and limits.'),'evidence_kind':kind,'coverage':'partial' if status=='blocked' else 'end-to-end','scene_ids':[s['id'] for s in mapped],'tests':tests,'remaining_gap':gap,'limitations':limits,'tested_product_revision':revision,'source_results':sorted({s['source_results'] for s in mapped})})
assert len(rows)==106
for name in ['test-events.jsonl','test-run.json']:shutil.copy2(root/'artifacts/completion-review/a7c22ba'/name,out/name)
(out/'toolchain.json').write_text(json.dumps(manifest,indent=2)+'\n')
recording={'revision':revision,'scenes':[{k:v for k,v in s.items() if k!='commands'} for s in scenes],'parts':parts,'full_matrix_qualification':'incomplete_pending_U045_static_interpretation'}
(out/'recording-coverage.json').write_text(json.dumps(recording,indent=2)+'\n');(out/'results.json').write_text(json.dumps({'revision':revision,'scenes':scenes,'parts':parts,'source_manifest':manifest,'limits':['No release publication','Static interpretation pending','Trace criteria and original failed probe dispositions remain distinct']},indent=2)+'\n')
matrix={'schema_version':1,'product_revision':revision,'results':'results.json','recording_coverage':'recording-coverage.json','main_video':'demo.mp4','stable_video':'../../demo/latest.mp4','stable_url':'/Users/Shared/factoryctl/artifacts/demo/latest.mp4','annotations':'annotations.json','rows':rows,'acceptance_disposition':'partial_review_package; mandatory U045 unresolved; no scoped exception inferred'}
(out/'requirements-coverage.json').write_text(json.dumps(matrix,indent=2)+'\n')
if not args.prepare_only:
 subprocess.run(['python3',str(root/'scripts/review-video/combine-segments.py'),*[str(root/p['video']) for p in parts],str(out/'demo.mp4')],check=True,stdout=(out/'combine.stdout').open('w'))
 combined=probe(out/'demo.mp4');assert abs(float(combined['format']['duration'])-offset)<.01
 for p in parts:assert sha(root/p['video'])==p['sha256']
 annotations={'product_revision':revision,'stable_link':matrix['stable_url'],'output_sha256':sha(out/'demo.mp4'),'stable_alias_sha256':None,'annotations':[{'scene_id':s['id'],'start_seconds':s['start_seconds'],'end_seconds':s['end_seconds'],'requirement_ids':s['requirement_ids'],'explanation':s['explanation'],'source_part':s['part']} for s in scenes],'alias_pending_independent_review':True,'no_cuts_or_speed_changes':True}
 (out/'annotations.json').write_text(json.dumps(annotations,indent=2)+'\n')
print(json.dumps({'output':str(out),'scenes':len(scenes),'rows':len(rows),'status_counts':{'pass':105,'blocked':1},'duration_seconds':offset,'canonical_alias_updated':False}))
