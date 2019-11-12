FROM scratch

COPY ./.bin/container-linux-node-labeler /container-linux-node-labeler

ENTRYPOINT [ "/container-linux-node-labeler" ]
