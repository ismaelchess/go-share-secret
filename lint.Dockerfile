FROM golangci/golangci-lint:v1.33.0
ARG NETRC
RUN echo "${NETRC}" > /root/.netrc
CMD [ "go", "version" ]
