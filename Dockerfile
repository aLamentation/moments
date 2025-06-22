# 构建阶段
FROM golang:1.23.3-alpine AS build

ARG VERSION
ARG COMMIT_ID

WORKDIR /app
RUN apk add --no-cache build-base tzdata

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./

# 👇 加上这一句：确保 public/ 存在
COPY front/.output/public ./public

RUN go build -tags prod -ldflags="-s -w -X main.version=${VERSION} -X main.commitId=${COMMIT_ID}" -o moments

FROM alpine
WORKDIR /app/data
RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV PORT=3000
ENV TZ=Asia/Shanghai
COPY --from=build /app/moments /app/moments
RUN chmod +x /app/moments
EXPOSE 3000
CMD ["/app/moments"]

#FROM golang:1.23.3-alpine AS backend
#ARG VERSION
#ARG COMMIT_ID
#WORKDIR /app
#RUN apk add --no-cache build-base tzdata
#COPY backend/go.mod .
#COPY backend/go.sum .
#RUN go mod download
#COPY backend/. .
#COPY --from=front /app/.output/public /app/public
#RUN go build -tags prod -ldflags="-s -w -X main.version=${VERSION} -X main.commitId=${COMMIT_ID}" -o /app/moments
#
#FROM alpine
#WORKDIR /app/data
#RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
#ENV PORT=3000
#ENV TZ=Asia/Shanghai
#COPY --from=backend /app/moments /app/moments
#RUN chmod +x /app/moments
#EXPOSE 3000
#CMD ["/app/moments"]
