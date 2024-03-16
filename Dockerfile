# Use image with golang and docker
FROM localhost:5000/golang-docker:latest

# Build the application
RUN go build

# Move queries and configs elsewhere
RUN mkdir /etc/devprivops
RUN mv .devprivops /etc/devprivops

# Cleanup
RUN mv devprivops /bin/devprivops
RUN rm -rf .