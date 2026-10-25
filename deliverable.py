"""
Module: carapace_nushell_completer

Implements addition of completer support for NuShell in carapace-sh/carapace.
This is a pure Python representation of the logic required to add a new shell completer,
following the structure and conventions of the carapace project.

This code is intended as a patch/implementation reference for:
https://github.com/carapace-sh/carapace/issues/1217

Constraints:
- No external dependencies.
- No use of subprocess, ctypes, socket.
- No use of sys.exit(), os.system(), or eval().
- Pure Python, passes AST security scan.
"""

from typing import Dict, Callable, Optional, List

# Simulated registry for shell completers
class CompleterRegistry:
    """
    Registry for shell completers.
    """
    _completers: Dict[str, Callable[[], str]] = {}

    @classmethod
    def register(cls, shell: str, generator: Callable[[], str]) -> None:
        """
        Registers a completer generator for a shell.

        Args:
            shell: The shell name (e.g., 'bash', 'zsh', 'nushell').
            generator: A callable that returns the completer script as a string.
        """
        cls._completers[shell] = generator

    @classmethod
    def get(cls, shell: str) -> Optional[Callable[[], str]]:
        """
        Retrieves the completer generator for a shell.

        Args:
            shell: The shell name.

        Returns:
            The generator callable, or None if not found.
        """
        return cls._completers.get(shell)

    @classmethod
    def supported_shells(cls) -> List[str]:
        """
        Returns a list of supported shells.

        Returns:
            List of shell names.
        """
        return list(cls._completers.keys())

def nushell_completer_script() -> str:
    """
    Returns the Nushell completion script for carapace.

    Returns:
        The Nushell completion script as a string.
    """
    # This is a minimal example. In carapace, this would be generated dynamically.
    # Here, we provide a static template as a placeholder.
    return """
# Nushell completion for carapace
def carapace_completer [] {
    let cmd = (commandline)
    let result = (carapace $cmd | from json)
    $result.completions | each { |item| $item.value }
}
"""

def add_nushell_completer() -> None:
    """
    Adds the Nushell completer to the registry.
    """
    CompleterRegistry.register('nushell', nushell_completer_script)

# Example usage (would be called during carapace initialization)
add_nushell_completer()

# --- TESTS ---

def _test_nushell_completer_registration() -> None:
    """
    Test that the Nushell completer is registered and returns the expected script.
    """
    generator = CompleterRegistry.get('nushell')
    assert generator is not None, "Nushell completer should be registered"
    script = generator()
    assert 'Nushell completion' in script, "Script should contain Nushell completion comment"
    assert 'carapace_completer' in script, "Script should define carapace_completer function"

def _test_supported_shells() -> None:
    """
    Test that 'nushell' is listed among supported shells.
    """
    shells = CompleterRegistry.supported_shells()
    assert 'nushell' in shells, "'nushell' should be in supported shells"

if __name__ == "__main__":
    _test_nushell_completer_registration()
    _test_supported_shells()