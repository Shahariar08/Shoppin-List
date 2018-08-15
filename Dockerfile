FROM busybox:glibc
ADD main /bin/main
EXPOSE 8800
CMD ["main"]
