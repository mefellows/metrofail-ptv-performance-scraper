FROM scratch
MAINTAINER Matt Fellows <matt.fellows@onegeek.com.au>
ADD ptvperf ptvperf
ENTRYPOINT ["/ptvperf"]
