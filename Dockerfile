ARG BASE_IMAGE=ubuntu:24.04

FROM golang AS build
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /app/kairos-init .


FROM scratch AS kairos-init
COPY --from=build /app/kairos-init /kairos-init

FROM ${BASE_IMAGE}
ARG MODEL=generic
ARG VARIANT=core
ARG FRAMEWORK_VERSION=""
ARG TRUSTED_BOOT=false

COPY --from=kairos-init /kairos-init /kairos-init
RUN /kairos-init -l debug -s install -m "${MODEL}" -v "${VARIANT}" -f "${FRAMEWORK_VERSION}" -t "${TRUSTED_BOOT}"
RUN /kairos-init -l debug -s init -m "${MODEL}" -v "${VARIANT}" -f "${FRAMEWORK_VERSION}" -t "${TRUSTED_BOOT}"
RUN rm /kairos-init