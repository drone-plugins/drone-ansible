FROM plugins/base:linux-arm64

LABEL maintainer="Drone.IO Community <drone-dev@googlegroups.com>" \
  org.label-schema.name="Drone Ansible" \
  org.label-schema.vendor="Drone.IO Community" \
  org.label-schema.schema-version="1.0"

RUN apk add --no-cache bash git curl rsync openssh-client py-pip py-requests python2-dev libffi-dev libressl libressl-dev build-base && \
  pip install -U pip && \
  pip install ansible==2.7.6 && \
  apk del python2-dev libffi-dev libressl-dev build-base

ADD release/linux/arm64/drone-ansible /bin/
ENTRYPOINT ["/bin/drone-ansible"]
