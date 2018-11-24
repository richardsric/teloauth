# Note: The newer busybox:glibc is missing libpthread.so.0.
FROM ngrok_root:latest
MAINTAINER iYOCHU Nigeria Ltd
#RUN apk add --no-cache ca-certificates
# Add config script.
COPY ngrok.yml /home/ngrok/.ngrok2/
# Add the executable
COPY . /
EXPOSE 4000

RUN chown -R ngrok:ngrok /home/ngrok \
&& chmod -R go=u,go-w /home/ngrok \
&& chmod go= /home/ngrok

USER ngrok
#ENTRYPOINT ["/entrypoint.sh"]
CMD ["/entrypoint.sh"]