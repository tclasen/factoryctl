#!/usr/bin/env python3
"""Qualify a frozen factoryctl binary through an actual PTY and native Codex.

Inputs are automated synthetic test answers, never operator submissions. Run into
an absent output directory. No native request/config/log bodies are persisted.
The same executable acts as an identifier-only transparent runtime tap when the
probe supplies its child-only FACTORYCTL_QUAL_TAP environment switch.
"""
import argparse
import hashlib
import json
import os
import pathlib
import platform
import pty
import queue
import re
import shutil
import subprocess
import sys
import termios
import threading
import time


def digest(path):
    return hashlib.sha256(pathlib.Path(path).read_bytes()).hexdigest()


def manifest(repo):
    paths = [p for root in ('cmd', 'internal') for p in (repo / root).rglob('*') if p.is_file()]
    paths += [repo / 'go.mod', repo / 'go.sum']
    return {str(p.relative_to(repo)): digest(p) for p in sorted(paths)}


def write_json(path, value):
    path.write_text(json.dumps(value, indent=2, sort_keys=True) + '\n')


def tap():
    """Forward protocol bytes; persist only allowlisted correlation metadata."""
    log = pathlib.Path(os.environ['FACTORYCTL_QUAL_TRACE']).open('a')
    lock = threading.Lock()
    methods = {}
    approvals = set()
    def emit(value):
        with lock:
            log.write(json.dumps(value) + '\n')
            log.flush()
    child_env = dict(os.environ)
    child_env.pop('FACTORYCTL_QUAL_TAP', None)
    runtime_args = list(sys.argv[1:])
    manual_reviewer = os.environ.get('FACTORYCTL_QUAL_REVIEWER') == 'user'
    if manual_reviewer and runtime_args and runtime_args[0] == 'app-server':
        runtime_args += ['-c', 'approvals_reviewer="user"']
        emit({'direction': 'runtime_invocation', 'approvals_reviewer_override': 'user',
              'scope': 'this qualification child invocation only'})
    proc = subprocess.Popen([os.environ['FACTORYCTL_QUAL_CODEX'], *runtime_args],
                            stdin=subprocess.PIPE, stdout=subprocess.PIPE,
                            stderr=subprocess.PIPE, bufsize=0, env=child_env)
    def up():
        for line in sys.stdin.buffer:
            try:
                item = json.loads(line)
                method = item.get('method')
                if method:
                    methods[str(item.get('id'))] = method
                if method in ('thread/start', 'turn/start', 'turn/interrupt'):
                    params = item.get('params', {})
                    emit({'direction': 'request', 'method': method,
                          'request_id': item.get('id'),
                          'thread_id': params.get('threadId'),
                          'turn_id': params.get('turnId')})
                if not method and item.get('id') in approvals:
                    value = item.get('result', {}).get('decision')
                    emit({'direction': 'approval_reply', 'request_id': item.get('id'),
                          'decision': value if value in ('decline', 'cancel', 'accept') else 'unsupported'})
            except (ValueError, TypeError, AttributeError):
                pass
            try:
                proc.stdin.write(line)
                proc.stdin.flush()
            except BrokenPipeError:
                break
        try:
            proc.stdin.close()
        except BrokenPipeError:
            pass
    def stderr():
        # Preserve native diagnostics in the product pipe, never in the tap log.
        for line in proc.stderr:
            sys.stderr.buffer.write(line)
            sys.stderr.buffer.flush()
    threading.Thread(target=up, daemon=True).start()
    threading.Thread(target=stderr, daemon=True).start()
    for line in proc.stdout:
        try:
            item = json.loads(line)
            method = item.get('method')
            request_method = methods.get(str(item.get('id')))
            if method and 'id' in item:
                emit({'direction':'native_request_method','method':method,'request_id':item['id']})
            if not method and request_method in ('thread/start', 'turn/start'):
                result = item.get('result', {})
                emit({'direction': 'response', 'method': request_method,
                      'request_id': item.get('id'),
                      'thread_id': result.get('thread', {}).get('id'),
                      'turn_id': result.get('turn', {}).get('id'),
                      'error_code': (item.get('error') or {}).get('code'),
                      'approvals_reviewer': result.get('approvalsReviewer') if result.get('approvalsReviewer') in ('user','auto_review','guardian_subagent') else None})
            if method in ('turn/completed', 'item/tool/requestUserInput',
                          'item/commandExecution/requestApproval',
                          'item/fileChange/requestApproval'):
                params = item.get('params', {})
                safe_decisions = [value for value in params.get('availableDecisions', [])
                                  if isinstance(value, str) and value in ('accept', 'decline', 'cancel')]
                if method.endswith('/requestApproval'):
                    approvals.add(item.get('id'))
                emit({'direction': 'event', 'method': method,
                      'request_id': item.get('id'),
                      'thread_id': params.get('threadId'),
                      'turn_id': params.get('turnId') or params.get('turn', {}).get('id'),
                      'turn_status': params.get('turn', {}).get('status'),
                      'available_decisions': safe_decisions,
                      'available_decisions_omitted': 'availableDecisions' not in params})
            call = item.get('params', {}).get('item', {})
            if method in ('item/started', 'item/completed') and call.get('type') == 'mcpToolCall':
                name = call.get('tool', '')
                if name in ('factory_node_read', 'factory_node_update'):
                    arguments = call.get('arguments', {})
                    if isinstance(arguments, str):
                        arguments = json.loads(arguments)
                    if not isinstance(arguments, dict):
                        arguments = {}
                    emit({'direction': 'mcp_call', 'method': method, 'tool': name,
                          'item_id': call.get('id'), 'status': call.get('status'),
                          'thread_id': item['params'].get('threadId'),
                          'turn_id': item['params'].get('turnId'),
                          'fixture_path': arguments.get('ref', {}).get('path') == 'tasks/native-edit.md',
                          'expected_revision': arguments.get('expected_revision') if re.fullmatch(r'[a-f0-9]{64}', str(arguments.get('expected_revision', ''))) else None,
                          'expected_definition': arguments.get('expected_definition') if re.fullmatch(r'[a-f0-9]{64}', str(arguments.get('expected_definition', ''))) else None,
                          'sets_expected_title': arguments.get('set', {}).get('title') == 'NATIVE_MCP_EDIT_VERIFIED'})
            if method == 'item/completed' and item.get('params', {}).get('item', {}).get('type') == 'commandExecution':
                command = item['params']['item']
                emit({'direction': 'command_completion', 'method': method,
                      'thread_id': item['params'].get('threadId'),
                      'turn_id': item['params'].get('turnId'),
                      'status': command.get('status'), 'exit_code': command.get('exitCode')})
        except (ValueError, TypeError, AttributeError):
            pass
        sys.stdout.buffer.write(line)
        sys.stdout.buffer.flush()
    return proc.wait()


