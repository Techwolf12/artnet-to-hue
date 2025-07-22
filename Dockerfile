FROM scratch
COPY artnet-to-hue /usr/bin/artnet-to-hue
ENTRYPOINT ["/usr/bin/artnet-to-hue"]