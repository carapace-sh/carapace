# cobra-zsh-gen

[![Build Status](https://travis-ci.org/rsteube/cobra-zsh-gen.svg?branch=master)](https://travis-ci.org/rsteube/cobra-zsh-gen)

This is essentially the content of spf13/cobra#646 wich improves the generation of [zsh-completion](https://github.com/zsh-users/zsh-completions/blob/master/zsh-completions-howto.org) scripts. This provides temporary access to that changes until the PR is merged.

## Usage

A wrapper struct named `ZshCommand` and the util function `Wrap` were added to enable execution from a different package while keeping the code close to the PR. So instead of calling the additional zsh related functions directly on `Command` it has to be wrapped first:

```go
import (
    zsh "github.com/rsteube/cobra-zsh-gen"
)

zsh.Wrap(issueListCmd).MarkZshCompPositionalArgumentCustom(1, "__lab_completion_remote")
zsh.Wrap(ciCreateCmd).MarkZshCompPositionalArgumentCustom(1, "__lab_completion_remote_branches origin")
zsh.Wrap(RootCmd).GenZshCompletion(os.Stdout)
```

Apart from that refer to the original [instructions](zsh_completions.md).
