FROM amery/docker-golang-modd:latest as intermediate

# custom user
ARG USER_NAME=nobody
ARG USER_UID
ARG USER_GID

RUN [ nobody = "$USER_NAME" ] || addgroup -S -g "$USER_GID" "$USER_NAME"
RUN [ nobody = "$USER_NAME" ] || adduser -S -D -G "$USER_NAME" -g "$USER_NAME" -u "$USER_UID" "$USER_NAME"
RUN cp /etc/passwd root/etc/

# empty image
FROM scratch
COPY --from=intermediate /go/root /

USER ${USER_NAME}
ENTRYPOINT [ "/modd", "-f", "/etc/modd.conf" ]
