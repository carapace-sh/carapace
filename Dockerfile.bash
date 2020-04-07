FROM golang

RUN echo "\n\
PS1=$'\e[0;36mcarapace \e[0m'\n\
source <(example _carapace bash)" \
       > /root/.bashrc

RUN ln -s /carapace/example/example /usr/local/bin/example


ENTRYPOINT [ "bash" ]

