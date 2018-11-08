FROM golang
EXPOSE 8888
RUN mkdir -p /app/baoguan
ADD . /app/baoguan
RUN chmod 777 -R /app/baoguan/
ENTRYPOINT ["/app/baoguan/baoguan"]
