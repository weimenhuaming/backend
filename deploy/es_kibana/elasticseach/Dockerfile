FROM docker.elastic.co/elasticsearch/elasticsearch:9.3.3

LABEL maintainer="chenaqi"

RUN elasticsearch-plugin install --batch \
https://get.infini.cloud/elasticsearch/analysis-ik/9.3.3 \
&& elasticsearch-plugin list