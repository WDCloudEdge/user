FROM alpine:3.4

ENV HATEAOS user
ENV USER_DATABASE mongodb
ENV MONGO_HOST user-db:27017

WORKDIR /
COPY user /

RUN	chmod +x /user && \
	chown 644 /user

USER ${SERVICE_USER}

CMD ["/user", "-port=8084"]
