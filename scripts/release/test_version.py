import unittest
from version import calculate


class VersionPolicyTests(unittest.TestCase):
    def test_initial_remains_point_one(self):
        self.assertEqual(calculate(["feat!: new format", "fix: bug"])["version"], "0.1.0")

    def test_aggregate_precedence(self):
        for messages in (["fix: bug", "feat: feature"], ["feat: feature", "perf: faster"]):
            self.assertEqual(calculate(messages, "0.1.9")["version"], "0.2.0")
        self.assertEqual(calculate(["perf: faster", "fix: bug"], "0.1.9")["version"], "0.1.10")

    def test_breaking_on_every_type(self):
        for kind in "feat fix docs refactor perf test build ci chore style revert".split():
            for message in [f"{kind}(scope)!: change", f"{kind}: change\n\nBREAKING CHANGE: new shape", f"{kind}: change\n\nBREAKING-CHANGE: new shape"]:
                self.assertEqual(calculate([message], "0.3.8")["version"], "0.4.0")

    def test_no_release(self):
        for messages in ([], ["docs: clarification", "chore: merge (#9)"]):
            self.assertFalse(calculate(messages, "0.2.0")["release_due"])
            with self.assertRaises(ValueError):
                calculate(messages, "0.2.0", "0.2.1")

    def test_reject_major_and_mismatches(self):
        for previous in ("1.0.0", "2.0.0", "0.01.0", "v0.1.0", "0.1.0-rc.1"):
            with self.assertRaises(ValueError):
                calculate(["feat: feature"], previous)
        for proposed in ("1.0.0", "2.1.0", "0.3.0", "0.01.0"):
            with self.assertRaises(ValueError):
                calculate(["feat: feature"], "0.1.0", proposed)
        self.assertEqual(calculate(["fix: bug"], "0.1.0", "0.1.1")["version"], "0.1.1")

    def test_malformed_messages(self):
        for message in ("feature: no", "feat: ", "feat(scope):", "Merge pull request", "", "feat\n: x"):
            with self.assertRaises(ValueError):
                calculate([message])
        for messages in ("feat: x", [None], {}):
            with self.assertRaises(ValueError):
                calculate(messages)


if __name__ == "__main__":
    unittest.main()
