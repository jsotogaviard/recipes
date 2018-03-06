# From image that contains go
FROM golang:1.9

# Create the directory
RUN mkdir /go/src/jsotogaviard-api-test

# Copy all the relevant information
ADD . /go/src/jsotogaviard-api-test

# Configure as current directory
WORKDIR /go/src/jsotogaviard-api-test

# Fetch the dependencies and build the binary file
RUN go get -v
RUN go build -o main .

# Launch the main
CMD ["/go/src/jsotogaviard-api-test/main"]