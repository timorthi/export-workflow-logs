# https://www.sethvargo.com/writing-github-actions-in-go/
FROM golang:1.19 AS builder

# Install upx to compress the binary
RUN apt-get update && apt-get -y install upx-ucl

ENV CGO_ENABLED=0

WORKDIR /src/export-workflow-logs
COPY . .

# Compile the action - the added flags instruct Go to produce a standalone binary
RUN go build \
  -a \
  -trimpath \
  -ldflags "-s -w -extldflags '-static'" \
  -installsuffix cgo \
  -tags netgo \
  -o /bin/action \
  .

# Strip any symbols (this is not a library) then compress the action
RUN strip /bin/action
RUN upx -q -9 /bin/action


FROM scratch

# Copy over SSL certificates from the first step - this is required
# if our code makes any outbound SSL connections because it contains
# the root CA bundle.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy over the compiled action from the first step
COPY --from=builder /bin/action /bin/action

# Specify the container's entrypoint as the action
ENTRYPOINT ["/bin/action"]
