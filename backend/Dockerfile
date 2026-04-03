FROM ubuntu:latest AS usersrc
RUN useradd -u 10001 climblive

FROM scratch

COPY --from=usersrc /etc/passwd /etc/passwd

USER climblive

COPY climblive /climblive
ENTRYPOINT ["/climblive"]