import unittest
from unittest.mock import patch
from check_commits import check


class CommitGuardTests(unittest.TestCase):
    def test_exact_range_and_full_body(self):
        with patch("check_commits.subprocess.check_output", return_value="docs: change\n\nBREAKING CHANGE: shape\n\0\nfix: bug\n\0\n") as run:
            result = check("a"*40, "b"*40, "feat: new feature")
        self.assertEqual(result["commits_checked"], 2)
        self.assertEqual(run.call_args.args[0][-1], "a"*40+".."+"b"*40)

    def test_new_branch_checks_tip_without_rewriting_history(self):
        with patch("check_commits.subprocess.check_output", return_value="chore: initial\0") as run:
            check("0"*40, "b"*40)
        self.assertEqual(run.call_args.args[0][-2:], ["-1", "b"*40])

    def test_reject_bad_commit_or_title(self):
        for output,title in [("bad message\0", None), ("fix: fine\0", "Bad PR title")]:
            with patch("check_commits.subprocess.check_output", return_value=output):
                with self.assertRaises(ValueError):
                    check("a"*40, "b"*40, title)

    def test_reject_revision_option_injection(self):
        with patch("check_commits.subprocess.check_output") as run:
            with self.assertRaises(ValueError):
                check("--all", "b"*40)
        run.assert_not_called()
