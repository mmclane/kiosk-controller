FROM alpine:3.15

ENV APP_DIR=/app
RUN mkdir -p $APP_DIR
WORKDIR $APP_DIR


COPY ./bin/kiosk_controller ./kiosk_controller
COPY kiosk_config.json ./kiosk_config.json

ENV KIOSK_CONFIG=$APP_DIR/kiosk_config.json

EXPOSE 8090

ENTRYPOINT [ "/app/kiosk_controller" ]
