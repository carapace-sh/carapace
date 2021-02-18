FROM cimg/go:1.15 as base
USER root

FROM base as bat
RUN curl -L https://github.com/sharkdp/bat/releases/download/v0.17.1/bat-v0.17.1-x86_64-unknown-linux-gnu.tar.gz \
  | tar -C /usr/local/bin/ --strip-components=1  -xvz bat-v0.17.1-x86_64-unknown-linux-gnu/bat \
 && chmod +x /usr/local/bin/bat

FROM base as elvish
RUN curl https://dl.elv.sh/linux-amd64/elvish-HEAD.tar.gz | tar -xvz \
 && mv elvish-* /usr/local/bin/elvish

FROM base as goreleaser
RUN curl -L https://github.com/goreleaser/goreleaser/releases/download/v0.155.2/goreleaser_Linux_x86_64.tar.gz | tar -xvz goreleaser \
 && mv goreleaser /usr/local/bin/goreleaser

FROM rust as ion
RUN git clone https://gitlab.redox-os.org/redox-os/ion/ \
 && cd ion \
 && RUSTUP=0 make # By default RUSTUP equals 1, which is for developmental purposes \
 && sudo make install prefix=/usr \
 && sudo make update-shells prefix=/usr

FROM base as nu
RUN curl -L https://github.com/nushell/nushell/releases/download/0.24.1/nu_0_24_1_linux.tar.gz | tar -xvz \
 && mv nu_0_24_1_linux/nushell-0.24.1/nu* /usr/local/bin

FROM base as oil
RUN curl https://www.oilshell.org/download/oil-0.8.5.tar.gz | tar -xvz \
 && cd oil-*/ \
 && ./configure \
 && make \
 && ./install

FROM base as shellcheck
RUN  scversion="stable" \
 && wget -qO- "https://github.com/koalaman/shellcheck/releases/download/${scversion?}/shellcheck-${scversion?}.linux.x86_64.tar.xz" | tar -xJv shellcheck-stable/shellcheck \
 && mv shellcheck-stable/shellcheck /usr/local/bin/ \
 && chmod +x /usr/local/bin/shellcheck

FROM base as mdbook
RUN curl -L "https://github.com/rust-lang/mdBook/releases/download/v0.4.4/mdbook-v0.4.4-x86_64-unknown-linux-gnu.tar.gz" | tar -xvz mdbook \
 && curl -L "https://github.com/Michael-F-Bryan/mdbook-linkcheck/releases/download/v0.7.0/mdbook-linkcheck-v0.7.0-x86_64-unknown-linux-gnu.tar.gz" | tar -xvz mdbook-linkcheck \
 && mv mdbook* /usr/local/bin/

FROM base
RUN wget -q https://packages.microsoft.com/config/ubuntu/20.04/packages-microsoft-prod.deb \
 && dpkg -i packages-microsoft-prod.deb \
 && rm packages-microsoft-prod.deb \
 && add-apt-repository universe

RUN apt-get update \
 && apt-get install -y bash-completion \ 
                       fish \
                       elvish \
                       powershell \
                       python3-pip \
                       zsh \
                       expect

RUN pip3 install --no-cache-dir --disable-pip-version-check xonsh \
 && ln -s $(which xonsh) /usr/bin/xonsh

RUN pwsh -Command "Install-Module PSScriptAnalyzer -Scope AllUsers -Force"

COPY --from=bat /usr/local/bin/* /usr/local/bin/
COPY --from=elvish /usr/local/bin/* /usr/local/bin/
COPY --from=goreleaser /usr/local/bin/* /usr/local/bin/
COPY --from=ion /ion/target/release/ion /usr/local/bin/
COPY --from=nu /usr/local/bin/* /usr/local/bin/
COPY --from=mdbook /usr/local/bin/* /usr/local/bin/
COPY --from=oil /usr/local/bin/* /usr/local/bin/
COPY --from=shellcheck /usr/local/bin/* /usr/local/bin/

RUN ln -s /carapace/example/example /usr/local/bin/example

# bash
RUN echo -e "\n\
PS1=$'\e[0;36mcarapace-bash \e[0m'\n\
source /usr/share/bash-completion/bash_completion \n\
source <(example _carapace)" \
       > ~/.bashrc

# fish
RUN mkdir -p ~/.config/fish \
 && echo -e "\n\
function fish_prompt \n\
    set_color cyan \n\
    echo -n 'carapace-fish ' \n\
    set_color normal\n\
end\n\
mkdir -p ~/.config/fish/completions\n\
example _carapace fish | source" \
       > ~/.config/fish/config.fish

# elvish
RUN mkdir -p ~/.elvish/lib \
 && echo -e "\
edit:prompt = { printf  'carapace-elvish ' } \n\
eval (example _carapace|slurp)" \
  > ~/.elvish/rc.elv

# ion
RUN mkdir -p ~/.config/ion \
 && echo -e "\
fn PROMPT\n\
    printf 'carapace-ion '\n\
end" \
  > ~/.config/ion/initrc

# oil
RUN mkdir -p ~/.config/oil \
 && echo -e "\n\
PS1='carapace-oil '\n\
source <(example _carapace)" \
       > ~/.config/oil/oshrc

# powershell
RUN mkdir -p ~/.config/powershell \
 && echo -e "\n\
function prompt {Write-Host \"carapace-powershell\" -NoNewLine -ForegroundColor 3; return \" \"}\n\
Set-PSReadlineKeyHandler -Key Tab -Function MenuComplete\n\
example _carapace | out-string | Invoke-Expression" \
       > ~/.config/powershell/Microsoft.PowerShell_profile.ps1

# xonsh
RUN mkdir -p ~/.config/xonsh \
 && echo -e "\n\
\$PROMPT='carapace-xonsh '\n\
\$COMPLETIONS_CONFIRM=True\n\
exec(\$(example _carapace xonsh))"\
  > ~/.config/xonsh/rc.xsh

# zsh
RUN echo -e "\n\
PS1=$'%{\e[0;36m%}carapace-zsh %{\e[0m%}'\n\
\n\
zstyle ':completion:*' menu select \n\
zstyle ':completion:*' matcher-list 'm:{a-zA-Z}={A-Za-z}' 'r:|=*' 'l:|=* r:|=*' \n\
\n\
autoload -U compinit && compinit \n\
source <(example _carapace zsh)"  > ~/.zshrc
