#!/usr/bin/env python3
"""Verify actual native trace criteria; preserve the original probe disposition."""
import json,pathlib,sys
core_path=pathlib.Path(sys.argv[1]);approval_path=pathlib.Path(sys.argv[2])
core=json.loads(core_path.read_text());approval=json.loads(approval_path.read_text())
assert core['passed'] and core['exit']==0 and core['sources_unchanged'] and core['binary_unchanged']
assert core['pty_isatty'] and core['same_native_thread'] and core['native_interrupted']
assert core['document_edit']['passed'] and core['document_edit']['only_title_bytes_changed']
assert core['document_edit']['native_revision_pinned_mcp_update_observed']
assert approval['binary_sha256']==core['binary_sha256'] and approval['sources_unchanged'] and approval['binary_unchanged']
events=[json.loads(line) for line in (approval_path.parent/'native-identities.jsonl').read_text().splitlines()]
requests=[e for e in events if e.get('direction')=='event' and e.get('method')=='item/commandExecution/requestApproval']
assert len(requests)==1
request=requests[0]
reply=next(e for e in events if e.get('direction')=='approval_reply' and e['request_id']==request['request_id'])
assert reply['decision'] in request['available_decisions'] and reply['decision'] in ['cancel','decline']
completion=next(e for e in events if e.get('direction')=='command_completion' and e.get('thread_id')==request['thread_id'] and e.get('turn_id')==request['turn_id'])
assert completion['status']=='declined' and completion['exit_code'] is None
assert any(e.get('method')=='turn/completed' and e.get('thread_id')==request['thread_id'] and e.get('turn_id')==request['turn_id'] and e.get('turn_status')=='interrupted' for e in events)
assert any(e.get('method')=='thread/start' and e.get('direction')=='response' and e.get('approvals_reviewer')=='user' for e in events)
print(json.dumps({'core_probe_passed':True,'native_revision_pinned_research_edit':True,'native_negative_approval_trace_verified':True,'offered':request['available_decisions'],'submitted':reply['decision'],'command_executed':False,'original_approval_probe_passed':approval['passed'],'original_approval_probe_stage':approval['stage'],'original_disposition_preserved':True,'limitation':'Original approval-only recorder waited for idle after interruption; whole-probe failure is retained. Unsupported MCP elicitation is not qualified.'},indent=2))
