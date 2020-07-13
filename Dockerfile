FROM scratch

LABEL maintainer="The Prometheus Authors <prometheus-developers@googlegroups.com>"

COPY promql-langserver /bin/promql-langserver

ENTRYPOINT [ "/bin/promql-langserver" ]
