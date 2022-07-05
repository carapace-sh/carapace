FROM cimg/go:1.18.1 as base
LABEL org.opencontainers.image.source https://github.com/rsteube/carapace
USER root

FROM base as bat
ARG version=0.20.0
RUN curl -L https://github.com/sharkdp/bat/releases/download/v${version}/bat-v${version}-x86_64-unknown-linux-gnu.tar.gz \
  | tar -C /usr/local/bin/ --strip-components=1  -xvz bat-v${version}-x86_64-unknown-linux-gnu/bat \
  && chmod +x /usr/local/bin/bat

FROM base as ble
RUN git clone --recursive https://github.com/akinomyoga/ble.sh.git \
 && apt-get update && apt-get install gawk \
 && make -C ble.sh

FROM base as elvish
ARG version=0.18.0
RUN curl https://dl.elv.sh/linux-amd64/elvish-v${version}.tar.gz | tar -xvz \
  && mv elvish-* /usr/local/bin/elvish

FROM base as goreleaser
ARG version=1.8.3
RUN curl -L https://github.com/goreleaser/goreleaser/releases/download/v${version}/goreleaser_Linux_x86_64.tar.gz | tar -xvz goreleaser \
  && mv goreleaser /usr/local/bin/goreleaser

FROM rsteube/ion-poc as ion-poc
#FROM rust as ion
#ARG version=master
#RUN git clone --single-branch --branch "${version}" --depth 1 https://gitlab.redox-os.org/redox-os/ion/ \
# && cd ion \
# && RUSTUP=0 make # By default RUSTUP equals 1, which is for developmental purposes \
# && sudo make install prefix=/usr \
# && sudo make update-shells prefix=/usr

