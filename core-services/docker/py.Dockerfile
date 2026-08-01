FROM ubuntu:24.04

# create a group + “system” user, no home directory
RUN useradd -ms /bin/bash application

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates python3.12 python3-pip && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /opt/app
COPY dist/ ./dist/

# dist/ is the wheelhouse built by the CI "build" stage (pip wheel . -w dist/):
# the app's own wheel plus every dependency's wheel, already resolved. Installing
# with --no-index/--find-links means this step needs no network access and never
# re-resolves versions that were already pinned and fetched upstream.
RUN pip install --no-cache-dir --break-system-packages --no-index --find-links=dist analytic-api_mrflick72 && \
    chown -R application:application /opt/app

USER application

ENTRYPOINT ["app"]