def sanitized(text):
    text = re.sub(r'(?s)(Native runtime request[^\n]*\n).*?(\nRespond with JSON[^\n]*)', r'\1[native request body omitted]\2', text)
    text = re.sub(r'(?i)(bearer\s+)[^\s"\']+', r'\1[REDACTED]', text)
    text = re.sub(r'\b(?:sk-|sk_live_|ghp_|github_pat_)[A-Za-z0-9_-]{12,}', '[REDACTED]', text)
    return text


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--binary', type=pathlib.Path, required=True)
    parser.add_argument('--expected-binary-sha256', required=True)
    parser.add_argument('--repo', type=pathlib.Path, required=True)
    parser.add_argument('--output', type=pathlib.Path, required=True)
    parser.add_argument('--document-edit', action='store_true', help='require a native revision-checked MCP Task edit')
    parser.add_argument('--native-approval-only', action='store_true', help='qualify only the native command approval route without MCP tools')
    parser.add_argument('--manual-approval-reviewer', action='store_true', help='select approvals_reviewer=user only for the qualification runtime invocation')
    parser.add_argument('--research-fixture', action='store_true', help='use a file-defined Research Fact with explicit synthetic provenance')
    parser.add_argument('--native-approval', action='store_true', help='attempt an innocuous escalation and decline the actual offered native request')
    parser.add_argument('--codex', default=shutil.which('codex'))
    parser.add_argument('--timeout', type=int, default=360)
    opts = parser.parse_args()
    if opts.native_approval_only and (not opts.native_approval or not opts.manual_approval_reviewer or opts.document_edit):
        parser.error('--native-approval-only requires --native-approval --manual-approval-reviewer and no document edit')
    if not opts.codex:
        parser.error('installed Codex executable not found')
    if not 30 <= opts.timeout <= 540:
        parser.error('--timeout must be between 30 and 540 seconds')
    binary, repo, output = opts.binary.resolve(), opts.repo.resolve(), opts.output.resolve()
    if output.exists():
        parser.error('choose an absent versioned output directory; failed attempts are preserved')
    before_binary = digest(binary)
    if before_binary != opts.expected_binary_sha256:
        parser.error('binary differs from the Coordinator BUILD READY digest')
    output.mkdir(parents=True)
    before_sources = manifest(repo)
    write_json(output / 'source-before.json', before_sources)
    runtime_version = subprocess.run([opts.codex, '--version'], text=True, capture_output=True, timeout=20)
    write_json(output / 'environment.json', {
        'host': platform.platform(), 'machine': platform.machine(),
        'python': platform.python_version(), 'runtime_executable': str(pathlib.Path(opts.codex).resolve()),
        'runtime_version_exit': runtime_version.returncode,
        'runtime_version': sanitized(runtime_version.stdout.strip()),
        'binary': str(binary), 'binary_sha256': before_binary,
        'automated_synthetic_test_input': True,
        'operator_input_or_acceptance_claimed': False,
        'cli_disposition': 'retain; no CLI retirement decision is claimed',
        'external_mcp_apps': 'not qualified by this terminal probe',
        'access_enforcement': 'advisory; prompt tool restrictions do not enforce filesystem isolation',
        'invocation_approvals_reviewer_override': 'user' if opts.manual_approval_reviewer else None,
    })
    scratch = output / 'scratch-project'
    (scratch / '.git/objects').mkdir(parents=True)
    (scratch / '.git/refs/heads').mkdir(parents=True)
    (scratch / '.git/HEAD').write_text('ref: refs/heads/main\n')
    init = subprocess.run([str(binary), '--path', str(scratch), 'init', '--project-id',
                           'runtime-qualification', '--starter', 'software'],
                          text=True, capture_output=True, timeout=30)
    write_json(output / 'init.json', {'exit': init.returncode, 'stdout': sanitized(init.stdout),
                                     'stderr': sanitized(init.stderr)})
    if init.returncode:
        return 1
    fixture = scratch / 'factory/tasks/native-edit.md'
    fixture_before = None
    fixture_before_projection = None
    if opts.document_edit:
        fixture_type = 'Research Fact' if opts.research_fixture else 'Task'
        extra_inputs = []
        summary = 'Native MCP revision checked edit qualification'
        if opts.research_fixture:
            (scratch / 'factory/definitions/research-fact.md').write_text('---\ntype: Record Definition\ntitle: Research Fact definition\ndefines: Research Fact\ndirectory: tasks\nfields:\n  citation: {kind: string, required: true}\n  certainty: {kind: enum, required: true, values: [observed, assumption, inference]}\n---\n# Research Fact\nSources are attributed claims; agents interpret their meaning.\n')
            extra_inputs = ['--input', 'citation=Synthetic qualification source; not a scientific claim',
                            '--input', 'certainty=observed']
            summary = 'Observed: this is a synthetic qualification fixture.\n\nAssumption: the citation identifies only test data.\n\nLLM inference: a successful edit demonstrates transport behavior, not research truth.'
        created = subprocess.run([str(binary), '--path', str(scratch), '--json', 'records', 'create',
                                  '--type', fixture_type, '--slug', 'native-edit', '--input', 'title=Native MCP fixture',
                                  '--input', 'summary=' + summary, *extra_inputs],
                                 text=True, capture_output=True, timeout=30)
        write_json(output / 'fixture-create.json', {'exit': created.returncode,
                   'stdout': sanitized(created.stdout), 'stderr': sanitized(created.stderr)})
        if created.returncode:
            return 1
        initial_show = subprocess.run([str(binary), '--path', str(scratch), '--json', 'records', 'show',
                                       'tasks/native-edit.md'], text=True, capture_output=True, timeout=30)
        if initial_show.returncode:
            write_json(output / 'fixture-before-show.json', {'exit': initial_show.returncode})
            return 1
        fixture_before_projection = json.loads(initial_show.stdout)
        write_json(output / 'fixture-before-show.json', fixture_before_projection)
        fixture_before = fixture.read_bytes()
        (output / 'fixture-before.md').write_bytes(fixture_before)
    prefs = output / 'runtime-preferences.md'
    prefs.write_text('---\nversion: 1\nruntime: codex\n---\nIsolated test invocation preferences.\n')
    settings = output / 'chat-settings.json'
    settings.write_text('{"version":1,"projects":{}}\n')
    trace = output / 'native-identities.jsonl'
    env = dict(os.environ)
    env.update(FACTORYCTL_QUAL_TAP='1', FACTORYCTL_QUAL_CODEX=str(pathlib.Path(opts.codex).resolve()),
               FACTORYCTL_QUAL_TRACE=str(trace))
    if opts.manual_approval_reviewer:
        env['FACTORYCTL_QUAL_REVIEWER'] = 'user'
    prompt = '''This is an authorized isolated automated interaction qualification using synthetic test answers. Use only factoryctl factory_user_ask and factory_user_wait tools; do not use shell, file, delegation or other tools. First ask a text question with exact prompt "CLI qualification: what is your project nickname?" and wait for the answer. Then ask a single_choice question with exact prompt "CLI qualification: choose a review mode or write a custom answer" and choices brief/Brief, detailed/Detailed, and wait. Return both exact answers on one line prefixed CLI_ROUNDTRIP. Do not invent answers. Keep this conversation for a follow-up.'''
    if opts.document_edit:
        prompt = ('First use factory_node_read for ref.path tasks/native-edit.md. Then use factory_node_update on that exact record with both returned expected_revision and expected_definition, changing only title to exactly NATIVE_MCP_EDIT_VERIFIED using set. Use no direct file or shell tools. Reread it through factory_node_read and confirm the observed persisted title with prefix NATIVE_EDIT_RESULT. After the edit, continue the interaction qualification below. ' + prompt.replace('Use only factoryctl factory_user_ask and factory_user_wait tools;', 'Use only factoryctl factory_node_read, factory_node_update, factory_user_ask and factory_user_wait tools;'))
    if opts.native_approval_only:
        prompt = 'Bounded native approval qualification: request exactly one innocuous /usr/bin/true command via exec_command with cmd /usr/bin/true and sandbox_permissions require_escalated, with justification asking permission for this qualification. No file changes, network, MCP tools, other commands, or alternate tools. The synthetic terminal will decline any offered native approval. If escalation is unavailable or declined, report that accurately and stop this turn; never retry, fall back, or execute without approval.'
    argv = [str(binary), '--path', str(scratch), 'chat', '--runtime', 'codex',
            '--runtime-executable', str(pathlib.Path(__file__).resolve()),
            '--preferences', str(prefs), '--chat-settings', str(settings), '--prompt', prompt]
    master, slave = pty.openpty()
    attrs = termios.tcgetattr(slave)
    attrs[3] &= ~termios.ECHO
    termios.tcsetattr(slave, termios.TCSANOW, attrs)
    is_terminal = os.isatty(slave)
    proc = subprocess.Popen(argv, stdin=slave, stdout=slave, stderr=slave,
                            start_new_session=True, env=env)
    os.close(slave)
    messages = queue.Queue()
    def reader():
        try:
            while True:
                data = os.read(master, 8192)
                if not data:
                    break
                messages.put(data.decode('utf8', errors='replace'))
        except OSError:
            pass
    threading.Thread(target=reader, daemon=True).start()
    transcript, actions, stage = '', [], 0
    started = time.monotonic()
    def send(value):
        os.write(master, (value + '\n').encode())
        actions.append({'seconds': round(time.monotonic() - started, 3),
                        'synthetic_terminal_input': value})
    def identities():
        if not trace.exists():
            return []
        rows = []
        for line in trace.read_text().splitlines():
            try:
                rows.append(json.loads(line))
            except ValueError:
                pass  # Last writer line may still be in progress.
        return rows
    failure = None
    approval_disposition = 'awaiting_native_surface' if opts.native_approval_only else 'not_requested'
    if opts.native_approval_only:
        stage = 30
    approval_reply_id = None
    approval_reply_decision = None
    try:
        while time.monotonic() - started < opts.timeout:
            if proc.poll() is not None and messages.empty():
                break
            try:
                transcript += messages.get(timeout=.2)
            except queue.Empty:
                pass
            if stage == 0 and 'Enter your answer' in transcript:
                send('juniper-qualification'); stage = 1
            elif stage == 1 and transcript.count('Enter your answer') >= 2:
                send('custom amber review'); stage = 2
            elif stage == 2 and 'Your follow-up:' in transcript:
                send('Without tools, repeat the two exact answers you received in this same conversation on one line prefixed FOLLOWUP_MEMORY.'); stage = 3
            elif stage == 3 and transcript.count('Your follow-up:') >= 2 and opts.native_approval:
                send('Bounded native approval qualification: request exactly one innocuous /usr/bin/true command via exec_command with cmd /usr/bin/true and sandbox_permissions require_escalated, with justification asking permission for this qualification. No file changes, network, other commands, or alternate tools. The synthetic terminal will decline any offered native approval. If escalation is unavailable or declined, report that accurately and stop this turn; never retry, fall back, or execute without approval.')
                approval_disposition = 'awaiting_native_surface'
                stage = 30
            elif stage == 30:
                offered = [row for row in identities() if row.get('method') == 'item/commandExecution/requestApproval']
                if offered and 'Respond with JSON matching this request' in transcript:
                    available = offered[-1].get('available_decisions', [])
                    decision = ('decline' if 'decline' in available or offered[-1].get('available_decisions_omitted') else
                                'cancel' if 'cancel' in available else None)
                    if decision:
                        approval_reply_id = offered[-1]['request_id']
                        approval_reply_decision = decision
                        send(json.dumps({'decision': decision}))
                        approval_disposition = decision + '_submitted_to_real_native_request'
                        stage = 31
                elif transcript.count('Your follow-up:') >= (1 if opts.native_approval_only else 3):
                    approval_disposition = 'unavailable_no_native_approval_surface_observed'
                    stage = 31
            elif stage == 31 and opts.native_approval_only and ('Your follow-up:' in transcript or any(row.get('method') == 'turn/completed' and row.get('turn_status') in ('interrupted','completed') for row in identities())):
                send('/quit'); stage = 6
            elif stage == 31 and transcript.count('Your follow-up:') >= 3:
                send('Cancellation qualification: ask one text question using factory_user_ask with prompt CLI cancellation test, then wait with factory_user_wait. Do not invent an answer or use other tools.')
                stage = 4
            elif stage == 3 and transcript.count('Your follow-up:') >= 2:
                send('Cancellation qualification: ask one text question using factory_user_ask with prompt CLI cancellation test, then wait with factory_user_wait. Do not invent an answer or use other tools.'); stage = 4
            elif stage == 4 and transcript.count('Enter your answer') >= 3:
                send('/cancel'); stage = 5
            elif stage == 5 and any(row.get('turn_status') == 'interrupted' for row in identities()):
                send('/quit'); stage = 6
        if proc.poll() is None:
            if stage != 6:
                failure = 'deadline reached before required interaction stages completed'
            send('/quit')
            try:
                proc.wait(timeout=8)
            except subprocess.TimeoutExpired:
                proc.terminate()
                proc.wait(timeout=5)
    except Exception as exc:
        failure = type(exc).__name__ + ': ' + str(exc)
        if proc.poll() is None:
            proc.terminate()
            proc.wait(timeout=5)
    finally:
        while not messages.empty():
            transcript += messages.get_nowait()
        os.close(master)
        (output / 'stdout.txt').write_text(sanitized(transcript))
        ids = identities()
        turns = [row for row in ids if row.get('direction') == 'request' and row.get('method') == 'turn/start']
        document_edit = {'requested': opts.document_edit, 'fixture_type': 'Research Fact' if opts.research_fixture else 'Task', 'synthetic_provenance': opts.research_fixture}
        if opts.document_edit:
            shown = subprocess.run([str(binary), '--path', str(scratch), '--json', 'records', 'show',
                                    'tasks/native-edit.md'], text=True, capture_output=True, timeout=30)
            write_json(output / 'fixture-show.json', {'exit': shown.returncode,
                       'stdout': sanitized(shown.stdout), 'stderr': sanitized(shown.stderr)})
            fixture_after = fixture.read_bytes()
            (output / 'fixture-after.md').write_bytes(fixture_after)
            def erase_title(raw):
                return re.sub(rb'(?m)^title:[^\n]*$', b'title: [TITLE]', raw, count=1)
            try:
                projection = json.loads(shown.stdout)
                node = projection.get('node', projection)
                verified_title = node.get('title') == 'NATIVE_MCP_EDIT_VERIFIED'
            except (ValueError, AttributeError):
                verified_title = False
            native_updates = [row for row in ids if row.get('direction') == 'mcp_call' and
                              row.get('tool') == 'factory_node_update' and row.get('fixture_path') and
                              row.get('sets_expected_title') and row.get('method') == 'item/completed' and
                              row.get('status') == 'completed' and
                              row.get('expected_revision') == hashlib.sha256(fixture_before).hexdigest() and
                              row.get('expected_definition') == fixture_before_projection.get('definition_hash')]
            document_edit['native_revision_pinned_mcp_update_observed'] = bool(native_updates)
            document_edit.update(cli_show_exit=shown.returncode,
                                 title_verified_by_fresh_cli=verified_title,
                                 only_title_bytes_changed=erase_title(fixture_before) == erase_title(fixture_after),
                                 actual_bytes_changed=fixture_before != fixture_after,
                                 before_sha256=hashlib.sha256(fixture_before).hexdigest(),
                                 after_sha256=hashlib.sha256(fixture_after).hexdigest())
            document_edit['passed'] = (shown.returncode == 0 and verified_title and
                document_edit['only_title_bytes_changed'] and document_edit['actual_bytes_changed'] and
                'NATIVE_EDIT_RESULT' in transcript and bool(native_updates))
        approval_replies = [row for row in ids if row.get('direction') == 'approval_reply' and
                            row.get('request_id') == approval_reply_id and row.get('decision') == approval_reply_decision and approval_reply_decision in ('decline','cancel')]
        native_approval = {'requested': opts.native_approval, 'disposition': approval_disposition,
                           'request_id': approval_reply_id, 'reply_decision': approval_reply_decision,
                           'negative_reply_observed': bool(approval_replies),
                           'decline_reply_observed': bool(approval_replies) and approval_reply_decision == 'decline',
                           'cancel_reply_observed': bool(approval_replies) and approval_reply_decision == 'cancel',
                           'manual_reviewer_invocation_override': opts.manual_approval_reviewer,
                           'effective_thread_reviewer': next((row.get('approvals_reviewer') for row in ids if row.get('method')=='thread/start' and row.get('direction')=='response'),None),
                           'command_completion_events': [row for row in ids if row.get('direction') == 'command_completion']}
        native_approval['qualified'] = bool(approval_replies) and not any(
            row.get('exit_code') == 0 for row in native_approval['command_completion_events'])
        after_sources = manifest(repo)
        write_json(output / 'source-after.json', after_sources)
        # Exact emitted answer checks are scoped to the assistant marker response;
        # merely seeing test input text in a displayed question is insufficient.
        def has_answers(marker):
            if marker not in transcript:
                return False
            response = transcript.split(marker, 1)[1].split('Your follow-up:', 1)[0]
            return 'juniper-qualification' in response and 'custom amber review' in response
        result = {'exit': proc.returncode, 'stage': stage, 'failure': failure,
                  'duration_seconds': round(time.monotonic() - started, 3),
                  'automated_real_terminal_input': True, 'operator_input_claimed': False,
                  'qualification_mode': 'native_approval_only' if opts.native_approval_only else 'mcp_document_and_pty' if opts.document_edit else 'pty',
                  'pty_isatty': is_terminal, 'terminal_echo_disabled': True,
                  'actions': actions, 'document_edit': document_edit, 'native_approval': native_approval,
                  'text_and_custom_roundtrip': has_answers('CLI_ROUNDTRIP'),
                  'followup_memory': has_answers('FOLLOWUP_MEMORY'),
                  'native_turn_count': len(turns),
                  'same_native_thread': len(turns) >= 3 and len({row.get('thread_id') for row in turns}) == 1,
                  'native_interrupt_requested': any(row.get('method') == 'turn/interrupt' for row in ids),
                  'native_interrupted': any(row.get('turn_status') == 'interrupted' for row in ids),
                  'binary_sha256': before_binary, 'binary_unchanged': digest(binary) == before_binary,
                  'source_manifest_sha256': digest(output / 'source-before.json'),
                  'sources_unchanged': before_sources == after_sources,
                  'native_ids': ids, 'cli_disposition': 'retain'}
        checks = ('pty_isatty', 'text_and_custom_roundtrip', 'followup_memory',
                  'same_native_thread', 'native_interrupt_requested', 'native_interrupted',
                  'binary_unchanged', 'sources_unchanged')
        result['passed'] = proc.returncode == 0 and stage == 6 and failure is None and all(result[k] for k in checks) and (not opts.document_edit or document_edit['passed'])
        if opts.native_approval_only:
            result['passed'] = proc.returncode == 0 and stage == 6 and failure is None and native_approval['qualified'] and result['binary_unchanged'] and result['sources_unchanged']
        write_json(output / 'result.json', result)
        print(json.dumps(result, indent=2))
    return 0 if result['passed'] else 1


if __name__ == '__main__':
    sys.exit(tap() if os.environ.get('FACTORYCTL_QUAL_TAP') == '1' else main())
