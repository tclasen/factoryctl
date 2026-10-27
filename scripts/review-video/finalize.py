#!/usr/bin/env python3
"""Convert complete Playwright capture; extract inspection frames and checksums."""
import hashlib
import json
from pathlib import Path
import subprocess
import sys

out = Path(sys.argv[1]).resolve()

def run(*args):
    subprocess.run(args, check=True)

run('ffmpeg', '-hide_banner', '-loglevel', 'error', '-n', '-i', str(out / 'cli-review.webm'),
    '-c:v', 'libx264', '-preset', 'slow', '-crf', '28', '-pix_fmt', 'yuv420p',
    '-movflags', '+faststart', '-an', str(out / 'cli-review.mp4'))
probe = subprocess.check_output(['ffprobe', '-v', 'error', '-show_entries',
    'format=duration,size:stream=codec_name,width,height', '-of', 'json',
    str(out / 'cli-review.mp4')], text=True)
(out / 'media.json').write_text(probe)
run('ffmpeg', '-hide_banner', '-loglevel', 'error', '-n', '-ss', '38', '-i',
    str(out / 'cli-review.mp4'), '-frames:v', '1', str(out / 'poster.png'))
run('ffmpeg', '-hide_banner', '-loglevel', 'error', '-n', '-i', str(out / 'cli-review.mp4'),
    '-vf', 'fps=1/36,scale=480:-1,tile=4x3', '-frames:v', '1', str(out / 'contact-sheet.jpg'))
coverage = json.loads((out / 'coverage.json').read_text())
lines = ['# Recorded CLI review chapters', '', 'Times are relative to the start of the continuous walkthrough.', '',
         '| Start | End | Chapter |', '| --- | --- | --- |']
for scene in coverage['scenes']:
    def time(seconds):
        return f'{int(seconds)//60:02}:{int(seconds)%60:02}'
    title = scene['title'].split(' / ', 1)[-1]
    lines.append(f"| {time(scene['start_seconds'])} | {time(scene['end_seconds'])} | {scene['id']:02} / {title} |")
(out / 'chapters.md').write_text('\n'.join(lines) + '\n')
files = ['cli-review.mp4', 'cli-review.webm', 'results.json', 'coverage.json', 'toolchain.json',
         'media.json', 'poster.png', 'contact-sheet.jpg', 'chapters.md']
(out / 'SHA256SUMS').write_text(''.join(
    hashlib.sha256((out / name).read_bytes()).hexdigest() + '  ' + name + '\n'
    for name in files))
print(probe)
