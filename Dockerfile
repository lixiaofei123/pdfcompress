FROM golang AS build
WORKDIR /build
COPY . .
ENV GOPROXY https://goproxy.io,direct
ENV CGO_ENABLED=0
RUN go build -o pdfcompress

FROM ubuntu
RUN apt-get -y update && \
  apt-get -y --no-install-recommends install ghostscript && \
  rm -rf /var/lib/apt/lists/*
RUN mkdir -p /opt/pdfcompress
RUN mkdir -p /opt/pdfcompress/input
RUN mkdir -p /opt/pdfcompress/output
COPY --from=build /build/pdfcompress /opt/pdfcompress/pdfcompress
COPY static /opt/pdfcompress/static
EXPOSE 8082
WORKDIR /opt/pdfcompress/
ENTRYPOINT ["./pdfcompress"]
