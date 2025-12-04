FROM fedora:43

WORKDIR /app/src
COPY . .

ENV RAW_PATH="/usr/local/go/bin:/root/.bun/bin"

RUN dnf install make unzip wget -y
RUN dnf clean all

RUN wget https://go.dev/dl/go1.25.5.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.25.5.linux-amd64.tar.gz
RUN curl -fsSL https://bun.sh/install | sh

RUN rm go1.25.5.linux-amd64.tar.gz

RUN PATH="$PATH:${RAW_PATH}" ./configure
RUN PATH="$PATH:${RAW_PATH}" make

RUN cp ./sample-app /app

WORKDIR /app
RUN rm -rf src/
RUN rm -rf /root/.bun
RUN rm -rf /usr/local/go

ENTRYPOINT [ "/app/sample-app" ]
