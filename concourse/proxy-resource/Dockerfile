ARG TARGET_IMAGE

FROM golang:alpine as builder
COPY ./in /go/src/github.com/concourse/trigger-resource/in
RUN ls /go/src/github.com/concourse/trigger-resource
RUN go build -o /assets/in github.com/concourse/trigger-resource/in

FROM $TARGET_IMAGE as resource
ARG TARGET_IMAGE
RUN mv /opt/resource /opt/target-resource
COPY --from=builder /assets /opt/resource
ENV PROXY_RESOURCE_PREFIX /opt/target-resource

FROM resource
