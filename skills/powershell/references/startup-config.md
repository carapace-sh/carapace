# PowerShell Startup and Configuration

In-depth reference for PowerShell's startup sequence, profile files, `powershell.config.json`, execution policy, environment variables, and configuration mechanisms.

## Profile Files

PowerShell supports four profile files, scoped by user and host. They are executed in order at startup for interactive sessions.

### Profile Locations

| Profile | Windows | Linux / macOS |
|---------|---------|---------------|
| **AllUsers, AllHosts** | `$PSHOME\Profile.ps1` | `$PSHOME/profile.ps1` |
| **AllUsers, CurrentHost** | `$PSHOME\Microsoft.PowerShell_profile.ps1` | `$PSHOME/Microsoft.PowerShell_profile.ps1` |
| **CurrentUser, AllHosts** | `$HOME\Documents\PowerShell\Profile.ps1` | `~/.config/powershell/profile.ps1` |
| **CurrentUser, CurrentHost** | `$HOME\Documents\PowerShell\Microsoft.PowerShell_profile.ps1` | `~/.config/powershell/Microsoft.PowerShell_profile.ps1` |

### Host-Specific Profiles

For alternative hosts (e.g., VS Code), the host-specific profile replaces `Microsoft.PowerShell` with the host name:

| Profile | Windows | Linux / macOS |
|---------|---------|---------------|
| AllUsers, CurrentHost (VS Code) | `$PSHOME\Microsoft.VSCode_profile.ps1` | `$PSHOME/Microsoft.VSCode_profile.ps1` |
| CurrentUser, CurrentHost (VS Code) | `$HOME\Documents\PowerShell\Microsoft.VSCode_profile.ps1` | `~/.config/powershell/Microsoft.VSCode_profile.ps1` |

### Loading Order

Profiles execute in this order (later profiles can override earlier ones):

1. **AllUsersAllHosts** — runs first
2. **AllUsersCurrentHost**
3. **CurrentUserAllHosts**
4. **CurrentUserCurrentHost** — runs last (most commonly referred to as "your profile")

### The $PROFILE Variable

The `$PROFILE` automatic variable stores profile paths as note properties:

```powershell
$PROFILE                        # CurrentUserCurrentHost
$PROFILE.CurrentUserCurrentHost # Same as above
$PROFILE.CurrentUserAllHosts   # CurrentUserAllHosts
$PROFILE.AllUsersCurrentHost   # AllUsersCurrentHost
$PROFILE.AllUsersAllHosts      # AllUsersAllHosts

# View all paths
$PROFILE | Select-Object *
```

### Creating and Editing Profiles

```powershell
# Create if it doesn't exist
if (!(Test-Path -Path $PROFILE)) {
    New-Item -ItemType File -Path $PROFILE -Force
}

# Edit in Notepad
notepad $PROFILE

# Edit AllUsers profile (requires admin on Windows)
notepad $PROFILE.AllUsersAllHosts
```

### Skipping Profiles

```powershell
pwsh -NoProfile
```

### Profiles and Remote Sessions

Profiles **don't run automatically** in remote sessions, and `$PROFILE` isn't populated remotely:

```powershell
# Run local profile in remote session
Invoke-Command -Session $s -FilePath $PROFILE

# Dot-source remote profile
Invoke-Command -Session $s -ScriptBlock {
    . "$HOME\Documents\PowerShell\Microsoft.PowerShell_profile.ps1"
}
```

### MSIX Installations

For MSIX installations, `AllUsers` profiles in `$PSHOME` are read-only. Only `CurrentUser` profiles are supported.

## $PSHOME

`$PSHOME` is the directory containing the executing `System.Management.Automation.dll` assembly. This is the PowerShell installation directory.

