# build static binary
FROM golang:1.15.2-alpine3.12 as builder 

# hadolint ignore=DL3018
RUN apk --no-cache add  \
    ca-certificates \
    git 

WORKDIR /go/src/github.com/bots-house/github-actions-upstream-watch

# download dependencies 
COPY go.mod go.sum ./
RUN go mod download 

COPY . .

ARG REVISION

# compile 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags="-w -s -extldflags \"-static\" -X \"main.revision=${REVISION}\"" -a \
      -tags timetzdata \
      -o /bin/github-actions-upstream-watch .

RUN mkdir /data

# run 
FROM scratch


COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/github-actions-upstream-watch /bin/github-actions-upstream-watch
COPY --from=builder /data /data 

VOLUME [ "/data" ]


ENV STATE=/data/state.sha
ENTRYPOINT [ "/bin/github-actions-upstream-watch" ]