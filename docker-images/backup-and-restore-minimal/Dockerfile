FROM golang:1.9

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