| Platform | Typical Path |
|----------|-------------|
| Windows | `C:\Program Files\PowerShell\7\` |
| Linux | `/opt/microsoft/powershell/7/` |
| macOS | `/usr/local/microsoft/powershell/7/` |

## powershell.config.json

The `powershell.config.json` file contains configuration settings for PowerShell. It replaces Windows Registry settings for cross-platform support.

### Configuration Scopes

| Scope | Location | Precedence |
|-------|----------|------------|
| **AllUsers** | `$PSHOME/powershell.config.json` | Lower (overridden by CurrentUser) |
| **CurrentUser** | `Split-Path $PROFILE.CurrentUserCurrentHost` | Higher |

### Scope Precedence (Full)

1. Windows Group Policy (Windows only) — highest
2. AllUsers `powershell.config.json`
3. CurrentUser `powershell.config.json` — lowest (wins)

### Available Configuration Keys

**Cross-Platform:**

| Key | Description |
|-----|-------------|
| `ExperimentalFeatures` | Array of experimental feature names to enable |
| `PSModulePath` | Module search paths (supports `%HOME%`, `%XDG_CONFIG_HOME%` env vars) |
| `Microsoft.PowerShell:ExecutionPolicy` | Execution policy (AllUsers sets LocalMachine, CurrentUser sets CurrentUser) |
| `DisableImplicitWinCompat` | Disable Windows PowerShell Compatibility |
| `WindowsPowerShellCompatibilityModuleDenyList` | Modules excluded from compatibility mode |
| `WindowsPowerShellCompatibilityNoClobberModuleList` | Modules that shouldn't be clobbered |

**Logging (Linux/macOS):**

| Key | Description |
|-----|-------------|
| `LogChannels` | Log channels to enable |
| `LogIdentity` | Log identity name |
| `LogKeywords` | Log keywords filter |
| `LogLevel` | Log level |

**Policy Settings (Windows):**

| Key | Description |
|-----|-------------|
| `PowerShellPolicies.ExecutionPolicy` | Script execution policy |
| `PowerShellPolicies.ConsoleSessionConfiguration` | Session configuration endpoint |
| `PowerShellPolicies.ModuleLogging` | Module logging settings |
| `PowerShellPolicies.ProtectedEventLogging` | Encrypted event logging |
| `PowerShellPolicies.ScriptBlockLogging` | Script block logging |
| `PowerShellPolicies.ScriptExecution` | Script execution settings |
| `PowerShellPolicies.Transcription` | Input/output transcription |
| `PowerShellPolicies.UpdatableHelp` | Updatable help settings |

### Example Configuration

```json
{
    "ExperimentalFeatures": [
        "PSCommandNotFoundSuggestion",
        "PSSubsystemPluginModel"
    ],
    "PSModulePath": "%HOME%/Documents/PowerShell/Modules:/usr/local/share/PowerShell/Modules",
    "Microsoft.PowerShell:ExecutionPolicy": "RemoteSigned"
}
```

### PSModulePath Configuration

The `PSModulePath` key supports environment variable embedding with `%` delimiters:

```json
{
    "PSModulePath": "%HOME%/Documents/PowerShell/Modules"
}
```

- Works on Windows, macOS, and Linux
- Available variables: `%HOME%`, `%XDG_CONFIG_HOME%`
- PowerShell variables cannot be embedded
- Case-sensitive on Linux/macOS
- Directory separators must be valid for the platform

## Execution Policy

Execution policy controls whether PowerShell loads configuration files and runs scripts. It is a safety feature, not a security boundary — users can bypass it by typing script contents at the command line.

### Available Policies

| Policy | Description |
|--------|-------------|
| **AllSigned** | All scripts must be digitally signed by a trusted publisher |
| **Bypass** | Nothing blocked, no warnings — for scripts in larger applications |
| **RemoteSigned** | Default on Windows — local scripts run, downloaded scripts must be signed |
| **Restricted** | Permits commands, blocks all script files (including profiles) |
| **Undefined** | No policy set — effective policy becomes Restricted (Windows) or Unrestricted (non-Windows) |
| **Unrestricted** | Default on non-Windows — runs all scripts, warns for non-local files |

### Scope Precedence

Scopes from highest to lowest precedence:

1. **MachinePolicy** — Group Policy for all users
2. **UserPolicy** — Group Policy for current user
3. **Process** — Current session only (`$Env:PSExecutionPolicyPreference`)
4. **CurrentUser** — Stored in user config
5. **LocalMachine** — Default scope, affects all users

```powershell
# View all effective policies
Get-ExecutionPolicy -List

# Get effective policy
Get-ExecutionPolicy

# Change policy
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# Set for single session
pwsh -ExecutionPolicy Bypass
```

### Non-Windows Behavior

On non-Windows platforms, the default is **Unrestricted** and cannot be changed. The behavior matches **Bypass** because those platforms don't implement Windows Security Zones.

### Profiles and Execution Policy

If the execution policy is `Restricted`, profiles **will not run**. Use `Set-ExecutionPolicy` to allow profile execution before creating a profile.

## Experimental Features

Enable experimental features via `powershell.config.json` or at runtime:

```powershell
# Enable at runtime
Enable-ExperimentalFeature PSCommandNotFoundSuggestion