FROM base as nushell
ARG version=0.61.0
RUN curl -L https://github.com/nushell/nushell/releases/download/${version}/nu_${version//./_}_linux.tar.gz | tar -xvz \
 && mv nu_${version//./_}_linux/nushell-${version}/nu* /usr/local/bin

FROM base as oil
ARG version=0.9.9
RUN apt-get update && apt-get install -y libreadline-dev
RUN curl https://www.oilshell.org/download/oil-${version}.tar.gz | tar -xvz \
  && cd oil-*/ \
  && ./configure \
  && make \
  && ./install

FROM base as shellcheck
ARG version=stable
RUN wget -qO- "https://github.com/koalaman/shellcheck/releases/download/${version?}/shellcheck-${version?}.linux.x86_64.tar.xz" | tar -xJv shellcheck-stable/shellcheck \
  && mv shellcheck-stable/shellcheck /usr/local/bin/ \
  && chmod +x /usr/local/bin/shellcheck

FROM base as starship
ARG version=1.6.3
RUN wget -qO- "https://github.com/starship/starship/releases/download/v${version}/starship-x86_64-unknown-linux-gnu.tar.gz" | tar -xvz starship \
 && mv starship /usr/local/bin/

FROM base as vivid
ARG version=0.8.0
RUN wget -qO- "https://github.com/sharkdp/vivid/releases/download/v${version}/vivid-v${version}-x86_64-unknown-linux-gnu.tar.gz" | tar -xvz vivid-v${version}-x86_64-unknown-linux-gnu/vivid \
 && mv vivid-v${version}-x86_64-unknown-linux-gnu/vivid /usr/local/bin/

FROM base as mdbook
ARG version=0.4.18
RUN curl -L "https://github.com/rust-lang/mdBook/releases/download/v${version}/mdbook-v${version}-x86_64-unknown-linux-gnu.tar.gz" | tar -xvz mdbook \
  && curl -L "https://github.com/Michael-F-Bryan/mdbook-linkcheck/releases/download/v0.7.0/mdbook-linkcheck-v0.7.0-x86_64-unknown-linux-gnu.tar.gz" | tar -xvz mdbook-linkcheck \
  && mv mdbook* /usr/local/bin/

FROM base
RUN wget -q https://packages.microsoft.com/config/ubuntu/20.04/packages-microsoft-prod.deb \
  && dpkg -i packages-microsoft-prod.deb \
  && rm packages-microsoft-prod.deb \
  && add-apt-repository universe

RUN apt-get update \
  && apt-get install -y fish \
  elvish \
  powershell \
  python3-pip \
  tcsh \
  zsh \
  expect

RUN pip3 install --no-cache-dir --disable-pip-version-check xonsh prompt_toolkit \
  && ln -s $(which xonsh) /usr/bin/xonsh

RUN pwsh -Command "Install-Module PSScriptAnalyzer -Scope AllUsers -Force"

COPY --from=bat /usr/local/bin/* /usr/local/bin/
COPY --from=ble /home/circleci/project/ble.sh /opt/ble.sh
COPY --from=elvish /usr/local/bin/* /usr/local/bin/
COPY --from=goreleaser /usr/local/bin/* /usr/local/bin/
#COPY --from=ion /ion/target/release/ion /usr/local/bin/
COPY --from=ion-poc /usr/local/bin/ion /usr/local/bin/
COPY --from=nushell /usr/local/bin/* /usr/local/bin/
COPY --from=mdbook /usr/local/bin/* /usr/local/bin/
COPY --from=oil /usr/local/bin/* /usr/local/bin/
COPY --from=shellcheck /usr/local/bin/* /usr/local/bin/
COPY --from=starship /usr/local/bin/* /usr/local/bin/
COPY --from=vivid /usr/local/bin/* /usr/local/bin/

RUN ln -s /carapace/example/example /usr/local/bin/example

RUN echo -e "\n\
[shell]\n\
disabled = false\n\
unknown_indicator = \"oil\"" \
  > ~/.config/starship.toml

# bash
RUN echo -e "\n\
export SHELL=bash\n\
export STARSHIP_SHELL=bash\n\
export LS_COLORS=\"\$(vivid generate dracula)\"\n\
[[ ! -z \$BLE ]] && source /opt/ble.sh/out/ble.sh \n\
eval \"\$(starship init bash)\"\n\
source <(\${TARGET} _carapace)" \
  > ~/.bashrc

# fish
RUN mkdir -p ~/.config/fish \
  && echo -e "\n\
  set SHELL 'fish'\n\
  set STARSHIP_SHELL 'fish'\n\
  set LS_COLORS (vivid generate dracula)\n\
  starship init fish | source \n\
  mkdir -p ~/.config/fish/completions\n\
  \$TARGET _carapace fish | source" \
  > ~/.config/fish/config.fish

# elvish
RUN mkdir -p ~/.elvish/lib \
  && echo -e "\
  set-env SHELL elvish\n\
  set-env STARSHIP_SHELL elvish\n\
  set-env LS_COLORS (vivid generate dracula)\n\
  set edit:prompt = { starship prompt }\n\
  eval (\$E:TARGET _carapace|slurp)" \
  > ~/.elvish/rc.elv

# ion
RUN mkdir -p ~/.config/ion \
  && echo -e "\
  fn PROMPT\n\
  printf 'carapace-ion '\n\
  end" \
  > ~/.config/ion/initrc

# nushell
RUN touch /carapace.nu \
  && mkdir -p ~/.config/nushell \
  && starship init nushell > ~/.config/nushell/starship.nu \
  && echo -e "\
ln -s \$env.TARGET /tmp/target \n\
/tmp/target _carapace | save /carapace.nu \n\
source /carapace.nu \n\
" > ~/.config/nushell/config.nu \
  && echo -e "\
source ~/.config/nushell/starship.nu \n\
" > ~/.config/nushell/env.nu

# oil
RUN mkdir -p ~/.config/oil \
  && echo -e "\n\
  export SHELL='oil'\n\
  export STARSHIP_SHELL='oil'\n\
  export LS_COLORS=\"\$(vivid generate dracula)\"\n\
  PS1=\"\$(starship prompt)\"\n\
  source <(\${TARGET} _carapace)" \
  > ~/.config/oil/oshrc

# powershell
RUN mkdir -p ~/.config/powershell \
  && echo -e "\n\
  \$env:SHELL = 'powershell'\n\
  \$env:STARSHIP_SHELL = 'powershell'\n\
  \$env:LS_COLORS = (&vivid generate dracula)\n\
  Invoke-Expression (&starship init powershell)\n\
  Set-PSReadlineKeyHandler -Key Tab -Function MenuComplete\n\
  & \$Env:TARGET _carapace | out-string | Invoke-Expression" \
  > ~/.config/powershell/Microsoft.PowerShell_profile.ps1

# tcsh
RUN  echo -e "\n\
  eval `starship init tcsh`\n\
  set autolist\n\
  eval "'`'"\${TARGET} _carapace"'`'"" \
  > ~/.tcshrc

# xonsh
RUN mkdir -p ~/.config/xonsh \
  && echo -e "\n\
\$SHELL=\"xonsh\"\n\
\$STARSHIP_SHELL=\"xonsh\"\n\
\$LS_COLROS=\$(vivid generate dracula)\n\
\$PROMPT=lambda: \$(starship prompt)\n\
\$COMPLETIONS_CONFIRM=True\n\
exec(\$(\$TARGET _carapace xonsh))"\
  > ~/.config/xonsh/rc.xsh

# zsh
RUN echo -e "\n\
  export SHELL=zsh\n\
  export STARSHIP_SHELL=zsh\n\
  export LS_COLORS=\"\$(vivid generate dracula)\"\n\
  eval \"\$(starship init zsh)\"\n\
  \n\
  zstyle ':completion:*' menu select \n\
  zstyle ':completion:*' matcher-list 'm:{a-zA-Z}={A-Za-z}' 'r:|=*' 'l:|=* r:|=*' \n\
  \n\
  autoload -U compinit && compinit \n\
  source <(\$TARGET _carapace zsh)"  > ~/.zshrc

ENV TERM xterm
RUN echo -e "#"'!'"/bin/bash\n\
  export PATH=\${PATH}:\$(dirname \"\${TARGET}\")\n\
  exec \"\$@\"" \
  > /entrypoint.sh \
  && chmod a+x /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]
