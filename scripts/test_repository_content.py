import importlib.util
import os
from pathlib import Path
import subprocess
import tempfile
import unittest

SCRIPT = Path(__file__).with_name('check-repository-content.py').resolve()

class RepositoryContentTest(unittest.TestCase):
    def test_guard_checks_index_and_deleted_history(self):
        with tempfile.TemporaryDirectory() as directory:
            def git(*args):
                return subprocess.run(['git', *args], cwd=directory, check=True, capture_output=True)
            def guard(*args):
                return subprocess.run(['python3', str(SCRIPT), *args], cwd=directory, capture_output=True)
            git('init', '-q')
            git('config', 'user.name', 'Guard test')
            git('config', 'user.email', 'guard@example.invalid')
            path = Path(directory)
            (path/'README.md').write_text('Source documentation\n')
            git('add', 'README.md')
            self.assertEqual(guard('--staged').returncode, 0)
            git('commit', '-qm', 'source')
            self.assertEqual(guard('--history', 'HEAD').returncode, 0)
            (path/'artifacts').mkdir()
            (path/'artifacts'/'result.json').write_text('{}')
            git('add', 'artifacts')
            self.assertEqual(guard('--staged').returncode, 1)
            git('commit', '-qm', 'forbidden historical artifact')
            git('rm', '-r', 'artifacts')
            git('commit', '-qm', 'delete artifact')
            self.assertEqual(guard('--staged').returncode, 0)
            self.assertEqual(guard('--history', 'HEAD').returncode, 1)
            (path/'demo.MP4').write_bytes(b'video fixture')
            git('add', 'demo.MP4')
            self.assertEqual(guard('--staged').returncode, 1)

if __name__ == '__main__':
    unittest.main()
