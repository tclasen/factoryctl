#!/usr/bin/env python3
"""Own an isolated current-binary server and record actual web editor journeys."""
import hashlib,json,os,pathlib,re,select,subprocess,sys,tempfile,time
root=pathlib.Path(__file__).resolve().parents[2]
out=pathlib.Path(os.environ['REVIEW_OUTPUT_DIR']).resolve()
if out.exists():raise SystemExit('Use a new output directory; preserve earlier attempts')
out.mkdir(parents=True)
binary=pathlib.Path(os.environ.get('FACTORYCTL_BINARY','/private/tmp/factoryctl-completion'))
fixture=pathlib.Path(tempfile.mkdtemp(prefix='factoryctl-completion-web-'))
subprocess.run(['git','init','-q',str(fixture)],check=True)
subprocess.run([str(binary),'--path',str(fixture),'init','--project-id','editor','--starter','software'],check=True,stdout=subprocess.DEVNULL)
server=subprocess.Popen([str(binary),'--path',str(fixture),'web','--listen','127.0.0.1:0'],stdout=subprocess.PIPE,stderr=subprocess.PIPE)
log=b'';started=time.monotonic();report={'binary_sha256':hashlib.sha256(binary.read_bytes()).hexdigest(),'fixture':str(fixture),'revision':subprocess.check_output(['git','rev-parse','HEAD'],cwd=root,text=True).strip()}
try:
 while time.monotonic()-started<15:
  if select.select([server.stdout],[],[],.2)[0]:log+=os.read(server.stdout.fileno(),4096)
  match=re.search(rb'http://127\.0\.0\.1:\d+',log)
  if match:break
  if server.poll() is not None:raise RuntimeError('Owned server exited before readiness')
 else:raise RuntimeError('Owned server readiness timed out')
 env=dict(os.environ,FACTORYCTL_WEB_URL=match[0].decode(),FACTORYCTL_WEB_ROOT=str(fixture),REVIEW_SOURCE_REVISION=report['revision'])
 command=['node',str(root/'scripts/review-video/completion-web-check.cjs')]
 report['command']=command
 result=subprocess.run(command,cwd=root,env=env,capture_output=True,timeout=600)
 (out/'runner.stdout').write_bytes(result.stdout);(out/'runner.stderr').write_bytes(result.stderr)
 report['exit_code']=result.returncode
 if result.returncode:raise RuntimeError(result.stderr.decode())
 child=next(out.glob('web-*/results.json'));data=json.loads(child.read_text());assert data['passed']
 video=pathlib.Path(data['video']);archive=child.parent/'demo.mp4'
 subprocess.run(['ffmpeg','-v','error','-n','-i',str(video),'-c:v','libx264','-preset','fast','-crf','24','-pix_fmt','yuv420p','-movflags','+faststart','-an',str(archive)],check=True)
 report.update(passed=True,results=str(child),video=str(archive),video_sha256=hashlib.sha256(archive.read_bytes()).hexdigest())
except Exception as error:
 report.update(passed=False,error=str(error));raise
finally:
 server.terminate()
 try:server.wait(timeout=5)
 except subprocess.TimeoutExpired:server.kill();server.wait()
 (out/'server.stdout').write_bytes(log);(out/'server.stderr').write_bytes(server.stderr.read())
 report['owned_server_stopped']=server.poll() is not None
 (out/'driver.json').write_text(json.dumps(report,indent=2)+'\n')
 print(json.dumps(report))
