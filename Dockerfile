FROM golang:1.13-alpine

# Add Maintainer Info
LABEL maintainer="Binh Nguyen"

# install AWSCLI
RUN apk --no-cache update && \
    apk --no-cache add python py-pip py-setuptools ca-certificates groff bash git jq file curl && \
    pip --no-cache-dir install awscli && \
    rm -rf /var/cache/apk/*

# Set the Current Working Directory inside the container
WORKDIR /app
COPY kubestorm .

# Expose port 8080 to the outside world
EXPOSE 8080

CMD ["./kubestorm"]
