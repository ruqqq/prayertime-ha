ARG BUILD_FROM
FROM $BUILD_FROM

ENV LANG C.UTF-8

# Copy data for add-on
COPY prayertime_ha /
COPY run.sh /
RUN chmod a+x /run.sh

CMD [ "/run.sh" ]
