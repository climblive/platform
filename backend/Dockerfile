FROM ubuntu:latest AS usersrc
RUN useradd -u 10001 climblive

FROM scratch

COPY --from=usersrc /etc/passwd /etc/passwd

COPY climblive /climblive
ENTRYPOINT ["/climblive"]
