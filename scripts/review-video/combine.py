#!/usr/bin/env python3
"""Join actual browser footage and the annotated CLI replay, without cuts or speed changes."""
import hashlib
import json
from pathlib import Path
import subprocess
import sys

browser, cli, target = map(lambda x: Path(x).resolve(), sys.argv[1:4])
assert not target.exists(), 'Completed combined demo exists; choose a fresh target'
target.parent.mkdir(parents=True, exist_ok=True)

def run(*args):
    subprocess.run(args, check=True)

def digest(p):
    return hashlib.sha256(p.read_bytes()).hexdigest()

def probe(p):
    return json.loads(subprocess.check_output(['ffprobe', '-v', 'error', '-show_entries',
        'format=duration,size:stream=codec_name,width,height,nb_frames,r_frame_rate', '-of', 'json', str(p)]))

originals = [(p, digest(p), probe(p)) for p in [browser, cli]]
assert all(info['streams'][0]['r_frame_rate'] == '25/1' for _, _, info in originals)
assert originals[0][2]['streams'][0]['width'] == 1440
assert originals[0][2]['streams'][0]['height'] == 1000
assert originals[1][2]['streams'][0]['width'] == 1440
assert originals[1][2]['streams'][0]['height'] == 900
run('ffmpeg', '-hide_banner', '-loglevel', 'error', '-n', '-i', str(browser), '-i', str(cli),
    '-filter_complex', '[0:v]setsar=1,setpts=PTS-STARTPTS[v0];[1:v]pad=1440:1000:0:0:color=0x0a1020,setsar=1,setpts=PTS-STARTPTS[v1];[v0][v1]concat=n=2:v=1:a=0[v]',
    '-map', '[v]', '-c:v', 'libx264', '-preset', 'slow', '-crf', '24', '-pix_fmt', 'yuv420p',
    '-movflags', '+faststart', '-an', str(target))
actual = probe(target)
expected_duration = sum(float(info['format']['duration']) for _, _, info in originals)
expected_frames = sum(int(info['streams'][0]['nb_frames']) for _, _, info in originals)
assert abs(float(actual['format']['duration']) - expected_duration) < 0.001
assert int(actual['streams'][0]['nb_frames']) == expected_frames
assert all(digest(p) == sha for p, sha, _ in originals), 'Source changed during combination'
run('ffmpeg', '-v', 'error', '-i', str(target), '-f', 'null', '-')
run('ffmpeg', '-hide_banner', '-loglevel', 'error', '-n', '-i', str(target), '-vf',
    'fps=1/43,scale=360:-1,tile=4x3', '-frames:v', '1', str(target.parent / 'contact-sheet.jpg'))
for second, name in [(43, 'browser'), (100, 'cli'), (455, 'boundary')]:
    run('ffmpeg', '-hide_banner', '-loglevel', 'error', '-n', '-ss', str(second), '-i', str(target),
        '-frames:v', '1', str(target.parent / ('review-' + name + '.png')))
print(json.dumps({'output': str(target), 'sha256': digest(target), 'media': actual,
    'cli_offset_seconds': float(originals[0][2]['format']['duration']),
    'sources': [{'path': str(p), 'sha256': sha, 'media': info} for p, sha, info in originals]}, indent=2))
