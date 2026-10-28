#!/usr/bin/env python3
"""Join complete Playwright segments without cuts, speed changes or source overwrites."""
import hashlib,json,pathlib,subprocess,sys
paths=[pathlib.Path(p).resolve() for p in sys.argv[1:]];target=paths.pop();assert len(paths)>1 and not target.exists()
def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()
def probe(p):return json.loads(subprocess.check_output(['ffprobe','-v','error','-show_entries','format=duration,size:stream=width,height,nb_frames,r_frame_rate','-of','json',str(p)]))
sources=[{'path':str(p),'sha256':sha(p),'media':probe(p)} for p in paths]
for s in sources:
 v=s['media']['streams'][0];assert v['width']==1440 and v['height'] in [900,1000] and v['r_frame_rate']=='25/1'
filters=[]
for i,s in enumerate(sources):filters.append(f'[{i}:v]pad=1440:1000:0:0:color=0x0a1020,setsar=1,setpts=PTS-STARTPTS[v{i}]')
filters.append(''.join(f'[v{i}]' for i in range(len(paths)))+f'concat=n={len(paths)}:v=1:a=0[v]')
args=['ffmpeg','-v','error','-n']
for p in paths:args+=['-i',str(p)]
args+=['-filter_complex',';'.join(filters),'-map','[v]','-c:v','libx264','-preset','fast','-crf','24','-pix_fmt','yuv420p','-movflags','+faststart','-an',str(target)]
subprocess.run(args,check=True);actual=probe(target)
assert abs(float(actual['format']['duration'])-sum(float(s['media']['format']['duration']) for s in sources))<.001
assert int(actual['streams'][0]['nb_frames'])==sum(int(s['media']['streams'][0]['nb_frames']) for s in sources)
assert all(sha(p)==s['sha256'] for p,s in zip(paths,sources))
subprocess.run(['ffmpeg','-v','error','-i',str(target),'-f','null','-'],check=True)
subprocess.run(['ffmpeg','-v','error','-n','-i',str(target),'-vf','fps=1/60,scale=360:-1,tile=4x4','-frames:v','1',str(target.parent/'contact-sheet.jpg')],check=True)
report={'sources':sources,'output':str(target),'sha256':sha(target),'media':actual};(target.parent/'combination.json').write_text(json.dumps(report,indent=2)+'\n');print(json.dumps(report))
