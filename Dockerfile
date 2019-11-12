FROM scratch

COPY ./.bin/container-linux-node-labeler-docker /container-linux-node-labeler

ENTRYPOINT [ "/container-linux-node-labeler" ]
