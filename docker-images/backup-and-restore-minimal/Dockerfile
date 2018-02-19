FROM golang:1.10

RUN \
  apt-get update && \
  apt-get install -y sudo unzip jq sshuttle parallel vim && \
  rm -rf /var/lib/apt/lists/*

RUN go get github.com/onsi/ginkgo/ginkgo
RUN go get github.com/onsi/gomega
RUN go get github.com/pivotal-cf/om
RUN go get -u github.com/golang/dep/cmd/dep

RUN (curl -L "https://cli.run.pivotal.io/stable?version=6.32.0&release=linux64-binary&source=github" | tar -zx) && \
   mv cf /usr/bin && \
   chmod +x /usr/bin/cf

RUN \
  wget https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-2.0.23-linux-amd64 --output-document="/usr/bin/bosh-cli" && \
  cd /usr/bin && \
  echo "2c22556ab2aee6da659beda7110e28c77bc0ab2a  bosh-cli" | sha1sum -c - && \
  chmod +x bosh-cli
