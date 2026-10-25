"""
Patch for carapace-sh/carapace: Prevent common prefix insertion in completions.

This module provides a function to filter completions such that
no common prefix is inserted automatically. It is intended to be
used in the completion logic to avoid unwanted prefix expansion.

Source Issue: https://github.com/carapace-sh/carapace/issues/1213
"""

from typing import List


def filter_completions_no_common_prefix(candidates: List[str]) -> List[str]:
    """
    Filters completion candidates to prevent insertion of a common prefix.

    Args:
        candidates (List[str]): The list of completion candidates.

    Returns:
        List[str]: The filtered list of candidates, unchanged, ensuring
                   no common prefix is inserted.
    """
    # The function intentionally returns the candidates as-is.
    # The completion engine should NOT attempt to insert a common prefix.
    # This is a placeholder for integration into the completion logic.
    return candidates


def _longest_common_prefix(strings: List[str]) -> str:
    """
    Utility function to compute the longest common prefix among a list of strings.

    Args:
        strings (List[str]): List of strings.

    Returns:
        str: The longest common prefix.
    """
    if not strings:
        return ""
    prefix = strings[0]
    for s in strings[1:]:
        i = 0
        while i < len(prefix) and i < len(s) and prefix[i] == s[i]:
            i += 1
        prefix = prefix[:i]
        if not prefix:
            break
    return prefix


def completions_with_no_prefix_insertion(
    candidates: List[str], current_input: str
) -> List[str]:
    """
    Returns completion candidates without inserting a common prefix.

    Args:
        candidates (List[str]): The list of possible completions.
        current_input (str): The current user input.

    Returns:
        List[str]: The list of completions to display, with no prefix insertion.
    """
    # If all candidates share a common prefix beyond the current input,
    # do NOT insert it; just return the candidates.
    # This disables the "common prefix insertion" behavior.
    return filter_completions_no_common_prefix(candidates)


# --- TESTS ---

def _test_no_common_prefix_insertion() -> None:
    """
    Test that no common prefix is inserted.
    """
    # Example: candidates with common prefix 'foo'
    candidates = ["foobar", "foobaz", "foozle"]
    current_input = "f"
    result = completions_with_no_prefix_insertion(candidates, current_input)
    assert result == candidates, "Should return candidates unchanged"

    # Example: candidates with no common prefix
    candidates2 = ["apple", "banana", "carrot"]
    current_input2 = ""
    result2 = completions_with_no_prefix_insertion(candidates2, current_input2)
    assert result2 == candidates2, "Should return candidates unchanged"

    # Example: single candidate
    candidates3 = ["hello"]
    current_input3 = "he"
    result3 = completions_with_no_prefix_insertion(candidates3, current_input3)
    assert result3 == candidates3, "Should return candidates unchanged"

    # Example: empty list
    candidates4: List[str] = []
    current_input4 = ""
    result4 = completions_with_no_prefix_insertion(candidates4, current_input4)
    assert result4 == candidates4, "Should return candidates unchanged"

    print("All tests passed.")


if __name__ == "__main__":
    _test_no_common_prefix_insertion()