# List all experimental features
Get-ExperimentalFeature
```

### Notable Experimental Features

| Feature | Description |
|---------|-------------|
| `PSCommandNotFoundSuggestion` | Suggest similar commands on failure |
| `PSSubsystemPluginModel` | Subsystem plugin model for prediction |
| `PSNativeCommandArgumentPassing` | Changed native command argument passing (now standard in PS 7.3+) |
| `PSNativeWindowsTildeExpansion` | Tilde expansion for native commands on Windows (PS 7.5+) |
| `PSNativeCommandPreserveBytePipe` | Preserve byte-stream for native command pipe (PS 7.4+) |

## Environment Variables

### PowerShell-Specific Environment Variables

| Variable | Description |
|----------|-------------|
| `PSExecutionPolicyPreference` | Sets Process-scope execution policy |
| `PSModulePath` | Additional module search paths |
| `CARAPACE_LOG` | Enable carapace debug logging |
| `CARAPACE_SANDBOX` | JSON mock context (set by sandbox tests) |
| `CARAPACE_LENIENT` | Allow unknown flags in carapace |
| `CARAPACE_MATCH` | Set to `CASE_INSENSITIVE` for case-insensitive matching |
| `CARAPACE_NOSPACE` | Additional nospace suffixes |
| `CARAPACE_UNFILTERED` | Skip prefix filtering |
| `CARAPACE_TOOLTIP` | Enable tooltip style |
| `CARAPACE_DESCRIPTION_LENGTH` | Max description length (default 80) |
| `NO_COLOR` | Disables ANSI colors (sets `$PSStyle.OutputRendering = PlainText`) |
| `TERM` | Terminal type — `dumb` disables VT, `xterm-mono`/`xterm` sets PlainText rendering |

### Terminal-Related Environment Variables

| Variable | Effect on PowerShell |
|----------|---------------------|
| `TERM=dumb` | `$Host.UI.SupportsVirtualTerminal = $false` |
| `TERM=xterm-mono` | `$PSStyle.OutputRendering = PlainText` |
| `NO_COLOR` exists | `$PSStyle.OutputRendering = PlainText` |

## Startup Sequence

The full PowerShell startup sequence for an interactive session:

1. **Process initialization** — PowerShell host starts
2. **Configuration loading** — `powershell.config.json` read (AllUsers, then CurrentUser)
3. **Execution policy check** — Determines if scripts can run
4. **Profile execution** (in order):
   - AllUsersAllHosts
   - AllUsersCurrentHost
   - CurrentUserAllHosts
   - CurrentUserCurrentHost
5. **Module auto-loading** — Modules in `$PSModulePath` are available
6. **PSReadLine initialization** — Key bindings and options set
7. **Prompt displayed** — User can interact

### Non-Interactive Startup

For non-interactive sessions (`pwsh -Command`, `pwsh -File`):
- Profiles are **not loaded** by default
- Use `-NoProfile` explicitly to ensure no profiles load
- Execution policy still applies

### Startup Parameters

```powershell
pwsh -NoProfile                    # Skip profiles
pwsh -ExecutionPolicy Bypass       # Set execution policy
pwsh -Command "Get-Process"        # Run command
pwsh -File script.ps1              # Run script
pwsh -ConfigurationFile config.json # Use specific config
```

## Group Policy (Windows)

On Windows, Group Policy settings take precedence over `powershell.config.json`:

| Policy Area | Settings |
|-------------|----------|
| `ModuleLogging` | Enable/disable, module list |
| `ProtectedEventLogging` | Enable/disable, encryption certificate |
| `ScriptBlockLogging` | Enable/disable, invocation logging |
| `Transcription` | Enable/disable, output directory, invocation header |
| `ScriptExecution` | Execution policy |
| `UpdatableHelp` | Default source path |

## References

- [about_Profiles](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_profiles)
- [about_PowerShell_Config](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_powershell_config)
- [about_Execution_Policies](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_execution_policies)
- [about_Environment_Variables](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_environment_variables)
- [Experimental Features](https://learn.microsoft.com/en-us/powershell/scripting/learn/experimental-features)
- [PowerShell Startup](https://learn.microsoft.com/en-us/powershell/scripting/learn/shell/running-commands)

## Related Skills

- **powershell** → `references/completion.md` — Register-ArgumentCompleter (typically configured in profile)
- **powershell** → `references/editor.md` — PSReadLine configuration (typically in profile)
- **powershell** → `references/styling.md` — $PSStyle configuration (typically in profile)
- **powershell** → `references/language.md` — Execution policy, argument passing
