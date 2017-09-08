FROM golang:1.8.3-alpine

RUN mkdir /app
ADD ./dist /app/

EXPOSE 3030

ENTRYPOINT /app/run.sh

# docker build -t build-monitor-v2:latest .
# docker run --name build-monitor-v2 -p 3030:3030 --rm build-monitor-v2