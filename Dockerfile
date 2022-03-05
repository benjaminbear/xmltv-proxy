FROM debian:11-slim

RUN apt-get update && apt-get install -y \
    wget \
    tzdata \
    && rm -rf /var/lib/apt/lists/* \
    mkdir /epg

WORKDIR /epg
COPY xmltv-proxy /epg/xmltv-proxy

EXPOSE 8080

CMD ["bash", "-c", "/epg/xmltv-proxy"]